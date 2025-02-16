package db_service

import (
	"fmt"
	"log"

	"github.com/Karan0009/go_wa_bot/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global dbClient object
var dbClient *gorm.DB

// NewDBClient creates a new DB client (connection)
func NewDBClient() (*gorm.DB, error) {
	config := config.AppConfig.PGDBConnection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSL,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

// InitializeDBClient initializes the global dbClient with a database connection
func InitializeDBClient() error {
	if dbClient != nil {
		return fmt.Errorf("dbClient is already initialized")
	}

	client, err := NewDBClient()
	if err != nil {
		return err
	}

	// Set the global dbClient
	dbClient = client
	return nil
}

// GetDBClient returns the global dbClient instance
func GetDBClient() *gorm.DB {
	if dbClient == nil {
		log.Fatal("dbClient is not initialized. Call InitializeDBClient first.")
	}
	return dbClient
}
