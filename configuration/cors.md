# CORS Конфигурация

Настройка Cross-Origin Resource Sharing (CORS).

## Переменные окружения

### CORS_ORIGIN

Основной разрешенный origin:

```env
CORS_ORIGIN=http://localhost:5173
```

### ALLOWED_ORIGINS

Список разрешенных origins (через запятую):

```env
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000,https://example.com
```

## Настройка для разработки

```env
CORS_ORIGIN=http://localhost:5173
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

## Настройка для production

```env
CORS_ORIGIN=https://yourdomain.com
ALLOWED_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com
```

## Разрешенные методы

- GET
- POST
- PUT
- PATCH
- DELETE
- OPTIONS

## Разрешенные заголовки

- Origin
- Content-Type
- Accept
- Authorization

## Безопасность

⚠️ **ВАЖНО:** В production настройте только ваши домены!

Не используйте:
```env
ALLOWED_ORIGINS=*
```

Используйте конкретные домены:
```env
ALLOWED_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com
```

