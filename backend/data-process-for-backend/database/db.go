package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error

	// Set logger mode based on environment (e.g., fallback to Silent/Warn in production)
	dbLogger := logger.Default.LogMode(logger.Info) // Use logger.Warn or logger.Silent for production
	if os.Getenv("ENV") == "production" {
		dbLogger = logger.Default.LogMode(logger.Error) // Only log true database errors
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Neon PostgreSQL: %v", err)
	}

	log.Println("Successfully connected to Neon DB via GORM!")
}