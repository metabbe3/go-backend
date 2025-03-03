package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/metabbe3/go-backend/models"
	"github.com/metabbe3/go-backend/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Load environment variables
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		utils.Warning("No .env file found, using default values")
	}
}

// GetEnv gets environment variables with a default fallback
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ConnectDB initializes the database connection
func ConnectDB() error { // 🔹 Change function to return error
	utils.InitLogger()
	loadEnv()

	utils.Info("Initializing database connection...")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		GetEnv("DB_USER", "root"),
		GetEnv("DB_PASSWORD", "secret"),
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "3306"),
		GetEnv("DB_NAME", "mydatabase"),
		GetEnv("DB_CHARSET", "utf8mb4"),
		GetEnv("DB_PARSE_TIME", "True"),
		GetEnv("DB_LOC", "Local"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Error(fmt.Sprintf("Database connection failed: %v", err))
		return err // 🔹 Return the error instead of logging fatal
	}

	utils.Info("Connected to MySQL database successfully!")
	autoMigrate()
	return nil // ✅ Return nil on success
}

// autoMigrate runs migrations for models
func autoMigrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Customer{}) // Add more models as needed
	if err != nil {
		utils.Error(fmt.Sprintf("Auto migration failed: %v", err))
		log.Fatalf("Failed to migrate database: %v", err)
	}
	utils.Info("Database migration completed!")
}
