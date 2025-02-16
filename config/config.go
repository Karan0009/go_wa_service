package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppConfig *Config // This will be our global config object

type PGDBConnection struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSL      string
}

type Config struct {
	PGDBConnection    PGDBConnection
	APP_ENV           string
	MEDIA_UPLOAD_PATH string
	GRPC_SERVER_PORT  string
	START_WA_CLIENT   bool
}

// LoadConfig initializes the AppConfig global variable with values from environment variables
func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return err
	}

	AppConfig = &Config{
		PGDBConnection: PGDBConnection{
			Host:     getEnv("PG_DB_HOST", "localhost"),
			User:     getEnv("PG_DB_USER", ""),
			Password: getEnv("PG_DB_PASSWORD", ""),
			DBName:   getEnv("PG_DB_NAME", "lekhakaar_main"),
			Port:     getEnv("PG_DB_PORT", "5432"),
			SSL:      getEnv("PG_DB_SSL_MODE", "false"),
		},
		APP_ENV:           getEnv("APP_ENV", "development"),
		MEDIA_UPLOAD_PATH: "../media_storage",
		GRPC_SERVER_PORT:  getEnv("GRPC_SERVER_PORT", "8088"),
		START_WA_CLIENT:   true,
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
