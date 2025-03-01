package models

import "time"

// User represents a user in the system
type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	RegistrationNo string    `json:"registration_no"`
	PhoneNo        string    `json:"phone_no"`
	Date           time.Time `json:"date"`
}
