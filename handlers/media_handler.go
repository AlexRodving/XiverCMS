package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xivercms/xivercms/database"
	"github.com/xivercms/xivercms/models"
)

const uploadDir = "./uploads"
const maxUploadSize = 10 << 20 // 10 MB

func init() {
	// Create uploads directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}
}

func UploadMedia(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// Check file size
	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), sanitizeFilename(file.Filename[:len(file.Filename)-len(ext)]), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Copy file
	_, err = io.Copy(dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Get file info
	fileInfo, _ := os.Stat(filePath)

	// Get image dimensions if it's an image
	width, height := 0, 0
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType([]byte{})
	}

	// Create media file record
	userId, _ := c.Get("userId")
	userID := userId.(uint)

	mediaFile := models.MediaFile{
		Name:        file.Filename,
		Mime:        mimeType,
		Size:        fileInfo.Size(),
		Format:      strings.TrimPrefix(ext, "."),
		URL:         fmt.Sprintf("/api/uploads/%s", filename),
		Path:        filePath,
		Provider:    "local",
		Width:       width,
		Height:      height,
		CreatedByID: &userID,
	}

	// Get alternative text and caption from form
	if alt := c.PostForm("alternative"); alt != "" {
		mediaFile.Alternative = alt
	}
	if caption := c.PostForm("caption"); caption != "" {
		mediaFile.Caption = caption
	}

	if err := database.DB.Create(&mediaFile).Error; err != nil {
		os.Remove(filePath) // Clean up file if DB save fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media file record"})
		return
	}

	c.JSON(http.StatusCreated, mediaFile)
}

func GetMediaFiles(c *gin.Context) {
	var mediaFiles []models.MediaFile
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.MediaFile{})

	// Search
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR alternative LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Filter by MIME type
	if mime := c.Query("mime"); mime != "" {
		query = query.Where("mime LIKE ?", mime+"%")
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).
		Preload("CreatedBy").
		Order("created_at DESC").
		Find(&mediaFiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mediaFiles,
		"meta": gin.H{
			"pagination": gin.H{
				"page":     page,
				"pageSize": pageSize,
				"total":    total,
			},
		},
	})
}

func GetMediaFile(c *gin.Context) {
	id := c.Param("id")
	var mediaFile models.MediaFile

	if err := database.DB.Preload("CreatedBy").First(&mediaFile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media file not found"})
		return
	}

	c.JSON(http.StatusOK, mediaFile)
}

func UpdateMediaFile(c *gin.Context) {
	id := c.Param("id")
	var mediaFile models.MediaFile

	if err := database.DB.First(&mediaFile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media file not found"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Alternative string `json:"alternative"`
		Caption     string `json:"caption"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		mediaFile.Name = req.Name
	}
	if req.Alternative != "" {
		mediaFile.Alternative = req.Alternative
	}
	if req.Caption != "" {
		mediaFile.Caption = req.Caption
	}

	if err := database.DB.Save(&mediaFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mediaFile)
}

func DeleteMediaFile(c *gin.Context) {
	id := c.Param("id")
	var mediaFile models.MediaFile

	if err := database.DB.First(&mediaFile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media file not found"})
		return
	}

	// Delete file from filesystem
	if err := os.Remove(mediaFile.Path); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete record from database
	if err := database.DB.Delete(&mediaFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media file deleted successfully"})
}

func ServeMediaFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(uploadDir, filename)

	// Security check - prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(uploadDir)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.File(filePath)
}

func sanitizeFilename(filename string) string {
	// Remove special characters and spaces
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "..", "")
	return filename
}
