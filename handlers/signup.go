package handlers

import (
	"net/http"

	"project/utils"

	"github.com/gin-gonic/gin"
)

// SignUpHandler registers a new user with Gmail, password, and date of birth.
func SignUpHandler(c *gin.Context) {
	var request struct {
		Username string `json:"username"` // Expected to be a Gmail address.
		Password string `json:"password"`
		DOB      string `json:"dob"` // Date of birth in format "YYYY-MM-DD"
	}

	// Bind the incoming JSON to the request struct.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if the user already exists.
	exists, err := utils.UserExists(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password using a secure hash function.
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
		return
	}

	// Create the new user record in the database.
	if err := utils.CreateUser(request.Username, hashedPassword, request.DOB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
