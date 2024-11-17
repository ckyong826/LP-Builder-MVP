package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASS", ""),
        DBName:     getEnv("DB_NAME", "postgres"),
        Port:       getEnv("PORT", "8080"),
    }

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Kuala_Lumpur",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
    )
}

// InitDB initializes the database connection
func InitDB(cfg *Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    log.Println("Database connection established")
    return db, nil
}