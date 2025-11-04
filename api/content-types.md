# API: Content Types

Content Types определяют структуру контента в системе.

## Получить список Content Types

**Endpoint:** `GET /api/content-types`

**Параметры запроса:**
- `page` - номер страницы (по умолчанию: 1)
- `pageSize` - размер страницы (по умолчанию: 10)
- `isVisible` - фильтр по видимости (true/false)

**Пример:**
```bash
curl http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "uid": "article",
      "kind": "collectionType",
      "displayName": "Article",
      "description": "Blog articles",
      "isVisible": true,
      "schema": {
        "title": {"type": "string"},
        "content": {"type": "text"}
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

## Получить Content Type

**Endpoint:** `GET /api/content-types/:uid`

**Пример:**
```bash
curl http://localhost:8080/api/content-types/article \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Создать Content Type

**Endpoint:** `POST /api/content-types`

**Тело запроса:**
```json
{
  "uid": "article",
  "displayName": "Article",
  "description": "Blog articles",
  "kind": "collectionType",
  "isVisible": true,
  "schema": {
    "title": {
      "type": "string",
      "required": true
    },
    "slug": {
      "type": "string",
      "unique": true
    },
    "content": {
      "type": "text"
    },
    "published": {
      "type": "boolean",
      "default": false
    }
  }
}
```

**Параметры:**
- `uid` (required) - уникальный идентификатор
- `displayName` (required) - отображаемое имя
- `description` - описание
- `kind` - тип: `collectionType` или `singleType` (по умолчанию: `collectionType`)
- `isVisible` - видимость (по умолчанию: true)
- `schema` (required) - JSON схема с определением полей

**Пример:**
```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "article",
    "displayName": "Article",
    "kind": "collectionType",
    "schema": {
      "title": {"type": "string"},
      "content": {"type": "text"}
    }
  }'
```

## Обновить Content Type

**Endpoint:** `PUT /api/content-types/:uid`

**Тело запроса:** (аналогично созданию)

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/content-types/article \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "displayName": "Updated Article",
    "schema": {
      "title": {"type": "string"},
      "content": {"type": "text"},
      "author": {"type": "string"}
    }
  }'
```

## Удалить Content Type

**Endpoint:** `DELETE /api/content-types/:uid`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/content-types/article \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Схема (Schema)

Schema определяет структуру полей для записей Content Type. Пример:

```json
{
  "title": {
    "type": "string",
    "required": true
  },
  "slug": {
    "type": "string",
    "unique": true
  },
  "content": {
    "type": "text"
  },
  "published": {
    "type": "boolean",
    "default": false
  },
  "tags": {
    "type": "array"
  },
  "metadata": {
    "type": "object"
  }
}
```

**Поддерживаемые типы:**
- `string` - текстовая строка
- `text` - многострочный текст
- `number` - число
- `boolean` - булево значение
- `date` - дата
- `array` - массив
- `object` - объект/JSON

