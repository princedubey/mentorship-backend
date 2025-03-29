package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

func InitializeDatabase() {
	dsn := os.Getenv("DATABASE_URL")

	// Configure GORM to be silent about slow queries
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err = gorm.Open(postgres.Open(dsn), config)

	if err != nil {
		log.Fatalf("Failed to connect to database:%v", err)
	}

	log.Println("Database connection established")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}
