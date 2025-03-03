package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/config"
	"github.com/metabbe3/go-backend/controllers"
	"github.com/metabbe3/go-backend/middleware"
	"github.com/metabbe3/go-backend/repositories"
	"github.com/metabbe3/go-backend/utils" // Import PasswordHasher package
)

// SetupRoutes initializes all routes
func SetupRoutes(router *gin.Engine) {
	// Ensure DB is initialized
	if config.DB == nil {
		panic("Database connection is not initialized")
	}

	// Initialize repositories with the global DB instance
	userRepo := repositories.NewUserRepository(config.DB)
	customerRepo := repositories.NewCustomerRepository(config.DB)

	// Initialize controllers with repositories and utils
	authController := controllers.NewAuthController(userRepo, utils.BcryptHasher{})
	userController := controllers.NewUserController(userRepo, utils.BcryptHasher{})
	customerController := controllers.NewCustomerController(customerRepo)

	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", authController.LoginUser)       // Login user
		auth.POST("/register", authController.RegisterUser) // Register new user
	}

	// Protected API routes (JWT required)
	api := router.Group("/api")
	api.Use(middleware.JWTAuthMiddleware()) // Apply JWT middleware to the API group
	{
		// User routes
		api.POST("/user", userController.CreateUser)       // Create user
		api.GET("/user/:email", userController.GetUser)    // Get user by ID
		api.PUT("/user/:email", userController.UpdateUser) // Update user by ID
		api.DELETE("/user/:id", userController.DeleteUser) // Delete user by ID
		api.GET("/users", userController.GetAllUsers)      // Get all users

		// Customer routes
		api.POST("/customer", customerController.CreateCustomer)       // Create customer
		api.GET("/customer/:id", customerController.GetCustomer)       // Get customer by ID
		api.PUT("/customer/:id", customerController.UpdateCustomer)    // Update customer by ID
		api.DELETE("/customer/:id", customerController.DeleteCustomer) // Delete customer by ID
		api.GET("/customers", customerController.GetAllCustomers)      // Get all customers

		// Dashboard route
		api.GET("/dashboard", func(c *gin.Context) {
			// This is just an example; you can add logic to return a dashboard summary
			c.JSON(200, gin.H{
				"message":   "Welcome to the Dashboard",
				"users":     "/api/users",
				"customers": "/api/customers",
			})
		})
	}

	// Print routes for debugging (optional)
	router.GET("/routes", func(c *gin.Context) {
		var routes []string
		for _, r := range router.Routes() {
			routes = append(routes, r.Method+" "+r.Path)
		}
		c.JSON(200, gin.H{"routes": routes})
	})
}
