package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

// AccessLevelMiddleware checks if user has required access level
// accessType: public, authenticated, moderator, admin
func AccessLevelMiddleware(accessType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Public access - no authentication required
		if accessType == "public" {
			c.Next()
			return
		}

		// Check if user is authenticated
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Authenticated access - any logged in user
		if accessType == "authenticated" {
			c.Next()
			return
		}

		// Load user with roles
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, userId).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Super admin has all access
		if user.IsSuperAdmin {
			c.Next()
			return
		}

		// Check if user has required role
		hasAccess := false
		for _, role := range user.Roles {
			// Moderator access - user has "Moderator" role
			if accessType == "moderator" {
				if role.Name == "Moderator" || role.Name == "Admin" {
					hasAccess = true
					break
				}
			}
			// Admin access - user has "Admin" role
			if accessType == "admin" {
				if role.Name == "Admin" {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient access level"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckContentTypeAccess checks if user can access a specific content type
func CheckContentTypeAccess(contentTypeUID string, c *gin.Context) bool {
	var contentType models.ContentType
	if err := database.DB.Where("uid = ?", contentTypeUID).First(&contentType).Error; err != nil {
		return false
	}

	// Public access
	if contentType.AccessType == "public" {
		return true
	}

	// Check authentication
	userId, exists := c.Get("userId")
	if !exists {
		return false
	}

	// Authenticated access
	if contentType.AccessType == "authenticated" {
		return true
	}

	// Load user with roles
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userId).Error; err != nil {
		return false
	}

	// Super admin
	if user.IsSuperAdmin {
		return true
	}

	// Check role-based access
	for _, role := range user.Roles {
		if contentType.AccessType == "moderator" {
			if role.Name == "Moderator" || role.Name == "Admin" {
				return true
			}
		}
		if contentType.AccessType == "admin" {
			if role.Name == "Admin" {
				return true
			}
		}
	}

	return false
}
