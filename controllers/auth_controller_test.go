package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/models"
	"github.com/metabbe3/go-backend/test"
	"github.com/metabbe3/go-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthController_RegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    string
		mockSetup  func(mockRepo *test.MockUserRepository)
		expectCode int
		expectMsg  string
	}{
		{
			name:    "Success - Valid Registration",
			request: `{"email":"test@example.com","password":"StrongPass123"}`,
			mockSetup: func(mockRepo *test.MockUserRepository) {
				mockRepo.On("CreateUser", mock.Anything).Return(nil).Once()
			},
			expectCode: http.StatusCreated,
			expectMsg:  "User registered successfully",
		},
		{
			name:       "Failure - Invalid JSON",
			request:    `{"email":"test@example.com", "password":}`,
			mockSetup:  func(mockRepo *test.MockUserRepository) {},
			expectCode: http.StatusBadRequest,
			expectMsg:  "Invalid request data",
		},
		{
			name:       "Failure - Invalid Email",
			request:    `{"email":"invalid-email","password":"StrongPass123"}`,
			mockSetup:  func(mockRepo *test.MockUserRepository) {},
			expectCode: http.StatusBadRequest,
			expectMsg:  "Invalid request data", // ✅ Update expected message
		},
		{
			name:       "Failure - Weak Password",
			request:    `{"email":"test@example.com","password":"123"}`,
			mockSetup:  func(mockRepo *test.MockUserRepository) {},
			expectCode: http.StatusBadRequest,
			expectMsg:  "Invalid request data", // ✅ Update expected message
		},
		{
			name:    "Failure - Database Error",
			request: `{"email":"test@example.com","password":"StrongPass123"}`,
			mockSetup: func(mockRepo *test.MockUserRepository) {
				mockRepo.On("CreateUser", mock.Anything).Return(errors.New("db error")).Once()
			},
			expectCode: http.StatusInternalServerError,
			expectMsg:  "Failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(test.MockUserRepository) // ✅ Mock Repository
			ctrl := AuthController{
				UserRepo: mockRepo,             // ✅ Inject interface-based mock
				Hasher:   utils.BcryptHasher{}, // ✅ Ensure Hasher is set
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.request))
			c.Request.Header.Set("Content-Type", "application/json")

			tt.mockSetup(mockRepo) // ✅ Set up mocks

			ctrl.RegisterUser(c) // ✅ Call function

			assert.Equal(t, tt.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectMsg)
		})

	}
}

func TestAuthController_LoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    string
		mockSetup  func(mockRepo *test.MockUserRepository)
		expectCode int
		expectMsg  string
	}{
		{
			name:    "Success - Valid Login",
			request: `{"email":"test@example.com","password":"StrongPass123"}`,
			mockSetup: func(mockRepo *test.MockUserRepository) {
				hashedPassword, _ := utils.BcryptHasher{}.HashPassword("StrongPass123")
				mockRepo.On("FindByEmail", "test@example.com").Return(&models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
					Role:     "user",
				}, nil).Once()
				mockRepo.On("UpdateUser", mock.Anything).Return(nil).Once()
			},
			expectCode: http.StatusOK,
			expectMsg:  "Login successful",
		},
		{
			name:       "Failure - Invalid JSON",
			request:    `{"email":"test@example.com", "password":}`,
			mockSetup:  func(mockRepo *test.MockUserRepository) {},
			expectCode: http.StatusBadRequest,
			expectMsg:  "Invalid request data",
		},
		{
			name:       "Failure - Invalid Email",
			request:    `{"email":"invalid-email","password":"StrongPass123"}`,
			mockSetup:  func(mockRepo *test.MockUserRepository) {},
			expectCode: http.StatusBadRequest,
			expectMsg:  "Invalid request data",
		},
		{
			name:    "Failure - User Not Found",
			request: `{"email":"notfound@example.com","password":"StrongPass123"}`,
			mockSetup: func(mockRepo *test.MockUserRepository) {
				mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("user not found")).Once()
			},
			expectCode: http.StatusUnauthorized,
			expectMsg:  "Invalid credentials",
		},
		{
			name:    "Failure - Incorrect Password",
			request: `{"email":"test@example.com","password":"WrongPass123"}`,
			mockSetup: func(mockRepo *test.MockUserRepository) {
				hashedPassword, _ := utils.BcryptHasher{}.HashPassword("StrongPass123")
				mockRepo.On("FindByEmail", "test@example.com").Return(&models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
					Role:     "user",
				}, nil).Once()
			},
			expectCode: http.StatusUnauthorized,
			expectMsg:  "Invalid credentials",
		},
		{
			name:    "Failure - Database Error on Token Update",
			request: `{"email":"test@example.com","password":"StrongPass123"}`,
			mockSetup: func(mockRepo *test.MockUserRepository) {
				hashedPassword, _ := utils.BcryptHasher{}.HashPassword("StrongPass123")
				mockRepo.On("FindByEmail", "test@example.com").Return(&models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
					Role:     "user",
				}, nil).Once()
				mockRepo.On("UpdateUser", mock.Anything).Return(errors.New("db error")).Once()
			},
			expectCode: http.StatusInternalServerError,
			expectMsg:  "Failed to update user token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(test.MockUserRepository)

			ctrl := AuthController{
				UserRepo: mockRepo,
				Hasher:   utils.BcryptHasher{},
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.request))
			c.Request.Header.Set("Content-Type", "application/json")

			tt.mockSetup(mockRepo)

			ctrl.LoginUser(c)

			assert.Equal(t, tt.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectMsg)
		})
	}
}

// func TestAuthController_LogoutUser(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	tests := []struct {
// 		name       string
// 		authHeader string
// 		mockSetup  func(mockRepo *test.MockUserRepository)
// 		expectCode int
// 		expectMsg  string
// 	}{
// 		{
// 			name:       "Success - Valid Token",
// 			authHeader: "Bearer valid_token",
// 			mockSetup: func(mockRepo *test.MockUserRepository) {
// 				mockRepo.On("FindByEmail", "test@example.com").Return(&models.User{
// 					ID:    1,
// 					Email: "test@example.com",
// 					Token: "valid_token",
// 				}, nil).Once()
// 				mockRepo.On("UpdateUser", mock.Anything).Return(nil).Once()

// 				// Mock token validation
// 				utils.ValidateToken = func(token string) (*utils.Claims, error) {
// 					if token == "valid_token" {
// 						return &utils.Claims{Username: "test@example.com"}, nil
// 					}
// 					return nil, errors.New("invalid token")
// 				}
// 			},
// 			expectCode: http.StatusOK,
// 			expectMsg:  "Logout successful",
// 		},
// 		{
// 			name:       "Failure - No Token Provided",
// 			authHeader: "",
// 			mockSetup:  func(mockRepo *test.MockUserRepository) {},
// 			expectCode: http.StatusUnauthorized,
// 			expectMsg:  "No token provided",
// 		},
// 		{
// 			name:       "Failure - Invalid Token",
// 			authHeader: "Bearer invalid_token",
// 			mockSetup: func(mockRepo *test.MockUserRepository) {
// 				utils.ValidateToken = func(token string) (*utils.Claims, error) {
// 					return nil, errors.New("invalid token")
// 				}
// 			},
// 			expectCode: http.StatusUnauthorized,
// 			expectMsg:  "Invalid token",
// 		},
// 		{
// 			name:       "Failure - User Not Found",
// 			authHeader: "Bearer valid_token",
// 			mockSetup: func(mockRepo *test.MockUserRepository) {
// 				mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("user not found")).Once()

// 				utils.ValidateToken = func(token string) (*utils.Claims, error) {
// 					if token == "valid_token" {
// 						return &utils.Claims{Username: "test@example.com"}, nil
// 					}
// 					return nil, errors.New("invalid token")
// 				}
// 			},
// 			expectCode: http.StatusUnauthorized,
// 			expectMsg:  "User not found",
// 		},
// 		{
// 			name:       "Failure - Database Error on Logout",
// 			authHeader: "Bearer valid_token",
// 			mockSetup: func(mockRepo *test.MockUserRepository) {
// 				mockRepo.On("FindByEmail", "test@example.com").Return(&models.User{
// 					ID:    1,
// 					Email: "test@example.com",
// 					Token: "valid_token",
// 				}, nil).Once()
// 				mockRepo.On("UpdateUser", mock.Anything).Return(errors.New("db error")).Once()

// 				utils.ValidateToken = func(token string) (*utils.Claims, error) {
// 					if token == "valid_token" {
// 						return &utils.Claims{Username: "test@example.com"}, nil
// 					}
// 					return nil, errors.New("invalid token")
// 				}
// 			},
// 			expectCode: http.StatusInternalServerError,
// 			expectMsg:  "Failed to logout",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(test.MockUserRepository)

// 			ctrl := AuthController{
// 				UserRepo: mockRepo,
// 				Hasher:   utils.BcryptHasher{},
// 			}

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)
// 			c.Request = httptest.NewRequest(http.MethodPost, "/logout", nil)
// 			c.Request.Header.Set("Authorization", tt.authHeader)

// 			tt.mockSetup(mockRepo)

// 			ctrl.LogoutUser(c)

// 			assert.Equal(t, tt.expectCode, w.Code)
// 			assert.Contains(t, w.Body.String(), tt.expectMsg)
// 		})
// 	}
// }
