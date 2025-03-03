package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/utils"
)

// JWTAuthMiddleware protects routes that require authentication
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Log request details
		utils.Info("JWT Middleware: Checking Authorization header")

		if authHeader == "" {
			utils.Warning("JWT Middleware: Missing Authorization header")
			utils.SendUnauthorized(c, "Authorization header is missing")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.Warning("JWT Middleware: Invalid authorization format")
			utils.SendUnauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.Error(fmt.Sprintf("JWT Middleware: Invalid or expired token - %v", err))
			utils.SendUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Attach user claims to the context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		utils.Info(fmt.Sprintf("JWT Middleware: Authentication successful for userID: %d", claims.UserID))
		c.Next()
	}
}
