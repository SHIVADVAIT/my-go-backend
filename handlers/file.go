package handlers

import (
	"fmt"
	"log"
	"net/http"
	"project/db"
	"project/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

// CreateUser handles storing the user data in the temp table
func CreateUser(c *gin.Context) {
	var user models.User

	// Bind the JSON data to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Insert into temp table
	err := WriteTempData(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User stored in temp table successfully!"})
}

// WriteTempData inserts the user data into the temp table
func WriteTempData(user models.User) error {
	// Format the current date as "dd/mm/yy"
	currentDate := time.Now().Format("02/01/06")

	query := `
		INSERT INTO temp (name, email, registration_no, phone_no, date)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.DB.Exec(query, user.Name, user.Email, user.RegistrationNo, user.PhoneNo, currentDate)
	if err != nil {
		return fmt.Errorf("❌ Error inserting data into temp: %v", err)
	}

	return nil
}

// TransferTempData moves data from temp to users table
func TransferTempData() {

	// Begin a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("❌ Error starting transaction: %v", err)
		return
	}

	// Insert into users from temp
	insertQuery := `
		INSERT INTO users (name, email, registration_no, phone_no, date)
		SELECT name, email, registration_no, phone_no, date FROM temp
	`
	_, err = tx.Exec(insertQuery)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Error transferring data to users: %v", err)
		return
	}

	// Delete transferred records from temp
	deleteQuery := `DELETE FROM temp`
	_, err = tx.Exec(deleteQuery)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Error deleting temp data: %v", err)
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Transaction commit failed: %v", err)
		return
	}

}

// StartScheduler initializes and starts the scheduler for periodic tasks
func StartScheduler() {
	// Create a new scheduler instance
	scheduler := gocron.NewScheduler(time.Local)

	// Schedule the TransferTempData function to run every 10 seconds
	scheduler.Every(10).Seconds().Do(TransferTempData)

	// Start the scheduler in a separate goroutine
	go scheduler.StartAsync()
}
