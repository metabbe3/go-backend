package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError structure for handling custom errors
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewError creates a new AppError
func NewError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// HandleError logs and responds with an error
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		// Custom application error
		log.Printf("App Error: %s", appErr.Error())
		c.JSON(appErr.Code, gin.H{
			"success": false,
			"message": appErr.Message,
			"code":    appErr.Code,
		})
	} else {
		// Generic internal server error
		log.Printf("Internal Server Error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"code":    http.StatusInternalServerError,
		})
	}
}
