package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/models"
	"github.com/metabbe3/go-backend/repositories"
	"github.com/metabbe3/go-backend/utils"
)

type UserController struct {
	UserRepo repositories.UserRepositoryInterface
	Hasher   utils.PasswordHasher
}

// NewUserController returns a new instance of UserController
func NewUserController(userRepo repositories.UserRepositoryInterface, hasher utils.PasswordHasher) *UserController {
	return &UserController{UserRepo: userRepo, Hasher: hasher}
}

// CreateUser handles user creation
func (ctrl *UserController) CreateUser(c *gin.Context) {
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

	utils.SendCreated(c, "User created successfully", gin.H{"email": user.Email})
}

// GetUser handles getting user details by email
func (ctrl *UserController) GetUser(c *gin.Context) {
	userEmail := c.Param("email") // Assuming email is passed as a parameter in the URL

	user, err := ctrl.UserRepo.FindByEmail(userEmail) // Search by email
	if err != nil {
		utils.SendNotFound(c, "User not found")
		return
	}

	utils.SendSuccess(c, "User details fetched successfully", gin.H{"user": user})
}

// UpdateUser handles updating user details by email
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userEmail := c.Param("email")
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

	user, err := ctrl.UserRepo.FindByEmail(userEmail)
	if err != nil {
		utils.SendNotFound(c, "User not found")
		return
	}

	// Update user fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = hashedPassword
	}

	if err := ctrl.UserRepo.UpdateUser(user); err != nil {
		utils.SendInternalServerError(c, "Failed to update user")
		return
	}

	utils.SendSuccess(c, "User updated successfully", gin.H{"user": user})
}

// DeleteUser handles deleting a user by ID
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	// Get user ID from URL parameters
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32) // Assuming ID is passed as a parameter
	if err != nil {
		utils.SendBadRequest(c, "Invalid user ID")
		return
	}

	// Call the DeleteUser method from UserRepository by ID
	err = ctrl.UserRepo.DeleteUser(uint(userID))
	if err != nil {
		utils.SendNotFound(c, "User not found")
		return
	}

	// Send success response
	utils.SendSuccess(c, "User deleted successfully", nil)
}

// GetAllUsers handles fetching all users with pagination
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	users, totalCount, err := ctrl.UserRepo.GetAllUsers(limit, offset)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch users")
		return
	}

	utils.SendSuccess(c, "Users fetched successfully", gin.H{
		"data":        users,
		"total_count": totalCount,
	})
}
