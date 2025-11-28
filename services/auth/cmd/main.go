package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecom/auth/internal"
	"github.com/yourusername/ecom/auth/internal/db"
)

func main() {
	// 1. Initialize Database
	// We assign the return value to a variable 'h' (handler) later
	// For now, we just ensure it connects.
	dbInstance := db.Init()

	// Check if dbInstance is valid (just for sanity)
	if dbInstance == nil {
		log.Fatal("Database instance is nil")
	}

	// Initialize the Handler with the DB connection
	authHandler := &internal.AuthHandler{
		DB: dbInstance,
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Auth Service",
			"status":  "active",
			"db":      "connected",
		})
	})

	// Register Routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// Run the server
	log.Println("Auth Service running on port 5001")
	r.Run(":5001")
}
