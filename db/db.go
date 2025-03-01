package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB is the global variable that holds the database connection pool
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get database credentials from environment variables
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PASSWORD")
	host := os.Getenv("SQL_HOST")
	port := os.Getenv("SQL_PORT")
	database := os.Getenv("SQL_DB")

	// MySQL connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Set connection pool settings
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// Ping the database to check if it's reachable
	if err := DB.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}
