package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Karan0009/go_wa_bot/config"
	db_service "github.com/Karan0009/go_wa_bot/modules/db"
	"github.com/Karan0009/go_wa_bot/modules/grpc_server"
	logging "github.com/Karan0009/go_wa_bot/modules/logger"
	"github.com/Karan0009/go_wa_bot/modules/wa_service"
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
	start_wa_client := config.AppConfig.START_WA_CLIENT

	var waService *wa_service.WAClientService
	if start_wa_client {
		waService, err = wa_service.NewWAClientService()
		if err != nil {
			fmt.Println("Failed to initialize WhatsApp client:", err)
			return
		}
		err = waService.Start()
		if err != nil {
			fmt.Println("Failed to start WhatsApp client:", err)
			return
		}
	}

	grpcServer, grpcServerErr := grpc_server.InitGrpcServer(config.AppConfig.GRPC_SERVER_PORT)
	if grpcServerErr != nil {
		log.Fatalf("❌ Error starting gRPC server: %v", err)
	}

	// Handle shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	grpcServer.GracefulStop()
	fmt.Println("✅ gRPC server stopped gracefully")
	if waService != nil {
		waService.Stop()
	}
	fmt.Println("WhatsApp client disconnected.")
}
