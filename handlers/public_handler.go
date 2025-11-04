package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/middleware"
	"github.com/xivercms/xivercms/models"
)

// PublicGetContentTypes - get content types (public or admin based on auth)
// If authenticated as admin, returns all content types including non-visible ones
// Otherwise, returns only visible and accessible content types
func PublicGetContentTypes(c *gin.Context) {
	var allContentTypes []models.ContentType
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// Check if user is authenticated and is admin
	userID, exists := c.Get("userID")
	isAdmin := false
	if exists && userID != nil {
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, userID).Error; err == nil {
			// Check if user is super admin or has admin role
			if user.IsSuperAdmin {
				isAdmin = true
			} else {
				for _, role := range user.Roles {
					if role.Name == "Admin" || role.Type == "custom" {
						isAdmin = true
						break
					}
				}
			}
		}
	}

	// Build query based on admin status
	baseQuery := database.DB.Model(&models.ContentType{})
	if !isAdmin {
		// Public users only see visible content types
		baseQuery = baseQuery.Where("is_visible = ?", true)
	}

	var total int64
	baseQuery.Count(&total)

	if err := baseQuery.Offset(offset).Limit(pageSize).Find(&allContentTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by access - only include types user can access (if not admin)
	accessibleTypes := []models.ContentType{}
	for _, ct := range allContentTypes {
		if isAdmin || middleware.CheckContentTypeAccess(ct.UID, c) {
			accessibleTypes = append(accessibleTypes, ct)
		}
	}

	// Recalculate total for accessible types only (if not admin)
	var accessibleTotal int64
	if isAdmin {
		accessibleTotal = total
	} else {
		allVisible := []models.ContentType{}
		database.DB.Model(&models.ContentType{}).Where("is_visible = ?", true).Find(&allVisible)
		for _, ct := range allVisible {
			if middleware.CheckContentTypeAccess(ct.UID, c) {
				accessibleTotal++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": accessibleTypes,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    accessibleTotal,
			},
		},
	})
}

// PublicGetContentType - get content type by UID (public or admin based on auth)
// If authenticated as admin, returns any content type including non-visible ones
// Otherwise, returns only visible and accessible content types
func PublicGetContentType(c *gin.Context) {
	uid := c.Param("uid")
	var contentType models.ContentType

	// Check if user is authenticated and is admin
	userID, exists := c.Get("userID")
	isAdmin := false
	if exists && userID != nil {
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, userID).Error; err == nil {
			if user.IsSuperAdmin {
				isAdmin = true
			} else {
				for _, role := range user.Roles {
					if role.Name == "Admin" || role.Type == "custom" {
						isAdmin = true
						break
					}
				}
			}
		}
	}

	// Build query based on admin status
	query := database.DB.Where("uid = ?", uid)
	if !isAdmin {
		query = query.Where("is_visible = ?", true)
	}

	if err := query.First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	// Check access based on accessType (if not admin)
	if !isAdmin && !middleware.CheckContentTypeAccess(uid, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, contentType)
}

// PublicGetContentEntries - get public content entries
// Handles /api/{uid} - simplified URL like Strapi
func PublicGetContentEntries(c *gin.Context) {
	contentTypeUID := c.Param("uid")

	// Skip if this is a reserved route (shouldn't happen if routes are ordered correctly)
	reservedRoutes := []string{"auth", "roles", "users", "permissions", "api-tokens", "upload", "media-files", "content-types", "admin", "audit-logs"}
	for _, reserved := range reservedRoutes {
		if contentTypeUID == reserved {
			c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
			return
		}
	}

	var contentType models.ContentType
	if err := database.DB.Where("uid = ? AND is_visible = ?", contentTypeUID, true).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	// Check access
	if !middleware.CheckContentTypeAccess(contentTypeUID, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.ContentEntry{}).Where("content_type_id = ? AND status = ?", contentType.ID, "published")

	// Search (simplified for SQLite compatibility)
	if search := c.Query("search"); search != "" {
		// Simple LIKE search on JSON text - works with both SQLite and PostgreSQL
		query = query.Where("data LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var entries []models.ContentEntry
	if err := query.Offset(offset).Limit(pageSize).
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": entries,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

// PublicGetContentEntry - get public content entry
// Handles /api/{uid}/{id} - simplified URL like Strapi
func PublicGetContentEntry(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")

	// Skip if this is a reserved route
	reservedRoutes := []string{"auth", "roles", "users", "permissions", "api-tokens", "upload", "media-files", "content-types", "admin", "audit-logs"}
	for _, reserved := range reservedRoutes {
		if contentTypeUID == reserved {
			c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
			return
		}
	}

	var contentType models.ContentType
	if err := database.DB.Where("uid = ? AND is_visible = ?", contentTypeUID, true).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	// Check access
	if !middleware.CheckContentTypeAccess(contentTypeUID, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var entry models.ContentEntry
	if err := database.DB.Where("id = ? AND content_type_id = ? AND status = ?", entryID, contentType.ID, "published").
		Preload("CreatedBy").
		First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	// Load relations if requested
	if c.Query("populate") == "true" {
		var relations []models.ContentRelation
		database.DB.Where("source_content_type_uid = ? AND source_entry_id = ?", contentTypeUID, entryID).Find(&relations)

		// Populate related entries
		populatedData := make(map[string]interface{})
		for k, v := range entry.Data {
			populatedData[k] = v
		}

		// Add relations to data
		for _, relation := range relations {
			var relatedEntry models.ContentEntry
			var relatedContentType models.ContentType
			database.DB.Where("uid = ?", relation.TargetContentTypeUID).First(&relatedContentType)
			if err := database.DB.Where("content_type_id = ? AND id = ? AND status = ?", relatedContentType.ID, relation.TargetEntryID, "published").First(&relatedEntry).Error; err == nil {
				fieldName := relation.SourceFieldName
				if relation.RelationType == "oneToMany" || relation.RelationType == "manyToMany" {
					if populatedData[fieldName] == nil {
						populatedData[fieldName] = []interface{}{}
					}
					arr := populatedData[fieldName].([]interface{})
					populatedData[fieldName] = append(arr, relatedEntry)
				} else {
					populatedData[fieldName] = relatedEntry
				}
			}
		}
		entry.Data = populatedData
	}

	c.JSON(http.StatusOK, entry)
}
