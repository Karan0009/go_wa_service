package logging

import (
	"fmt"
	"log/slog"
	"os"
	"wa_bot_service/config"
)

// NewLogger creates and initializes a new logger instance with a label.
func NewLogger(label string) *slog.Logger {
	// Create or open the log file
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}

	// Create a log file handler
	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{})

	// Default logger instance with file handler
	logger := slog.New(fileHandler)
	if config.AppConfig.APP_ENV == "development" {
		// Create a console handler for development environment
		consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
		// Add the console handler to the logger
		logger = slog.New(consoleHandler)
	}

	// Add the label to the logger
	logger = logger.With(slog.String("label", label))

	return logger
}
