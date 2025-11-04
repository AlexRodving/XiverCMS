package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

func GetAPITokens(c *gin.Context) {
	var tokens []models.APIToken
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.APIToken{})

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mask tokens in response (show only first 8 characters)
	for i := range tokens {
		if len(tokens[i].Token) > 8 {
			tokens[i].Token = tokens[i].Token[:8] + "..."
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tokens,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

func GetAPIToken(c *gin.Context) {
	id := c.Param("id")
	var token models.APIToken

	if err := database.DB.Preload("CreatedBy").First(&token, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API token not found"})
		return
	}

	// Mask token
	if len(token.Token) > 8 {
		token.Token = token.Token[:8] + "..."
	}

	c.JSON(http.StatusOK, token)
}

type CreateAPITokenRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	ExpiresAt   *time.Time `json:"expiresAt"`
}

func CreateAPIToken(c *gin.Context) {
	var req CreateAPITokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	tokenString := "xvc_" + hex.EncodeToString(tokenBytes)

	userId, _ := c.Get("userId")
	userID := userId.(uint)

	if req.Type == "" {
		req.Type = "read-only"
	}

	apiToken := models.APIToken{
		Name:        req.Name,
		Token:       tokenString,
		Description: req.Description,
		Type:        req.Type,
		ExpiresAt:   req.ExpiresAt,
		CreatedByID: &userID,
	}

	if err := database.DB.Create(&apiToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return full token only on creation
	response := apiToken
	c.JSON(http.StatusCreated, response)
}

func UpdateAPIToken(c *gin.Context) {
	id := c.Param("id")
	var token models.APIToken

	if err := database.DB.First(&token, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API token not found"})
		return
	}

	var req CreateAPITokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		token.Name = req.Name
	}
	if req.Description != "" {
		token.Description = req.Description
	}
	if req.Type != "" {
		token.Type = req.Type
	}
	if req.ExpiresAt != nil {
		token.ExpiresAt = req.ExpiresAt
	}

	if err := database.DB.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mask token
	if len(token.Token) > 8 {
		token.Token = token.Token[:8] + "..."
	}

	c.JSON(http.StatusOK, token)
}

func DeleteAPIToken(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.APIToken{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API token deleted successfully"})
}
