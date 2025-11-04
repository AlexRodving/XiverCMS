# Создание Handlers

Руководство по созданию HTTP handlers в XiverCMS.

## Структура Handler

```go
package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/xivercms/xivercms/database"
    "github.com/xivercms/xivercms/models"
)

func GetMyModels(c *gin.Context) {
    var models []models.MyModel
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    offset := (page - 1) * pageSize
    
    var total int64
    database.DB.Model(&models.MyModel{}).Count(&total)
    
    if err := database.DB.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": models,
        "meta": gin.H{
            "pagination": gin.H{
                "page":     page,
                "pageSize": pageSize,
                "total":    total,
            },
        },
    })
}
```

## CRUD операции

### GET (список)

```go
func GetMyModels(c *gin.Context) {
    var models []models.MyModel
    // ... пагинация и фильтрация
    database.DB.Find(&models)
    c.JSON(http.StatusOK, gin.H{"data": models})
}
```

### GET (один)

```go
func GetMyModel(c *gin.Context) {
    id := c.Param("id")
    var model models.MyModel
    
    if err := database.DB.First(&model, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
        return
    }
    
    c.JSON(http.StatusOK, model)
}
```

### POST (создание)

```go
func CreateMyModel(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    model := models.MyModel{
        Name: req.Name,
    }
    
    if err := database.DB.Create(&model).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, model)
}
```

### PUT (обновление)

```go
func UpdateMyModel(c *gin.Context) {
    id := c.Param("id")
    var model models.MyModel
    
    if err := database.DB.First(&model, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
        return
    }
    
    var req struct {
        Name string `json:"name"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if req.Name != "" {
        model.Name = req.Name
    }
    
    database.DB.Save(&model)
    c.JSON(http.StatusOK, model)
}
```

### DELETE

```go
func DeleteMyModel(c *gin.Context) {
    id := c.Param("id")
    
    if err := database.DB.Delete(&models.MyModel{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
```

## Регистрация маршрутов

Добавьте в `routes/routes.go`:

```go
myModels := protected.Group("/my-models")
{
    myModels.GET("", handlers.GetMyModels)
    myModels.GET("/:id", handlers.GetMyModel)
    myModels.POST("", handlers.CreateMyModel)
    myModels.PUT("/:id", handlers.UpdateMyModel)
    myModels.DELETE("/:id", handlers.DeleteMyModel)
}
```

## Обработка ошибок

```go
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}
```

## Валидация

Используйте Gin binding:

```go
type CreateRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

var req CreateRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

