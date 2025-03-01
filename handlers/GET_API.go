package handlers

import (
	"fmt"
	"net/http"
	"project/db" // Importing the db package for accessing the DB connection
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllUsers handles the retrieval of all users from the 'users' table
func GetAllUsers(c *gin.Context) {
	// Call FetchAllUsers to retrieve all users data from the database
	users, err := FetchAllUsers()
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error fetching users: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the users data
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// FetchAllUsers retrieves all user data from the 'users' table
func FetchAllUsers() ([]map[string]interface{}, error) {
	// Prepare the SQL query to fetch all user data
	query := `SELECT id, name, email, registration_no, phone_no, date FROM users`

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error fetching users from database: %v\n", err)
		return nil, fmt.Errorf("error fetching users data: %v", err)
	}
	defer rows.Close()

	// Slice to store all users data
	var users []map[string]interface{}

	// Loop through the rows and scan the data into the user slice
	for rows.Next() {
		var id int
		var name, email, registrationNo, phoneNo, date string
		if err := rows.Scan(&id, &name, &email, &registrationNo, &phoneNo, &date); err != nil {
			// Log the error for debugging
			fmt.Printf("Error scanning user data: %v\n", err)
			return nil, fmt.Errorf("error scanning user data: %v", err)
		}

		// Parse the date from the database, format it to dd/mm/yy
		parsedDate, err := time.Parse("02/01/06", date)
		if err != nil {
			// If date parsing fails, return an error
			fmt.Printf("Error parsing date: %v\n", err)
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		formattedDate := parsedDate.Format("02/01/06") // Format to dd/mm/yy

		// Prepare the user map with formatted date
		user := map[string]interface{}{
			"id":              id,
			"name":            name,
			"email":           email,
			"registration_no": registrationNo,
			"phone_no":        phoneNo,
			"date":            formattedDate, // Use the formatted date
		}
		users = append(users, user)
	}

	// Check for any row iteration error
	if err := rows.Err(); err != nil {
		// Log the error for debugging
		fmt.Printf("Error iterating over rows: %v\n", err)
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	// Log the users data before returning it
	fmt.Printf("Users fetched: %+v\n", users)

	return users, nil
}

// GetUserByID handles the retrieval of a user by their ID
func GetUserByID(c *gin.Context) {
	// Get the user ID from the URL parameter
	userID := c.Param("id")

	// Call FetchUser to retrieve the user data from the database
	user, err := FetchUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the user data
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// FetchUserByID retrieves user data by ID from the 'users' table
func FetchUserByID(userID string) (map[string]interface{}, error) {
	// Prepare the SQL query to fetch user data by ID
	query := `SELECT name, email, registration_no, phone_no, date FROM users WHERE id = ?`

	// Execute the query with the provided user ID
	row := db.DB.QueryRow(query, userID)

	// Map to store the fetched user data
	var name, email, registrationNo, phoneNo, date string
	if err := row.Scan(&name, &email, &registrationNo, &phoneNo, &date); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user data: %v", err)
	}

	// Parse and format the date into dd/mm/yy
	parsedDate, err := time.Parse("02/01/06", date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %v", err)
	}
	formattedDate := parsedDate.Format("02/01/06")

	// Return the user data as a map
	user := map[string]interface{}{
		"id":              userID,
		"name":            name,
		"email":           email,
		"registration_no": registrationNo,
		"phone_no":        phoneNo,
		"date":            formattedDate,
	}
	return user, nil
}
