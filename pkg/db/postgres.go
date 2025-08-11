package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URI")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	return db
}
