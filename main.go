package main

import (
	"log"
	"project/db"
	"project/handlers"
	"project/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize the database
	db.InitDB()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve the secret key

	// Start the scheduler for transferring data from temp to users
	go handlers.StartScheduler() // This will run in the background
	// Setup Gin router
	r := routes.SetupRouter()

	// Enable CORS
	r.Use(cors.Default())

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
