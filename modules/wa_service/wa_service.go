package wa_service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wa_bot_service/config"
	"wa_bot_service/db/models"
	logging "wa_bot_service/modules/logger"
	"wa_bot_service/modules/utils"
	"wa_bot_service/services/raw_transaction_service"
	test_series_raw_question_service "wa_bot_service/services/test_series_service"
	"wa_bot_service/services/user_service"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WAClientService struct {
	client *whatsmeow.Client
	logger *slog.Logger
}

func NewWAClientService() (*WAClientService, error) {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:wa_service_datastore.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}
	session, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(session, clientLog)

	return &WAClientService{client: client, logger: logging.NewLogger("wa_service")}, nil
}

func (svc *WAClientService) Start() error {
	svc.client.AddEventHandler(func(evt interface{}) {
		svc.eventHandler(svc.client, evt)
	})

	if svc.client.Store.ID == nil {
		qrChan, _ := svc.client.GetQRChannel(context.Background())
		if err := svc.client.Connect(); err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				svc.logger.Info("QR code", "code", evt.Code)
			} else {
				svc.logger.Info("Login event:", "evnet", evt.Event)
			}
		}
	} else {
		if err := svc.client.Connect(); err != nil {
			return err
		}
	}
	return nil
}

func (svc *WAClientService) Stop() {
	svc.client.Disconnect()
}

func (svc *WAClientService) eventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsFromMe {
			return
		}

		jid := v.Info.Chat
		if v.Info.Chat.Server == "status" || v.Info.Chat.Server == "broadcast" {
			return
		}
		fullPhoneNumber := svc.GetPhoneNumberWithCountryCode(jid.String())
		phoneNumber := svc.GetPhoneNumberWithoutCountryCode(jid.String())
		user, err := user_service.GetUserByPhoneNumber(phoneNumber)
		if err != nil {
			if defaultMessage, ok := WAMessageTemplates["defaultErrorMessage"].(string); ok {
				svc.SendMessage(fullPhoneNumber, defaultMessage)
				return
			} else {
				svc.logger.Error("defaultErrorMessage template is not a string")
			}
			return
		}
		msgLowerCase := strings.ToLower(v.Message.GetConversation())
		if msgLowerCase == "1" {
			err := svc.handleRegisterUser(fullPhoneNumber)
			if err != nil {
				svc.logger.Error(err.Error())
			}
		} else if user != nil {
			// TODO: WRITE LIMITER HERE

			svc.logger.Info("jid", "jid", jid)
			isAashi := user.ID.String() == "fa4ebfb0-aa82-11ef-990b-037e8c7d24b0"
			// Handle media messages
			if v.Message.GetImageMessage() != nil {
				var mediaType = "raw_transaction"
				if isAashi {
					mediaType = "raw_question"
				}
				filePath, downloadMediaError := downloadMedia(client, v.Info.Chat, v.Message.GetImageMessage(), mediaType)
				if downloadMediaError == nil {
					keyword := "uploads/"
					index := strings.Index(*filePath, keyword)
					var formattedFilePath string
					if index != -1 {
						str := *filePath
						formattedFilePath = str[index:]
					} else {
						// todo: do better handling
						fmt.Println("Keyword not found in path")
					}
					if isAashi {
						// todo: add test series message handler
						if err := svc.testSeriesMessageHandler(true, &formattedFilePath); err != nil {
							if transactions, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
								svc.SendMessage(fullPhoneNumber, transactions["invalid_input"])
								return
							} else {
								svc.logger.Error("transactions template is not a map[string]string")
								return
							}
						}
						return
					} else {
						if err := svc.transactionMessageHandler(user, true, msgLowerCase, &formattedFilePath); err != nil {
							if transactions, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
								svc.SendMessage(fullPhoneNumber, transactions["invalid_input"])
								return
							} else {
								svc.logger.Error("transactions template is not a map[string]string")
								return
							}
						}
						if transactionsWaMsg, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
							svc.SendMessage(fullPhoneNumber, transactionsWaMsg["input_received"])
							return
						} else {
							svc.logger.Error("transactions template is not a map[string]string")
							return
						}
						// todo: add transaction message handler
					}

				} else {
					if transactions, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
						svc.SendMessage(fullPhoneNumber, transactions["invalid_input"])
						return
					} else {
						svc.logger.Error("transactions template is not a map[string]string")
						return
					}
				}
			} else if msgLowerCase != "" {
				if err := svc.transactionMessageHandler(user, false, msgLowerCase, nil); err != nil {
					if transactions, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
						svc.SendMessage(fullPhoneNumber, transactions["invalid_input"])
						return
					} else {
						svc.logger.Error("transactions template is not a map[string]string")
						return
					}
				}
				if transactionsWaMsg, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
					svc.SendMessage(fullPhoneNumber, transactionsWaMsg["input_received"])
					return
				} else {
					svc.logger.Error("transactions template is not a map[string]string")
					return
				}
			} else if v.Message.GetVideoMessage() != nil || v.Message.GetDocumentMessage() != nil || v.Message.GetAudioMessage() != nil {
				// downloadMedia(client, v.Info.Chat, v.Message.GetDocumentMessage(), "document")
				if transactions, ok := WAMessageTemplates["transactions"].(map[string]string); ok {
					svc.SendMessage(fullPhoneNumber, transactions["invalid_input"])
					return
				} else {
					svc.logger.Error("transactions template is not a map[string]string")
					return
				}
			} else {
				svc.logger.Error("unhandled wa event", evt)
				return
			}
		} else {
			if defaultMessage, ok := WAMessageTemplates["defaultMessage"].(string); ok {
				svc.SendMessage(fullPhoneNumber, defaultMessage)
				return
			} else {
				svc.logger.Error("defaultMessage template is not a string")
				return
			}
		}
	}
}

func (svc *WAClientService) testSeriesMessageHandler(hasMedia bool, filePath *string) error {
	if filePath == nil {
		return fmt.Errorf("nil filePath")
	}
	testSeriesRawQuestionService := test_series_raw_question_service.NewTestSeriesRawQuestionService()
	data := "no-media"
	if hasMedia {
		data = *filePath
	}
	_, err := testSeriesRawQuestionService.CreateTestSeriesRawQuestion(data)

	return err
}

func (svc *WAClientService) transactionMessageHandler(user *models.User, hasMedia bool, messageText string, filePath *string) error {
	if hasMedia && filePath == nil {
		return fmt.Errorf("nil filePath")
	}
	rawTransactionService := raw_transaction_service.NewRawTransactionService()
	rawTransactionType := models.RawTransactionType.WAText
	if hasMedia {
		rawTransactionType = models.RawTransactionType.WAImage
	}
	transactionData := func() string {
		if hasMedia {
			return *filePath
		}
		return messageText
	}()
	status := func() string {
		if hasMedia {
			return models.RawTransactionStatuses.PendingTextExtraction
		}
		return models.RawTransactionStatuses.Pending
	}()
	_, err := rawTransactionService.CreateRawTransaction(user.ID, rawTransactionType, transactionData, status)

	return err
}

func downloadMedia(client *whatsmeow.Client, chat types.JID, media interface{}, mediaPurpose string) (*string, error) {
	var fileExtension string
	var downloadable whatsmeow.DownloadableMessage

	switch m := media.(type) {
	case *waE2E.ImageMessage:
		fileExtension = ".jpeg"
		downloadable = m
	case *waE2E.VideoMessage:
		fileExtension = ".mp4"
		downloadable = m
	case *waE2E.DocumentMessage:
		fileExtension = filepath.Ext(m.GetFileName())
		downloadable = m
	case *waE2E.AudioMessage:
		fileExtension = ".ogg"
		downloadable = m
	default:
		fmt.Println("Unsupported media type")
		return nil, fmt.Errorf("unsupported media type")
	}

	// Create a file to save the media
	mediaStoragePath := config.AppConfig.MEDIA_UPLOAD_PATH
	uploadsFolderPath := fmt.Sprintf("%s/uploads", mediaStoragePath)
	absFolderPath, absFolderPathErr := filepath.Abs(uploadsFolderPath)
	if absFolderPathErr != nil {
		fmt.Println("Failed to get absolute path:", absFolderPathErr)
		return nil, fmt.Errorf("failed to get absolute path: %v", absFolderPathErr)
	}
	err := utils.EnsureDirExists(absFolderPath)
	if err != nil {
		fmt.Println("Failed to create folder:", err)
		return nil, fmt.Errorf("failed to create folder: %v", err)
	}
	fileName := fmt.Sprintf("%s/%d_%s_%s%s", uploadsFolderPath, time.Now().UnixMilli(), mediaPurpose, chat.User, fileExtension)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Download the media file using DownloadToFile
	err = client.DownloadToFile(downloadable, file)
	if err != nil {
		fmt.Println("Failed to download media:", err)
		return nil, fmt.Errorf("failed to download media: %v", err)
	}

	fmt.Printf("Downloaded %s: %s\n", mediaPurpose, fileName)
	return &fileName, nil
}

func (svc *WAClientService) SendMessage(to string, message string) error {
	jid := types.NewJID(to, "s.whatsapp.net")

	_, err := svc.client.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: &message,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	fmt.Printf("Message sent to %s: %s\n", to, message)
	return nil
}

// todo: send mark read
// func (svc *WAClientService) markRead(msgIDs []string, recipientJID types.JID, senderJID types.JID) error {
// 	err := svc.client.MarkRead(msgIDs, time.Now(), recipientJID, senderJID)
// 	if err != nil {
// 		return fmt.Errorf("failed to mark message as read: %v", err)
// 	} else {
// 		svc.logger.Info("Message marked as read (blue tick)!")
// 		return nil
// 	}
// }

func (svc *WAClientService) handleRegisterUser(fullPhoneNumber string) error {
	phoneNumber := fullPhoneNumber[2:]

	alreadyExist, err := user_service.GetUserByPhoneNumber(phoneNumber)
	if err == nil && alreadyExist != nil {
		if registerTemplates, ok := WAMessageTemplates["register"].(map[string]string); ok {
			msgErr := svc.SendMessage(fullPhoneNumber, registerTemplates["already_registered"])
			return fmt.Errorf("already registered: %v", msgErr)
		} else {
			return fmt.Errorf("register template is not a map[string]string")
		}
	} else if err != nil {
		if registerTemplates, ok := WAMessageTemplates["register"].(map[string]string); ok {
			msgErr := svc.SendMessage(fullPhoneNumber, registerTemplates["error"])
			return fmt.Errorf("error in registeration: %v", msgErr)
		} else {
			return fmt.Errorf("register template is not a map[string]string")
		}
	}

	user, err := user_service.CreateUser(phoneNumber, "+91")
	if err != nil || user == nil {
		registerTemplates, ok := WAMessageTemplates["register"].(map[string]string)
		if !ok {
			return fmt.Errorf("register template is not a map[string]string")
		}
		msgErr := svc.SendMessage(fullPhoneNumber, registerTemplates["error"])
		return fmt.Errorf("error in registeration: %v", msgErr)
	}
	if registerTemplates, ok := WAMessageTemplates["register"].(map[string]string); ok {
		svc.SendMessage(fullPhoneNumber, registerTemplates["success"])
	} else {
		return fmt.Errorf("register template is not a map[string]string")
	}

	return nil
}

func (svc *WAClientService) GetPhoneNumberWithCountryCode(jid string) string {
	// Remove the server part (@s.whatsapp.net)
	phoneWithCountry := strings.Split(jid, "@")[0]

	return phoneWithCountry
}

func (svc *WAClientService) GetPhoneNumberWithoutCountryCode(jid string) string {
	// Remove the server part (@s.whatsapp.net)
	phoneWithCountry := strings.Split(jid, "@")[0]

	// Get the last 10 digits
	if len(phoneWithCountry) > 10 {
		return phoneWithCountry[len(phoneWithCountry)-10:] // Extract last 10 digits
	}
	return phoneWithCountry // Return as is if it's already 10 digits
}
