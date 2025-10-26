package middleware

import (
	"net/http"

	"matchee/services/internal/entity"
	"matchee/services/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	jwtService service.JWTService
}

func NewAdminMiddleware(jwtService service.JWTService) *AdminMiddleware {
	return &AdminMiddleware{jwtService: jwtService}
}

func (m *AdminMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("Authorization header required"))
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("Invalid token"))
			c.Abort()
			return
		}

		// Check if user has admin role
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, entity.ErrorResponse(http.StatusForbidden, "Admin access required"))
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
