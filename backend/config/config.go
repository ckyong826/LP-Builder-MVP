package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadConfig() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    log.Println("Environment variables loaded successfully")
}

// DB is the global database connection variable.
var DB *gorm.DB

// Connect initializes the database connection.
func InitDB() {
	// Read environment variables for database configuration
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")

	// Create the Data Source Name (DSN) for the database connection
	dsn := "host=" + DB_HOST + " user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " port=" + DB_PORT + " sslmode=disable TimeZone=Asia/Kuala_Lumpur"

	log.Println(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	// Attempt to connect to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("Database connection established")
}
