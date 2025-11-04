# API: Content Entries

Управление записями контента.

## Получить список записей

**Endpoint:** `GET /api/content-types/:uid/entries`

**Параметры:**
- `page` - номер страницы
- `pageSize` - размер страницы
- `status` - фильтр по статусу (draft, published)

**Пример:**
```bash
curl http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: Bearer YOUR_TOKEN"
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
        "email": "user@example.com"
      },
      "updatedBy": {
        "id": 1,
        "email": "user@example.com"
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

## Получить запись

**Endpoint:** `GET /api/content-types/:uid/entries/:id`

**Пример:**
```bash
curl http://localhost:8080/api/content-types/article/entries/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Создать запись

**Endpoint:** `POST /api/content-types/:uid/entries`

**Тело запроса:**
```json
{
  "data": {
    "title": "My Article",
    "slug": "my-article",
    "content": "Article content",
    "published": true
  },
  "status": "published"
}
```

**Параметры:**
- `data` (required) - данные записи согласно схеме Content Type
- `status` - статус: `draft` или `published` (по умолчанию: `draft`)

**Пример:**
```bash
curl -X POST http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "title": "My First Article",
      "slug": "my-first-article",
      "content": "This is the content",
      "published": true
    },
    "status": "published"
  }'
```

## Обновить запись

**Endpoint:** `PUT /api/content-types/:uid/entries/:id`

**Тело запроса:** (аналогично созданию)

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/content-types/article/entries/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "title": "Updated Article",
      "content": "Updated content"
    },
    "status": "published"
  }'
```

## Удалить запись

**Endpoint:** `DELETE /api/content-types/:uid/entries/:id`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/content-types/article/entries/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## История изменений

**Endpoint:** `GET /api/content-types/:uid/entries/:id/history`

**Пример:**
```bash
curl http://localhost:8080/api/content-types/article/entries/1/history \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Ответ:**
```json
[
  {
    "id": 1,
    "contentEntryId": 1,
    "data": {
      "title": "Updated Article",
      "content": "Updated content"
    },
    "changeType": "updated",
    "changeNote": "Entry updated",
    "changedBy": {
      "id": 1,
      "email": "user@example.com"
    },
    "createdAt": "2024-01-01T12:00:00Z"
  },
  {
    "id": 2,
    "contentEntryId": 1,
    "data": {
      "title": "My Article",
      "content": "Article content"
    },
    "changeType": "created",
    "changeNote": "Entry created",
    "changedBy": {
      "id": 1,
      "email": "user@example.com"
    },
    "createdAt": "2024-01-01T00:00:00Z"
  }
]
```

## Статусы записей

- `draft` - черновик (не опубликован)
- `published` - опубликован (доступен публично)

При изменении статуса на `published`, автоматически устанавливается `publishedAt`.

