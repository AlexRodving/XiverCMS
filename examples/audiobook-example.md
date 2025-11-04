# Пример: Аудиокнига с связями

Пример создания Content Type "Аудиокнига" с связями к авторам, сериям, циклам и лайкам.

## Структура Content Types

### 1. Author (Автор)

```json
{
  "uid": "author",
  "displayName": "Author",
  "kind": "collectionType",
  "schema": {
    "name": {
      "type": "string",
      "required": true
    },
    "bio": {
      "type": "text"
    },
    "photo": {
      "type": "media"
    }
  }
}
```

### 2. Series (Серия)

```json
{
  "uid": "series",
  "displayName": "Series",
  "kind": "collectionType",
  "schema": {
    "title": {
      "type": "string",
      "required": true
    },
    "description": {
      "type": "text"
    }
  }
}
```

### 3. Cycle (Цикл)

```json
{
  "uid": "cycle",
  "displayName": "Cycle",
  "kind": "collectionType",
  "schema": {
    "title": {
      "type": "string",
      "required": true
    },
    "description": {
      "type": "text"
    }
  }
}
```

### 4. Audiobook (Аудиокнига)

```json
{
  "uid": "audiobook",
  "displayName": "Audiobook",
  "kind": "collectionType",
  "schema": {
    "title": {
      "type": "string",
      "required": true
    },
    "duration": {
      "type": "number",
      "required": true,
      "description": "Duration in minutes"
    },
    "authors": {
      "type": "relation",
      "relationType": "manyToMany",
      "targetContentType": "author",
      "required": true
    },
    "series": {
      "type": "relation",
      "relationType": "manyToOne",
      "targetContentType": "series"
    },
    "cycle": {
      "type": "relation",
      "relationType": "manyToOne",
      "targetContentType": "cycle"
    },
    "likes": {
      "type": "number",
      "default": 0
    },
    "cover": {
      "type": "media"
    },
    "audioFile": {
      "type": "media",
      "required": true
    }
  }
}
```

## Создание через API

### 1. Создать Author Content Type

```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "author",
    "displayName": "Author",
    "kind": "collectionType",
    "schema": {
      "name": {"type": "string", "required": true},
      "bio": {"type": "text"},
      "photo": {"type": "media"}
    }
  }'
```

### 2. Создать Series Content Type

```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "series",
    "displayName": "Series",
    "kind": "collectionType",
    "schema": {
      "title": {"type": "string", "required": true},
      "description": {"type": "text"}
    }
  }'
```

### 3. Создать Cycle Content Type

```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "cycle",
    "displayName": "Cycle",
    "kind": "collectionType",
    "schema": {
      "title": {"type": "string", "required": true},
      "description": {"type": "text"}
    }
  }'
```

### 4. Создать Audiobook Content Type

```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "audiobook",
    "displayName": "Audiobook",
    "kind": "collectionType",
    "schema": {
      "title": {"type": "string", "required": true},
      "duration": {"type": "number", "required": true},
      "authors": {
        "type": "relation",
        "relationType": "manyToMany",
        "targetContentType": "author"
      },
      "series": {
        "type": "relation",
        "relationType": "manyToOne",
        "targetContentType": "series"
      },
      "cycle": {
        "type": "relation",
        "relationType": "manyToOne",
        "targetContentType": "cycle"
      },
      "likes": {"type": "number", "default": 0},
      "cover": {"type": "media"},
      "audioFile": {"type": "media", "required": true}
    }
  }'
```

## Создание записей

### 1. Создать автора

```bash
curl -X POST http://localhost:8080/api/content-types/author/entries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "name": "Джордж Оруэлл",
      "bio": "Английский писатель и публицист"
    },
    "status": "published"
  }'
```

### 2. Создать серию

```bash
curl -X POST http://localhost:8080/api/content-types/series/entries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "title": "Классика антиутопии",
      "description": "Серия классических антиутопических произведений"
    },
    "status": "published"
  }'
```

### 3. Создать аудиокнигу с связями

```bash
# Сначала создаем запись
curl -X POST http://localhost:8080/api/content-types/audiobook/entries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "title": "1984",
      "duration": 720,
      "likes": 0,
      "audioFile": "/api/uploads/1984.mp3",
      "cover": "/api/uploads/1984-cover.jpg"
    },
    "status": "published"
  }'

# Получаем ID созданной записи (например, 1)
# Затем создаем связи

# Связь с автором (ID автора = 1)
curl -X POST http://localhost:8080/api/content-types/audiobook/entries/1/relations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fieldName": "authors",
    "targetContentTypeUid": "author",
    "targetEntryId": 1,
    "relationType": "manyToMany"
  }'

# Связь с серией (ID серии = 1)
curl -X POST http://localhost:8080/api/content-types/audiobook/entries/1/relations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fieldName": "series",
    "targetContentTypeUid": "series",
    "targetEntryId": 1,
    "relationType": "manyToOne"
  }'
```

## Получение с связями

```bash
# Получить аудиокнигу со всеми связями
curl "http://localhost:8080/api/content-types/audiobook/entries/1?populate=true" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Ответ будет включать связанные записи:

```json
{
  "id": 1,
  "data": {
    "title": "1984",
    "duration": 720,
    "likes": 0,
    "authors": [
      {
        "id": 1,
        "data": {
          "name": "Джордж Оруэлл",
          "bio": "Английский писатель и публицист"
        }
      }
    ],
    "series": {
      "id": 1,
      "data": {
        "title": "Классика антиутопии",
        "description": "Серия классических антиутопических произведений"
      }
    }
  }
}
```

## Увеличение лайков

```bash
# Получить текущую запись
ENTRY=$(curl http://localhost:8080/api/content-types/audiobook/entries/1 \
  -H "Authorization: Bearer YOUR_TOKEN")

# Увеличить лайки
LIKES=$(echo $ENTRY | jq '.data.likes + 1')

curl -X PUT http://localhost:8080/api/content-types/audiobook/entries/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"data\": {
      \"title\": \"1984\",
      \"duration\": 720,
      \"likes\": $LIKES
    }
  }"
```

## Рейтинг на основе лайков

Для реализации рейтинга можно:

1. Использовать поле `likes` для подсчета
2. Создать отдельный Content Type "Rating" для более детального рейтинга
3. Использовать агрегацию при получении списка:

```bash
# Получить все аудиокниги отсортированные по лайкам
curl "http://localhost:8080/api/content-types/audiobook/entries?sort=likes:desc" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

