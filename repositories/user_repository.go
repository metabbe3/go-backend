package repositories

import (
	"github.com/metabbe3/go-backend/models"
	"gorm.io/gorm"
)

// UserRepositoryInterface defines the methods to interact with the User model
type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetAllUsers(limit, offset int) ([]models.User, int, error) // Updated
}

// UserRepository is a concrete implementation of the UserRepositoryInterface
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser saves a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user details (e.g., saving JWT token)
func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

// DeleteUser deletes a user by their ID
func (r *UserRepository) DeleteUser(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

// GetAllUsers retrieves all users from the database with limit, offset, and total count
func (r *UserRepository) GetAllUsers(limit, offset int) ([]models.User, int, error) {
	var users []models.User
	var totalCount int64

	// Fetch the users with limit and offset
	if err := r.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// Get the total count of users
	if err := r.DB.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	return users, int(totalCount), nil
}
