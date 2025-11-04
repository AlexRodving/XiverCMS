# Добавление моделей

Руководство по созданию новых моделей данных в XiverCMS.

## Создание модели

Создайте новый файл в `models/`:

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type MyModel struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Ваши поля
    Name        string `json:"name" gorm:"not null"`
    Description string `json:"description"`
}
```

## Миграция

Добавьте модель в `database/database.go`:

```go
func Migrate() {
    err := DB.AutoMigrate(
        &models.User{},
        // ... другие модели
        &models.MyModel{},  // Добавьте вашу модель
    )
    // ...
}
```

## GORM теги

### Основные теги

- `primaryKey` - первичный ключ
- `not null` - обязательное поле
- `uniqueIndex` - уникальный индекс
- `index` - обычный индекс
- `default:value` - значение по умолчанию
- `type:jsonb` - JSON поле (PostgreSQL)
- `type:json` - JSON поле (SQLite)

### Связи

```go
// One-to-Many
CreatedByID *uint `json:"createdById"`
CreatedBy   *User `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`

// Many-to-Many
Roles []Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
```

## Примеры

### С JSON полем

```go
type MyModel struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    
    Data JSONB `json:"data" gorm:"type:jsonb"`
}
```

### С связью к User

```go
type MyModel struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    
    CreatedByID *uint `json:"createdById"`
    CreatedBy   *User `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
}
```

## Использование

```go
// Создание
model := models.MyModel{
    Name: "My Model",
}
database.DB.Create(&model)

// Получение
var model models.MyModel
database.DB.First(&model, id)

// Обновление
database.DB.Save(&model)

// Удаление (soft delete)
database.DB.Delete(&model)
```

