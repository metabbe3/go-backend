package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response format structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Code:    http.StatusOK,
		Data:    data,
	})
}

// SendCreated sends a success response with HTTP 201
func SendCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Code:    http.StatusCreated,
		Data:    data,
	})
}

// SendError sends an error response with a given status code
func SendError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Code:    statusCode,
	})
}

// SendValidationError sends a validation error response
func SendValidationError(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": message,
		"code":    http.StatusBadRequest,
		"errors":  errors,
	})
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
		Code:    http.StatusUnauthorized,
	})
}

// SendForbidden sends a 403 Forbidden response
func SendForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
		Code:    http.StatusForbidden,
	})
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
		Code:    http.StatusNotFound,
	})
}

// SendInternalServerError sends a 500 Internal Server Error response
func SendInternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Code:    http.StatusInternalServerError,
	})
}

// SendBadRequest sends a 400 Bad Request response
func SendBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Code:    http.StatusBadRequest,
	})
}
