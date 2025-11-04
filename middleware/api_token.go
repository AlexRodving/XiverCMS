package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

// APITokenMiddleware validates API tokens for programmatic access
func APITokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next() // Continue to next middleware (JWT auth)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.Next() // Not an API token, continue to JWT
			return
		}

		// Check if it's an API token (starts with xvc_)
		if !strings.HasPrefix(parts[1], "xvc_") {
			c.Next() // Not an API token, continue to JWT
			return
		}

		// Validate API token
		var token models.APIToken
		if err := database.DB.Where("token = ?", parts[1]).First(&token).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API token"})
			c.Abort()
			return
		}

		// Check if token is expired
		if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API token expired"})
			c.Abort()
			return
		}

		// Update last used timestamp
		now := time.Now()
		token.LastUsedAt = &now
		database.DB.Save(&token)

		// Set token info in context
		c.Set("apiTokenId", token.ID)
		c.Set("apiTokenType", token.Type)
		c.Set("apiTokenUserId", token.CreatedByID)

		// For read-only tokens, restrict to GET requests
		if token.Type == "read-only" && c.Request.Method != "GET" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Read-only token cannot perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
