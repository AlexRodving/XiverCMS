package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

// RBACMiddleware checks if user has permission to perform action
func RBACMiddleware(action, subject string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Preload("Roles").Preload("Roles.Permissions").First(&user, userId).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Super admin has all permissions
		if user.IsSuperAdmin {
			c.Next()
			return
		}

		// Check if user has permission through roles
		hasPermission := false
		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				// Check if permission matches
				if permissionMatches(permission, action, subject, c) {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func permissionMatches(permission models.Permission, action, subject string, c *gin.Context) bool {
	// Check action (wildcard support)
	if permission.Action != "all" && permission.Action != action {
		return false
	}

	// Check subject (wildcard support)
	if permission.Subject != "all" && permission.Subject != subject {
		// Check if subject matches pattern (e.g., "content-type:article")
		if !strings.HasPrefix(permission.Subject, subject+":") {
			return false
		}
	}

	// TODO: Check conditions and properties from JSON fields
	// This would require parsing the JSON and evaluating conditions

	return true
}

// RequirePermission is a helper to create middleware with specific action and subject
func RequirePermission(action, subject string) gin.HandlerFunc {
	return RBACMiddleware(action, subject)
}
