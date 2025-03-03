package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/models"
	"github.com/metabbe3/go-backend/repositories"
	"github.com/metabbe3/go-backend/utils"
)

// Change from `*repositories.UserRepository` to `repositories.UserRepositoryInterface`
type AuthController struct {
	UserRepo repositories.UserRepositoryInterface
	Hasher   utils.PasswordHasher // Use an interface instead of direct utils.HashPassword call
}

// Change `*repositories.UserRepository` to `repositories.UserRepositoryInterface`
func NewAuthController(userRepo repositories.UserRepositoryInterface, hasher utils.PasswordHasher) *AuthController {
	return &AuthController{UserRepo: userRepo, Hasher: hasher}
}

// RegisterUser handles user registration
func (ctrl *AuthController) RegisterUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, "Invalid request data", err.Error())
		return
	}

	hashedPassword, err := ctrl.Hasher.HashPassword(req.Password)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to hash password")
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := ctrl.UserRepo.CreateUser(&user); err != nil {
		utils.SendInternalServerError(c, "Failed to create user")
		return
	}

	utils.SendCreated(c, "User registered successfully", gin.H{"email": user.Email})
}

// LoginUser handles user authentication
func (ctrl *AuthController) LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, "Invalid request data", err.Error())
		return
	}

	// Validate input
	if err := utils.ValidateUserInput(req.Email, req.Password); err != nil {
		utils.SendValidationError(c, "Validation error", err.Error())
		return
	}

	// Find user
	user, err := ctrl.UserRepo.FindByEmail(req.Email)
	if err != nil {
		utils.SendUnauthorized(c, "Invalid credentials")
		return
	}

	// Validate password
	if err := ctrl.Hasher.ComparePasswords(user.Password, req.Password); err != nil {
		utils.SendUnauthorized(c, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to generate token")
		return
	}

	// Save token in DB
	user.Token = token
	if err := ctrl.UserRepo.UpdateUser(user); err != nil { // ✅ Remove '&' since user is already a pointer
		utils.SendInternalServerError(c, "Failed to update user token")
		return
	}

	utils.SendSuccess(c, "Login successful", gin.H{"token": token})
}

// LogoutUser handles user logout by removing the JWT token from the database
func (ctrl *AuthController) LogoutUser(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.SendUnauthorized(c, "No token provided")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateToken(token)
	if err != nil {
		utils.SendUnauthorized(c, "Invalid token")
		return
	}

	// Find user by email
	user, err := ctrl.UserRepo.FindByEmail(claims.Username)
	if err != nil {
		utils.SendUnauthorized(c, "User not found")
		return
	}

	// Clear token from DB
	user.Token = ""
	if err := ctrl.UserRepo.UpdateUser(user); err != nil {
		utils.SendInternalServerError(c, "Failed to logout")
		return
	}

	utils.SendSuccess(c, "Logout successful", nil)
}
