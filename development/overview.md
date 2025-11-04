# Разработка - Обзор

Руководство по разработке и расширению XiverCRM.

## Основные разделы

- [Добавление моделей](models.md) - создание новых моделей данных
- [Создание handlers](handlers.md) - обработка HTTP запросов
- [Middleware](middleware.md) - создание middleware
- [Тестирование](testing.md) - написание тестов

## Структура проекта

```
XiverCRM/
├── config/          # Конфигурация
├── models/          # Модели данных
├── database/        # Работа с БД
├── handlers/        # HTTP handlers
├── middleware/      # Middleware
├── routes/          # Маршруты
└── auth/            # Аутентификация
```

## Разработка

### Запуск в режиме разработки

```bash
# Backend
GIN_MODE=debug go run main.go

# Frontend
cd frontend && npm run dev
```

### Горячая перезагрузка

Для Go используйте:
```bash
go install github.com/cosmtrek/air@latest
air
```

Для Vue используйте Vite (автоматически).

## Стиль кода

- Следуйте Go conventions
- Используйте gofmt для форматирования
- Комментируйте публичные функции
- Пишите тесты для новой функциональности

