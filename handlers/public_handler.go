package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/middleware"
	"github.com/xivercms/xivercms/models"
)

// PublicGetContentTypes - get public content types
func PublicGetContentTypes(c *gin.Context) {
	var allContentTypes []models.ContentType
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// Get all visible content types
	baseQuery := database.DB.Model(&models.ContentType{}).Where("is_visible = ?", true)

	var total int64
	baseQuery.Count(&total)

	if err := baseQuery.Offset(offset).Limit(pageSize).Find(&allContentTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by access - only include types user can access
	accessibleTypes := []models.ContentType{}
	for _, ct := range allContentTypes {
		if middleware.CheckContentTypeAccess(ct.UID, c) {
			accessibleTypes = append(accessibleTypes, ct)
		}
	}

	// Recalculate total for accessible types only
	var accessibleTotal int64
	allVisible := []models.ContentType{}
	database.DB.Model(&models.ContentType{}).Where("is_visible = ?", true).Find(&allVisible)
	for _, ct := range allVisible {
		if middleware.CheckContentTypeAccess(ct.UID, c) {
			accessibleTotal++
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

// PublicGetContentType - get public content type by UID
func PublicGetContentType(c *gin.Context) {
	uid := c.Param("uid")
	var contentType models.ContentType

	if err := database.DB.Where("uid = ? AND is_visible = ?", uid, true).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	// Check access based on accessType
	if !middleware.CheckContentTypeAccess(uid, c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, contentType)
}

// PublicGetContentEntries - get public content entries
func PublicGetContentEntries(c *gin.Context) {
	contentTypeUID := c.Param("uid")
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

	query := database.DB.Where("content_type_id = ? AND status = ?", contentType.ID, "published")

	// Search
	if search := c.Query("search"); search != "" {
		query = query.Where("data::text ILIKE ?", "%"+search+"%")
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
func PublicGetContentEntry(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")

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
