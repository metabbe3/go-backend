package test

import (
	"github.com/metabbe3/go-backend/models"
	"github.com/metabbe3/go-backend/repositories"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepositoryInterface
var _ repositories.UserRepositoryInterface = (*MockUserRepository)(nil)

// CreateUser mocks the CreateUser function
func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindByEmail mocks the FindByEmail function
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	// If the first return value is nil, return the error as the second return value
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// UpdateUser mocks the UpdateUser function
func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser function
func (m *MockUserRepository) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetAllUsers mocks the GetAllUsers function
func (m *MockUserRepository) GetAllUsers(limit, offset int) ([]models.User, int, error) {
	args := m.Called(limit, offset)
	// If the first return value is nil, return an empty slice and the error
	if args.Get(0) == nil {
		return nil, 0, args.Error(1)
	}
	// Return the users, total count and error if applicable
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}
