# Тестирование

Руководство по написанию тестов для XiverCMS.

## Настройка тестов

Создайте `*_test.go` файлы рядом с кодом.

## Unit тесты

```go
package handlers_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetMyModel(t *testing.T) {
    // Подготовка
    // ...
    
    // Выполнение
    // ...
    
    // Проверка
    assert.Equal(t, expected, actual)
}
```

## Интеграционные тесты

```go
func TestCreateModel(t *testing.T) {
    // Настройка тестовой БД
    // ...
    
    // Создание модели
    // ...
    
    // Проверка результата
    // ...
}
```

## Использование testify

```bash
go get github.com/stretchr/testify
```

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
    assert.Equal(t, expected, actual)
    require.NotNil(t, obj)
}
```

## Запуск тестов

```bash
# Все тесты
go test ./...

# Конкретный пакет
go test ./handlers

# С покрытием
go test -cover ./...

# Вербозный вывод
go test -v ./...
```

## Mock база данных

Для тестов используйте тестовую БД:

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // Миграции
    db.AutoMigrate(&models.MyModel{})
    
    return db
}
```

