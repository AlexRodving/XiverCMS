package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xivercrm/xivercrm/database"
	"github.com/xivercrm/xivercrm/models"
)

func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.AuditLog{})

	// Filter by action
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}

	// Filter by subject
	if subject := c.Query("subject"); subject != "" {
		query = query.Where("subject = ?", subject)
	}

	// Filter by user
	if userId := c.Query("userId"); userId != "" {
		query = query.Where("user_id = ?", userId)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).
		Preload("User").
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

func GetContentHistory(c *gin.Context) {
	entryID := c.Param("id")
	var history []models.ContentHistory

	if err := database.DB.Where("content_entry_id = ?", entryID).
		Preload("ChangedBy").
		Order("created_at DESC").
		Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// Helper function to create audit log
func CreateAuditLog(c *gin.Context, action, subject string, subjectID *uint, description string, metadata map[string]interface{}) {
	var userID *uint
	if userId, exists := c.Get("userId"); exists {
		id := userId.(uint)
		userID = &id
	}

	log := models.AuditLog{
		Action:      action,
		Subject:     subject,
		SubjectID:   subjectID,
		Description: description,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		UserID:      userID,
		Metadata:    models.JSONB(metadata),
	}

	database.DB.Create(&log)
}

// Helper function to create content history entry
func CreateContentHistory(entryID uint, changeType, changeNote string, data models.JSONB, changedByID *uint) {
	history := models.ContentHistory{
		ContentEntryID: entryID,
		Data:           data,
		ChangeType:     changeType,
		ChangeNote:     changeNote,
		ChangedByID:    changedByID,
	}

	database.DB.Create(&history)
}
