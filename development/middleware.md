# Middleware

Руководство по созданию middleware в XiverCMS.

## Создание Middleware

```go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Логика перед обработкой запроса
        
        // Вызов следующего handler
        c.Next()
        
        // Логика после обработки запроса
    }
}
```

## Примеры

### Логирование

```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        log.Printf("%s %s - %v", c.Request.Method, c.Request.URL.Path, duration)
    }
}
```

### Проверка прав

```go
func RequireAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        userId, exists := c.Get("userId")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        
        var user models.User
        database.DB.First(&user, userId)
        
        if !user.IsSuperAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## Использование

В `routes/routes.go`:

```go
protected := r.Group("/api")
protected.Use(middleware.AuthMiddleware())
protected.Use(middleware.MyMiddleware())  // Ваш middleware
{
    // routes
}
```

Или для конкретной группы:

```go
admin := protected.Group("/admin")
admin.Use(middleware.RequireAdmin())
{
    // admin routes
}
```

## Прерывание выполнения

Используйте `c.Abort()` для остановки выполнения:

```go
if condition {
    c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
    c.Abort()
    return
}
```

## Передача данных

Используйте `c.Set()` и `c.Get()`:

```go
// В middleware
c.Set("myKey", myValue)

// В handler
value, exists := c.Get("myKey")
```

