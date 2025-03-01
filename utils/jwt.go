package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT generates a JWT token for a given username.
func GenerateJWT(username string) (string, error) {
	// Get the secret key from environment variables.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in environment variables")
	}

	// Create a new token object with signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // token expires in 24 hours
	})

	// Sign and get the complete encoded token as a string.
	return token.SignedString([]byte(jwtSecret))
}
