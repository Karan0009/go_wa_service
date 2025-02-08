package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wa_bot_service/config"
	db_service "wa_bot_service/modules/db"
	logging "wa_bot_service/modules/logger"
	wa_service "wa_bot_service/modules/wa_service"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	logging.NewLogger("main").Info("Config loaded")
	err = db_service.InitializeDBClient()
	if err != nil {
		log.Fatalf("Error initializing DB client: %v", err)
	}
	logging.NewLogger("main").Info("Connected to DB")
	start_wa_client := true

	if !start_wa_client {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		return
	}
	waService, err := wa_service.NewWAClientService()
	if err != nil {
		fmt.Println("Failed to initialize WhatsApp client:", err)
		return
	}

	err = waService.Start()
	if err != nil {
		fmt.Println("Failed to start WhatsApp client:", err)
		return
	}

	// Handle shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	waService.Stop()
	fmt.Println("WhatsApp client disconnected.")
}
