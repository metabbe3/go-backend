package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/config"
	"github.com/metabbe3/go-backend/routes"
)

func runTests() {
	cmd := exec.Command("go", "test", "./controllers", "-cover")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		log.Fatalf("❌ Tests failed: %v", err)
	}
}

// InitializeApp sets up the application for running and testing.
func InitializeApp() (*gin.Engine, error) {
	fmt.Println("🚀 Starting application...")

	fmt.Println("🚀 Testing Functions...")
	runTests()

	// Initialize the database
	if err := config.ConnectDB(); err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		return nil, err
	}

	fmt.Println("✅ Database connected successfully!")

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	r := gin.Default()

	// Trust only localhost as a proxy (adjust if needed)
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Setup routes
	routes.SetupRoutes(r)

	return r, nil
}

func main() {
	r, err := InitializeApp()
	if err != nil {
		log.Fatalf("❌ Application initialization failed: %v", err)
	}

	// Start the server
	port := ":8080"
	fmt.Println("🌍 Server running on http://localhost" + port)
	if err := r.Run(port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
