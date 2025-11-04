package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

func GetRoles(c *gin.Context) {
	var roles []models.Role
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Role{})

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": roles,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

// GetPublicRoles - get public roles available for registration
func GetPublicRoles(c *gin.Context) {
	var roles []models.Role

	if err := database.DB.Where("type = ?", "public").
		Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func GetRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.Preload("Permissions").Preload("Users").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	PermissionIDs []uint `json:"permissionIds"`
}

func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role already exists
	var existing models.Role
	if err := database.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role with this name already exists"})
		return
	}

	if req.Type == "" {
		req.Type = "custom"
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		var permissions []models.Permission
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
		database.DB.Model(&role).Association("Permissions").Append(permissions)
	}

	database.DB.Preload("Permissions").First(&role, role.ID)

	c.JSON(http.StatusCreated, role)
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Type != "" {
		role.Type = req.Type
	}

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update permissions if provided
	if req.PermissionIDs != nil {
		var permissions []models.Permission
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
		database.DB.Model(&role).Association("Permissions").Replace(permissions)
	}

	database.DB.Preload("Permissions").First(&role, role.ID)

	c.JSON(http.StatusOK, role)
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// Permission handlers

func GetPermissions(c *gin.Context) {
	var permissions []models.Permission

	if err := database.DB.Preload("Roles").Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

func GetPermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission

	if err := database.DB.Preload("Roles").First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	c.JSON(http.StatusOK, permission)
}

type CreatePermissionRequest struct {
	Action     string                 `json:"action" binding:"required"`
	Subject    string                 `json:"subject" binding:"required"`
	Properties map[string]interface{} `json:"properties"`
	Conditions map[string]interface{} `json:"conditions"`
}

func CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var propertiesJSON, conditionsJSON string
	if req.Properties != nil {
		bytes, _ := json.Marshal(req.Properties)
		propertiesJSON = string(bytes)
	}
	if req.Conditions != nil {
		bytes, _ := json.Marshal(req.Conditions)
		conditionsJSON = string(bytes)
	}

	permission := models.Permission{
		Action:     req.Action,
		Subject:    req.Subject,
		Properties: propertiesJSON,
		Conditions: conditionsJSON,
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

func UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission

	if err := database.DB.First(&permission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Action != "" {
		permission.Action = req.Action
	}
	if req.Subject != "" {
		permission.Subject = req.Subject
	}
	if req.Properties != nil {
		bytes, _ := json.Marshal(req.Properties)
		permission.Properties = string(bytes)
	}
	if req.Conditions != nil {
		bytes, _ := json.Marshal(req.Conditions)
		permission.Conditions = string(bytes)
	}

	if err := database.DB.Save(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

func DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Permission{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission deleted successfully"})
}
