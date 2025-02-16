package raw_transaction_service

import (
	"errors"

	"github.com/Karan0009/go_wa_bot/db/models"
	db_service "github.com/Karan0009/go_wa_bot/modules/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RawTransactionService struct {
	db *gorm.DB
}

func NewRawTransactionService() *RawTransactionService {
	return &RawTransactionService{db: db_service.GetDBClient()}
}

func (s *RawTransactionService) CreateRawTransaction(userID uuid.UUID, transactionType, transactionData, status string) (*models.RawTransaction, error) {
	rawTransaction := &models.RawTransaction{
		UserID:             userID,
		RawTransactionType: transactionType,
		RawTransactionData: transactionData,
		Status:             status,
	}

	if err := s.db.Omit("deleted_at").Create(rawTransaction).Error; err != nil {
		return nil, err
	}

	return rawTransaction, nil
}

func (s *RawTransactionService) GetRawTransactionByID(id uuid.UUID) (*models.RawTransaction, error) {
	var rawTransaction models.RawTransaction
	if err := s.db.First(&rawTransaction, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &rawTransaction, nil
}

func (s *RawTransactionService) UpdateTransactionStatus(id uuid.UUID, newStatus string) error {
	if err := s.db.Model(&models.RawTransaction{}).
		Where("id = ?", id).
		Update("status", newStatus).Error; err != nil {
		return err
	}
	return nil
}

func (s *RawTransactionService) GetPendingTransactions() ([]models.RawTransaction, error) {
	var transactions []models.RawTransaction
	err := s.db.Where("status = ?", models.RawTransactionStatuses.Pending).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *RawTransactionService) DeleteRawTransaction(id uuid.UUID) error {
	if err := s.db.Delete(&models.RawTransaction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
