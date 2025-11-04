# Публичные API

XiverCMS предоставляет публичные API endpoints для доступа к **данным контента** без аутентификации (если Content Type настроен как публичный).

⚠️ **Важно:** Endpoints для получения схем Content Types (`/api/content-types`) **не являются публичными** и требуют аутентификации. Они используются только в админ-панели. Схемы содержат внутреннюю структуру данных и не должны быть доступны публично.

## Публичные endpoints для данных записей

**URL:** `/api/:uid`

**Описание:** Возвращает **записи** (конкретные данные) типа контента - только опубликованные записи.

**Пример:**
```bash
GET /api/articles
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "createdAt": "2025-11-04T10:00:00Z",
      "updatedAt": "2025-11-04T10:00:00Z",
      "publishedAt": "2025-11-04T10:00:00Z",
      "data": {
        "title": "Моя первая статья",
        "content": "Содержание статьи..."
      },
      "status": "published"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 1
    }
  }
}
```

**Использование:**
- Для получения опубликованного контента
- Для отображения на фронтенде
- Для интеграции с другими приложениями

## Типы доступа (AccessType)

Каждый Content Type может иметь один из следующих типов доступа:

- **`public`** - полностью публичный доступ, не требует аутентификации
- **`authenticated`** - требует авторизованного пользователя (любая роль)
- **`moderator`** - требует роль Moderator или Admin
- **`admin`** - требует роль Admin или Super Admin

## Базовый URL

```
http://localhost:8080/api
```

## Получить публичные записи Content Type

**Endpoint:** `GET /api/:uid`

**Упрощенный URL** для получения записей Content Type.

**Пример:** `GET /api/articles` вместо `GET /api/content-types/articles/entries`

Возвращает только опубликованные (`status = "published"`) записи. Доступ контролируется через `accessType` Content Type.

**Параметры:**
- `page` - номер страницы (по умолчанию: 1)
- `pageSize` - размер страницы (по умолчанию: 10)
- `search` - поиск по содержимому записи
- `populate` - загрузить связанные записи (true/false)

**Доступ:**
- Если `accessType = "public"` - не требует аутентификации
- Если `accessType = "authenticated"` - требует JWT токен
- Если `accessType = "moderator"` - требует роль Moderator/Admin
- Если `accessType = "admin"` - требует роль Admin/Super Admin

**Пример (публичный доступ):**
```bash
curl http://localhost:8080/api/articles
```

**Пример (с аутентификацией):**
```bash
curl http://localhost:8080/api/articles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Пример (с поиском):**
```bash
curl "http://localhost:8080/api/articles?search=news"
```

**Пример (с populate):**
```bash
curl "http://localhost:8080/api/articles/1?populate=true"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "contentTypeId": 1,
      "data": {
        "title": "My Article",
        "content": "Article content",
        "published": true
      },
      "status": "published",
      "publishedAt": "2024-01-01T00:00:00Z",
      "createdBy": {
        "id": 1,
        "username": "author"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 1
    }
  }
}
```

**Ошибки:**
- `404 Not Found` - если Content Type не найден
- `403 Forbidden` - если недостаточно прав доступа

## Получить публичную запись

**Endpoint:** `GET /api/:uid/:id`

**Упрощенный URL** для получения одной записи.

**Пример:** `GET /api/articles/1` вместо `GET /api/content-types/articles/entries/1`

Возвращает запись только если она опубликована. Доступ контролируется через `accessType` Content Type.

**Параметры:**
- `populate` - загрузить связанные записи (true/false)

**Пример (публичный доступ):**
```bash
curl http://localhost:8080/api/articles/1
```

**Пример (с populate):**
```bash
curl "http://localhost:8080/api/articles/1?populate=true"
```

**Ответ:**
```json
{
  "id": 1,
  "contentTypeId": 1,
  "data": {
    "title": "My Article",
    "content": "Article content",
    "published": true
  },
  "status": "published",
  "publishedAt": "2024-01-01T00:00:00Z",
  "createdBy": {
    "id": 1,
    "email": "author@example.com"
  },
  "updatedBy": {
    "id": 1,
    "email": "author@example.com"
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Ошибки:**
- `