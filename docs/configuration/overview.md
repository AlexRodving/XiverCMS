# Конфигурация

XiverCRM настраивается через переменные окружения и конфигурационные файлы.

## Переменные окружения

Все настройки задаются через файл `.env` в корне проекта.

## Основные разделы конфигурации

- [Переменные окружения](environment.md) - описание всех переменных
- [База данных](database.md) - настройка БД
- [Сервер](server.md) - настройка сервера
- [CORS](cors.md) - настройка CORS

## Быстрая настройка

1. Скопируйте `.env.example` в `.env`
2. Отредактируйте необходимые параметры
3. Перезапустите приложение

## Пример конфигурации

```env
# Server
PORT=8080
GIN_MODE=release

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=xivercrm

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# CORS
CORS_ORIGIN=http://localhost:5173
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

