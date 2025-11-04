package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

func GetContentTypes(c *gin.Context) {
	var contentTypes []models.ContentType
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.ContentType{})

	// Filter by visibility
	if isVisible := c.Query("isVisible"); isVisible != "" {
		visible, _ := strconv.ParseBool(isVisible)
		query = query.Where("is_visible = ?", visible)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&contentTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": contentTypes,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

func GetContentType(c *gin.Context) {
	uid := c.Param("uid")
	var contentType models.ContentType

	if err := database.DB.Where("uid = ?", uid).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	c.JSON(http.StatusOK, contentType)
}

type CreateContentTypeRequest struct {
	UID         string                 `json:"uid" binding:"required"`
	Kind        string                 `json:"kind"`
	DisplayName string                 `json:"displayName" binding:"required"`
	Description string                 `json:"description"`
	IsVisible   bool                   `json:"isVisible"`
	AccessType  string                 `json:"accessType"`
	Schema      map[string]interface{} `json:"schema" binding:"required"`
}

type UpdateContentTypeRequest struct {
	UID         string                 `json:"uid"` // Optional for updates
	Kind        string                 `json:"kind"`
	DisplayName string                 `json:"displayName"`
	Description string                 `json:"description"`
	IsVisible   bool                   `json:"isVisible"`
	AccessType  string                 `json:"accessType"`
	Schema      map[string]interface{} `json:"schema"`
}

func CreateContentType(c *gin.Context) {
	var req CreateContentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if UID already exists
	var existing models.ContentType
	if err := database.DB.Where("uid = ?", req.UID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Content type with this UID already exists"})
		return
	}

	contentType := models.ContentType{
		UID:         req.UID,
		Kind:        req.Kind,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsVisible:   req.IsVisible,
		AccessType:  req.AccessType,
		Schema:      models.JSONB(req.Schema),
	}

	if contentType.Kind == "" {
		contentType.Kind = "collectionType"
	}
	if contentType.AccessType == "" {
		contentType.AccessType = "public"
	}

	if err := database.DB.Create(&contentType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contentType)
}

func UpdateContentType(c *gin.Context) {
	uid := c.Param("uid")
	var contentType models.ContentType

	if err := database.DB.Where("uid = ?", uid).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	var req UpdateContentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DisplayName != "" {
		contentType.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		contentType.Description = req.Description
	}
	if req.Schema != nil {
		contentType.Schema = models.JSONB(req.Schema)
	}
	if req.AccessType != "" {
		contentType.AccessType = req.AccessType
	}
	contentType.IsVisible = req.IsVisible

	if err := database.DB.Save(&contentType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contentType)
}

func DeleteContentType(c *gin.Context) {
	uid := c.Param("uid")
	if err := database.DB.Where("uid = ?", uid).Delete(&models.ContentType{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content type deleted successfully"})
}

// Content Entry Handlers

func GetContentEntries(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	var contentType models.ContentType

	if err := database.DB.Where("uid = ?", contentTypeUID).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	var entries []models.ContentEntry
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	query := database.DB.Where("content_type_id = ?", contentType.ID)

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.ContentEntry{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).
		Preload("CreatedBy").Preload("UpdatedBy").
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

func GetContentEntry(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")

	var contentType models.ContentType
	if err := database.DB.Where("uid = ?", contentTypeUID).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	var entry models.ContentEntry
	if err := database.DB.Where("id = ? AND content_type_id = ?", entryID, contentType.ID).
		Preload("CreatedBy").Preload("UpdatedBy").
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
			if err := database.DB.Where("content_type_id = ? AND id = ?", relatedContentType.ID, relation.TargetEntryID).First(&relatedEntry).Error; err == nil {
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

type CreateContentEntryRequest struct {
	Data   map[string]interface{} `json:"data" binding:"required"`
	Status string                 `json:"status"`
}

func CreateContentEntry(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	var contentType models.ContentType

	if err := database.DB.Where("uid = ?", contentTypeUID).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	var req CreateContentEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("userId")
	userID := userId.(uint)

	// Separate relation fields from regular data
	entryData := make(map[string]interface{})
	relationData := make(map[string]interface{})

	for k, v := range req.Data {
		// Check if this is a relation field
		schemaMap := map[string]interface{}(contentType.Schema)
		if fieldDef, exists := schemaMap[k]; exists {
			if fieldMap, ok := fieldDef.(map[string]interface{}); ok {
				if fieldType, _ := fieldMap["type"].(string); fieldType == "relation" {
					relationData[k] = v
					continue
				}
			}
		}
		entryData[k] = v
	}

	entry := models.ContentEntry{
		ContentTypeID: contentType.ID,
		Data:          models.JSONB(entryData),
		Status:        req.Status,
		CreatedByID:   &userID,
		UpdatedByID:   &userID,
	}

	if entry.Status == "" {
		entry.Status = "draft"
	}

	if entry.Status == "published" {
		now := time.Now()
		entry.PublishedAt = &now
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Preload("CreatedBy").Preload("UpdatedBy").First(&entry, entry.ID)

	// Process relations if any
	if len(relationData) > 0 {
		processRelations(c, contentTypeUID, entry.ID, relationData, contentType.Schema)
	}

	// Create audit log
	CreateAuditLog(c, "create", "content-entry", &entry.ID, "Created content entry", map[string]interface{}{
		"contentType": contentTypeUID,
		"status":      entry.Status,
	})

	// Create content history
	CreateContentHistory(entry.ID, "created", "Entry created", entry.Data, entry.CreatedByID)

	c.JSON(http.StatusCreated, entry)
}

func UpdateContentEntry(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")

	var contentType models.ContentType
	if err := database.DB.Where("uid = ?", contentTypeUID).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	var entry models.ContentEntry
	if err := database.DB.Where("id = ? AND content_type_id = ?", entryID, contentType.ID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	var req CreateContentEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Separate relation fields from regular data
	relationData := make(map[string]interface{})

	if req.Data != nil {
		// Get current data
		currentData := make(map[string]interface{})
		if entry.Data != nil {
			for k, v := range entry.Data {
				currentData[k] = v
			}
		}

		// Merge with new data
		for k, v := range req.Data {
			// Check if this is a relation field
			schemaMap := map[string]interface{}(contentType.Schema)
			if fieldDef, exists := schemaMap[k]; exists {
				if fieldMap, ok := fieldDef.(map[string]interface{}); ok {
					if fieldType, _ := fieldMap["type"].(string); fieldType == "relation" {
						relationData[k] = v
						continue
					}
				}
			}
			currentData[k] = v
		}
		entry.Data = models.JSONB(currentData)
	}

	if req.Status != "" {
		entry.Status = req.Status

		if req.Status == "published" && entry.PublishedAt == nil {
			now := time.Now()
			entry.PublishedAt = &now
		}
	}

	userId, _ := c.Get("userId")
	userID := userId.(uint)
	entry.UpdatedByID = &userID

	if err := database.DB.Save(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Preload("CreatedBy").Preload("UpdatedBy").First(&entry, entry.ID)

	// Process relations if any
	if len(relationData) > 0 {
		processRelations(c, contentTypeUID, entry.ID, relationData, contentType.Schema)
	}

	// Create audit log
	changeType := "updated"
	if req.Status == "published" && entry.PublishedAt != nil {
		changeType = "published"
	}
	CreateAuditLog(c, "update", "content-entry", &entry.ID, "Updated content entry", map[string]interface{}{
		"contentType": contentTypeUID,
		"status":      entry.Status,
	})

	// Create content history
	CreateContentHistory(entry.ID, changeType, "Entry updated", entry.Data, entry.UpdatedByID)

	c.JSON(http.StatusOK, entry)
}

func DeleteContentEntry(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")

	var contentType models.ContentType
	if err := database.DB.Where("uid = ?", contentTypeUID).First(&contentType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content type not found"})
		return
	}

	if err := database.DB.Where("id = ? AND content_type_id = ?", entryID, contentType.ID).
		Delete(&models.ContentEntry{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}

// processRelations processes relation fields and creates ContentRelation records
func processRelations(c *gin.Context, contentTypeUID string, entryID uint, data map[string]interface{}, schema models.JSONB) {
	if schema == nil {
		return
	}

	// Get schema as map
	schemaMap := map[string]interface{}(schema)

	// Process each field in data
	for fieldName, fieldValue := range data {
		// Check if field is defined in schema as relation
		fieldDef, exists := schemaMap[fieldName]
		if !exists {
			continue
		}

		fieldMap, ok := fieldDef.(map[string]interface{})
		if !ok {
			continue
		}

		fieldType, _ := fieldMap["type"].(string)
		if fieldType != "relation" {
			continue
		}

		// Get relation configuration
		targetContentType, _ := fieldMap["targetContentType"].(string)
		relationType, _ := fieldMap["relationType"].(string)
		if targetContentType == "" {
			continue
		}

		// Delete existing relations for this field
		database.DB.Where("source_content_type_uid = ? AND source_entry_id = ? AND source_field_name = ?",
			contentTypeUID, entryID, fieldName).Delete(&models.ContentRelation{})

		// Create new relations
		if relationType == "oneToMany" || relationType == "manyToMany" {
			// Handle array of IDs
			if ids, ok := fieldValue.([]interface{}); ok {
				for idx, id := range ids {
					var targetID uint
					if idFloat, ok := id.(float64); ok {
						targetID = uint(idFloat)
					} else if idMap, ok := id.(map[string]interface{}); ok {
						if idVal, exists := idMap["id"]; exists {
							if idFloat, ok := idVal.(float64); ok {
								targetID = uint(idFloat)
							}
						}
					}

					if targetID > 0 {
						relation := models.ContentRelation{
							SourceContentTypeUID: contentTypeUID,
							SourceEntryID:        entryID,
							SourceFieldName:      fieldName,
							TargetContentTypeUID: targetContentType,
							TargetEntryID:        targetID,
							RelationType:         relationType,
							Order:                idx,
						}
						database.DB.Create(&relation)
					}
				}
			}
		} else {
			// Handle single ID
			var targetID uint
			if idFloat, ok := fieldValue.(float64); ok {
				targetID = uint(idFloat)
			} else if idMap, ok := fieldValue.(map[string]interface{}); ok {
				if idVal, exists := idMap["id"]; exists {
					if idFloat, ok := idVal.(float64); ok {
						targetID = uint(idFloat)
					}
				}
			}

			if targetID > 0 {
				relation := models.ContentRelation{
					SourceContentTypeUID: contentTypeUID,
					SourceEntryID:        entryID,
					SourceFieldName:      fieldName,
					TargetContentTypeUID: targetContentType,
					TargetEntryID:        targetID,
					RelationType:         relationType,
				}
				database.DB.Create(&relation)
			}
		}

		// Relation is stored separately in ContentRelation table
		// Field is already excluded from entry.Data
	}
}
