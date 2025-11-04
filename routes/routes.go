package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xivercrm/xivercrm/handlers"
	"github.com/xivercrm/xivercrm/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve uploaded media files
	r.Static("/api/uploads", "./uploads")

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/register", handlers.Register)
	}

	// Protected routes - supports both JWT and API tokens
	protected := r.Group("/api")
	protected.Use(middleware.APITokenMiddleware()) // Try API token first
	protected.Use(middleware.AuthMiddleware())     // Then try JWT
	{
		// Auth
		protected.GET("/auth/me", handlers.Me)

		// Users
		users := protected.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
			users.PUT("/:id/password", handlers.ChangePassword)
		}

		// Roles & Permissions
		roles := protected.Group("/roles")
		{
			roles.GET("", handlers.GetRoles)
			roles.GET("/:id", handlers.GetRole)
			roles.POST("", handlers.CreateRole)
			roles.PUT("/:id", handlers.UpdateRole)
			roles.DELETE("/:id", handlers.DeleteRole)
		}

		permissions := protected.Group("/permissions")
		{
			permissions.GET("", handlers.GetPermissions)
			permissions.GET("/:id", handlers.GetPermission)
			permissions.POST("", handlers.CreatePermission)
			permissions.PUT("/:id", handlers.UpdatePermission)
			permissions.DELETE("/:id", handlers.DeletePermission)
		}

		// API Tokens
		apiTokens := protected.Group("/api-tokens")
		{
			apiTokens.GET("", handlers.GetAPITokens)
			apiTokens.GET("/:id", handlers.GetAPIToken)
			apiTokens.POST("", handlers.CreateAPIToken)
			apiTokens.PUT("/:id", handlers.UpdateAPIToken)
			apiTokens.DELETE("/:id", handlers.DeleteAPIToken)
		}

		// Media Library
		media := protected.Group("/upload")
		{
			media.POST("", handlers.UploadMedia)
		}

		mediaFiles := protected.Group("/media-files")
		{
			mediaFiles.GET("", handlers.GetMediaFiles)
			mediaFiles.GET("/:id", handlers.GetMediaFile)
			mediaFiles.PUT("/:id", handlers.UpdateMediaFile)
			mediaFiles.DELETE("/:id", handlers.DeleteMediaFile)
		}

		// Content Types
		contentTypes := protected.Group("/content-types")
		{
			contentTypes.GET("", handlers.GetContentTypes)
			contentTypes.GET("/:uid", handlers.GetContentType)
			contentTypes.POST("", handlers.CreateContentType)
			contentTypes.PUT("/:uid", handlers.UpdateContentType)
			contentTypes.DELETE("/:uid", handlers.DeleteContentType)
		}

		// Content Entries
		contentEntries := protected.Group("/content-types/:uid/entries")
		{
			contentEntries.GET("", handlers.GetContentEntries)
			contentEntries.GET("/:id", handlers.GetContentEntry)
			contentEntries.POST("", handlers.CreateContentEntry)
			contentEntries.PUT("/:id", handlers.UpdateContentEntry)
			contentEntries.DELETE("/:id", handlers.DeleteContentEntry)
			contentEntries.GET("/:id/history", handlers.GetContentHistory)

			// Relations
			contentEntries.GET("/:id/relations", handlers.GetRelations)
			contentEntries.POST("/:id/relations", handlers.CreateRelation)
			contentEntries.DELETE("/:id/relations/:relationId", handlers.DeleteRelation)
			contentEntries.GET("/:id/relations/:field", handlers.GetRelatedEntries)
		}

		// Audit Logs
		auditLogs := protected.Group("/audit-logs")
		{
			auditLogs.GET("", handlers.GetAuditLogs)
		}
	}
}
