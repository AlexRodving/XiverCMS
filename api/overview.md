# API Обзор

XiverCMS предоставляет RESTful API для управления контентом, пользователями и настройками системы.

## Базовый URL

```
http://localhost:8080/api
```

## Аутентификация

Большинство endpoints требуют аутентификации. Используйте один из методов:

### JWT Token

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

### API Token

```http
Authorization: xvc_YOUR_API_TOKEN
```

## Формат ответов

### Успешный ответ

```json
{
  "data": { ... },
  "meta": {
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 100
    }
  }
}
```

### Ошибка

```json
{
  "error": "Error message"
}
```

## Коды статусов

- `200 OK` - Успешный запрос
- `201 Created` - Ресурс создан
- `400 Bad Request` - Неверный запрос
- `401 Unauthorized` - Требуется аутентификация
- `403 Forbidden` - Недостаточно прав
- `404 Not Found` - Ресурс не найден
- `409 Conflict` - Конфликт (например, дубликат)
- `500 Internal Server Error` - Ошибка сервера

## Пагинация

Многие endpoints поддерживают пагинацию:

```http
GET /api/users?page=1&pageSize=20
```

Параметры:
- `page` - номер страницы (по умолчанию: 1)
- `pageSize` - размер страницы (по умолчанию: 10)

## Фильтрация

Многие endpoints поддерживают фильтрацию:

```http
GET /api/content-types?isVisible=true
GET /api/media-files?mime=image%
GET /api/audit-logs?action=create&subject=content-entry
```

## Поиск

Некоторые endpoints поддерживают поиск:

```http
GET /api/users?search=john
GET /api/media-files?search=photo
```

## Основные разделы API

- [Аутентификация](authentication.md)
- [Пользователи](users.md)
- [Content Types](content-types.md)
- [Content Entries](content-entries.md)
- [Media Library](media-library.md)
- [Роли и права](roles-permissions.md)
- [API Tokens](api-tokens.md)
- [Audit Logs](audit-logs.md)

