package db

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/ecom/auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	// 1. Get credentials from Environment Variables (Best Practice)
	// Default to localhost and host port 5433 for local development (outside Docker).
	// When services run inside Docker Compose the env var `DB_HOST` should be set
	// to the service name (e.g. "postgres") and the container-internal port 5432 will be used.
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "ecom_db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		// docker-compose maps host port 5433 -> container 5432 for convenience
		port = "5433"
	}

	// 2. Build Connection String (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	// 3. Connect to Postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 4. AutoMigrate (Magic)
	// This automatically creates the "users" table based on your Struct
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("🚀 Connected to Database & Migrated!")
	return db
}
