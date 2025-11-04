package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

// GetRelations returns relations for a content entry
func GetRelations(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")
	fieldName := c.Query("field") // Optional: filter by field name

	query := database.DB.Where("source_content_type_uid = ? AND source_entry_id = ?", contentTypeUID, entryID)

	if fieldName != "" {
		query = query.Where("source_field_name = ?", fieldName)
	}

	var relations []models.ContentRelation
	if err := query.Order("order ASC").Find(&relations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, relations)
}

// CreateRelation creates a relation between entries
func CreateRelation(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")

	var req struct {
		FieldName            string `json:"fieldName" binding:"required"`
		TargetContentTypeUID string `json:"targetContentTypeUid" binding:"required"`
		TargetEntryID        uint   `json:"targetEntryId" binding:"required"`
		RelationType         string `json:"relationType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RelationType == "" {
		req.RelationType = "manyToOne"
	}

	// Convert entryID from string to uint
	entryIDUint, err := strconv.ParseUint(entryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID"})
		return
	}

	relation := models.ContentRelation{
		SourceContentTypeUID: contentTypeUID,
		SourceEntryID:        uint(entryIDUint),
		SourceFieldName:      req.FieldName,
		TargetContentTypeUID: req.TargetContentTypeUID,
		TargetEntryID:        req.TargetEntryID,
		RelationType:         req.RelationType,
	}

	if err := database.DB.Create(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, relation)
}

// DeleteRelation deletes a relation
func DeleteRelation(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")
	relationID := c.Param("relationId")

	if err := database.DB.Where("id = ? AND source_content_type_uid = ? AND source_entry_id = ?",
		relationID, contentTypeUID, entryID).Delete(&models.ContentRelation{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relation deleted successfully"})
}

// GetRelatedEntries returns related entries for a content entry
func GetRelatedEntries(c *gin.Context) {
	contentTypeUID := c.Param("uid")
	entryID := c.Param("id")
	fieldName := c.Param("field")

	// Get relations
	var relations []models.ContentRelation
	if err := database.DB.Where("source_content_type_uid = ? AND source_entry_id = ? AND source_field_name = ?",
		contentTypeUID, entryID, fieldName).Find(&relations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get related entries
	result := make([]map[string]interface{}, 0)
	for _, relation := range relations {
		var entry models.ContentEntry
		if err := database.DB.Where("content_type_id = (SELECT id FROM content_types WHERE uid = ?) AND id = ?",
			relation.TargetContentTypeUID, relation.TargetEntryID).First(&entry).Error; err == nil {
			result = append(result, map[string]interface{}{
				"entry":    entry,
				"relation": relation,
			})
		}
	}

	c.JSON(http.StatusOK, result)
}
