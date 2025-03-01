package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	BaseURL = "https://custom-chatbot-api.p.rapidapi.com"
)

var ChatbotAPIKey string

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get API key from environment variable
	ChatbotAPIKey = os.Getenv("RAPIDAPI_CHATBOT_KEY")
	if ChatbotAPIKey == "" {
		fmt.Println("Warning: RAPIDAPI_CHATBOT_KEY is not set in .env")
	}
}
