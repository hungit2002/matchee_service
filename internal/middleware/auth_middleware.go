package middleware

import (
	"net/http"
	"strings"

	"matchee/services/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService service.JWTService
}

func NewAuthMiddleware(jwtService service.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userPhone", claims.Phone)
		c.Set("userFullName", claims.FullName)

		c.Next()
	}
}

// OptionalAuth middleware validates JWT token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userPhone", claims.Phone)
		c.Set("userFullName", claims.FullName)

		c.Next()
	}
}

// GetCurrentUserID extracts user ID from context
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint64)
	return id, ok
}

// GetCurrentUserEmail extracts user email from context
func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetCurrentUserPhone extracts user phone from context
func GetCurrentUserPhone(c *gin.Context) (string, bool) {
	phone, exists := c.Get("userPhone")
	if !exists {
		return "", false
	}

	phoneStr, ok := phone.(string)
	return phoneStr, ok
}

// GetCurrentUserFullName extracts user full name from context
func GetCurrentUserFullName(c *gin.Context) (string, bool) {
	fullName, exists := c.Get("userFullName")
	if !exists {
		return "", false
	}

	fullNameStr, ok := fullName.(string)
	return fullNameStr, ok
}
