package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/pkg/jwt"
	"github.com/zaem/backend-go/pkg/response"
)

// Protect verifies JWT and loads user into context
func Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, 401, "Not authorized. Please login.")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.Verify(token)
		if err != nil {
			response.Error(c, 401, "Invalid or expired token. Please login again.")
			c.Abort()
			return
		}

		var user models.User
		err = database.DB.
			Select("id", "name", "email", "phone", "role", "is_active").
			First(&user, "id = ?", claims.UserID).Error
		if err != nil {
			response.Error(c, 401, "User no longer exists.")
			c.Abort()
			return
		}

		if !user.IsActive {
			response.Error(c, 403, "Your account has been disabled.")
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Set("user_role", string(user.Role))
		c.Next()
	}
}

// AdminOnly allows only ADMIN role
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Error(c, 403, "Admin access required.")
			c.Abort()
			return
		}

		roleStr, _ := role.(string)
		if roleStr != string(models.RoleAdmin) {
			response.Error(c, 403, "Admin access required.")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ManagerOrAdmin allows MANAGER or ADMIN
func ManagerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		roleStr, _ := role.(string)
		if roleStr != string(models.RoleAdmin) && roleStr != string(models.RoleManager) {
			response.Error(c, 403, "Access required.")
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalAuth — doesn't fail if no token, just sets user if valid
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.Verify(token)
		if err != nil {
			c.Next()
			return
		}

		var user models.User
		err = database.DB.
			Select("id", "name", "email", "phone", "role", "is_active").
			First(&user, "id = ?", claims.UserID).Error
		if err == nil && user.IsActive {
			c.Set("user", &user)
			c.Set("user_id", user.ID)
			c.Set("user_role", string(user.Role))
		}
		c.Next()
	}
}