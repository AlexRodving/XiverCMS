package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xivercrm/xivercrm/auth"
	"github.com/xivercrm/xivercrm/config"
	"github.com/xivercrm/xivercrm/database"
	"github.com/xivercrm/xivercrm/models"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expiration, _ := time.ParseDuration(config.AppConfig.JWTExpiration)
	token, err := auth.GenerateToken(user.ID, user.Email, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Load roles
	database.DB.Model(&user).Association("Roles").Find(&user.Roles)

	// Create audit log
	CreateAuditLog(c, "login", "user", &user.ID, "User logged in", nil)

	c.JSON(http.StatusOK, gin.H{
		"jwt": token,
		"user": gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"username":     user.Username,
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"isSuperAdmin": user.IsSuperAdmin,
			"roles":        user.Roles,
		},
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Assign "Authenticated" role
	var authenticatedRole models.Role
	if err := database.DB.Where("name = ?", "Authenticated").First(&authenticatedRole).Error; err == nil {
		database.DB.Model(&user).Association("Roles").Append(&authenticatedRole)
	}

	expiration, _ := time.ParseDuration(config.AppConfig.JWTExpiration)
	token, err := auth.GenerateToken(user.ID, user.Email, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"jwt": token,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"username":  user.Username,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}

func Me(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.Preload("Roles").Preload("Roles.Permissions").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"username":     user.Username,
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"isSuperAdmin": user.IsSuperAdmin,
		"roles":        user.Roles,
	})
}
