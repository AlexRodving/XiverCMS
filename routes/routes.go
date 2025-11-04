package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/handlers"
	"github.com/xivercms/xivercms/middleware"
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
		// Auth routes (must be before dynamic content routes)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/register", handlers.Register)

		// Public roles for registration (must be before dynamic content routes)
		public.GET("/roles/public", handlers.GetPublicRoles)

		// Public content access - uses OptionalAuthMiddleware to check auth if provided
		// Access is controlled by accessType in ContentType (public, authenticated, moderator, admin)
		// Simplified URLs like Strapi: /api/{content-type} and /api/{content-type}/{id}
		publicContent := public.Group("")
		publicContent.Use(middleware.OptionalAuthMiddleware())
		{
			// Get all content types (list of available types)
			// Must be registered before /:uid to avoid conflicts
			publicContent.GET("/content-types", handlers.PublicGetContentTypes)

			// Get specific content type schema
			// Must be registered before /:uid to avoid conflicts
			publicContent.GET("/content-types/:uid", handlers.PublicGetContentType)

			// Public API: Get single entry by ID (must be before /:uid)
			// URL: /api/{uid}/{id} (e.g., /api/articles/1, /api/books/123)
			publicContent.GET("/:uid/:id", handlers.PublicGetContentEntry)

			// Public API: Get all entries for a content type (registered last to catch remaining routes)
			// URL: /api/{uid} (e.g., /api/articles, /api/books)
			// This will match any /api/{something} that isn't matched above
			publicContent.GET("/:uid", handlers.PublicGetContentEntries)
		}
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

		// Content Entries Management (protected - requires auth)
		// NOTE: Public read access to published entries is via /api/content-types/:uid/entries (public routes above)
		// These endpoints allow managing entries (create, update, delete)
		// For reading all entries including drafts, use admin endpoints below
		contentEntries := protected.Group("/admin/content-types/:uid/entries")
		{
			// Admin endpoints for reading all entries (including drafts)
			contentEntries.GET("", handlers.GetContentEntries)   // Can get all statuses with auth
			contentEntries.GET("/:id", handlers.GetContentEntry) // Can get any entry with auth

			// Management endpoints
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
