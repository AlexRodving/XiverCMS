package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xivercrm/xivercrm/auth"
	"github.com/xivercrm/xivercrm/config"
	"github.com/xivercrm/xivercrm/database"
	"github.com/xivercrm/xivercrm/middleware"
	"github.com/xivercrm/xivercrm/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize auth
	auth.InitAuth(config.AppConfig.JWTSecret)

	// Connect to database
	database.Connect()

	// Run migrations
	database.Migrate()

	// Seed initial data
	database.Seed()

	// Setup Gin router
	gin.SetMode(config.AppConfig.GinMode)
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := ":" + config.GetPort()
	log.Printf("Server starting on port %s", port)

	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
