package utils

import (
	"database/sql"
	"fmt"

	"project/db" // Import your db package which initializes the DB connection

	"golang.org/x/crypto/bcrypt"
)

// User represents a record in the signupusers table.
type User struct {
	Username       string
	HashedPassword string
	DOB            string
}

// HashPassword generates a bcrypt hash of the provided password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with a hashed password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser inserts a new record into the signupusers table.
func CreateUser(username, hashedPassword, dob string) error {
	query := "INSERT INTO signupusers (username, password, dob) VALUES (?, ?, ?)"
	_, err := db.DB.Exec(query, username, hashedPassword, dob)
	return err
}

// GetUserByUsername retrieves a record from the signupusers table by username.
func GetUserByUsername(username string) (*User, error) {
	query := "SELECT username, password, dob FROM signupusers WHERE username = ?"
	row := db.DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.Username, &user.HashedPassword, &user.DOB)
	if err == sql.ErrNoRows {
		return nil, nil // No record found
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

// UserExists checks if a record with the given username exists in the signupusers table.
func UserExists(username string) (bool, error) {
	query := "SELECT COUNT(*) FROM signupusers WHERE username = ?"
	var count int
	err := db.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		// Log the detailed error
		fmt.Printf("Error in UserExists query: %v\n", err)
		return false, err
	}
	return count > 0, nil
}
