# Связи между Content Types

XiverCMS поддерживает связи между Content Types, аналогично Strapi.

## Типы связей

Поддерживаются все стандартные типы связей:

- **One to One** - одна запись связана с одной записью
- **One to Many** - одна запись связана со многими записями
- **Many to One** - много записей связаны с одной записью
- **Many to Many** - много записей связаны со многими записями

## Создание связи

### Через визуальный редактор

1. При создании Content Type добавьте поле типа **Relation**
2. Выберите **Relation Type** (One to One, One to Many, Many to One, Many to Many)
3. Выберите **Target Content Type** (целевой Content Type)
4. Сохраните Content Type

### Через JSON схему

```json
{
  "authors": {
    "type": "relation",
    "relationType": "manyToMany",
    "targetContentType": "author"
  },
  "series": {
    "type": "relation",
    "relationType": "manyToOne",
    "targetContentType": "series"
  }
}
```

## Использование связей

### Создание записи со связями

При создании записи передайте ID связанных записей:

```json
{
  "data": {
    "title": "1984",
    "duration": 720,
    "authors": [1, 2],    // Массив ID для manyToMany
    "series": 1,          // Один ID для manyToOne
    "cycle": 1
  },
  "status": "published"
}
```

Связи создаются автоматически при создании/обновлении записи.

### Получение записей со связями

Используйте параметр `populate=true`:

```bash
curl "http://localhost:8080/api/content-types/audiobook/entries/1?populate=true" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Ответ будет включать связанные записи:

```json
{
  "id": 1,
  "data": {
    "title": "1984",
    "duration": 720
  },
  "authors": [
    {
      "id": 1,
      "data": {
        "name": "Джордж Оруэлл"
      }
    }
  ],
  "series": {
    "id": 1,
    "data": {
      "title": "Классика антиутопии"
    }
  }
}
```

## API для управления связями

### Получить связи записи

```bash
GET /api/content-types/:uid/entries/:id/relations
```

### Создать связь

```bash
POST /api/content-types/:uid/entries/:id/relations
{
  "fieldName": "authors",
  "targetContentTypeUid": "author",
  "targetEntryId": 1,
  "relationType": "manyToMany"
}
```

### Удалить связь

```bash
DELETE /api/content-types/:uid/entries/:id/relations/:relationId
```

### Получить связанные записи по полю

```bash
GET /api/content-types/:uid/entries/:id/relations/:field
```

## Примеры использования

### Аудиокнига с авторами, серией и циклом

См. [Полный пример](../examples/audiobook-example.md)

### Создание Content Type с связями

```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "audiobook",
    "displayName": "Audiobook",
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
      "likes": {"type": "number", "default": 0}
    }
  }'
```

### Создание записи со связями

```bash
curl -X POST http://localhost:8080/api/content-types/audiobook/entries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "title": "1984",
      "duration": 720,
      "likes": 0,
      "authors": [1, 2],
      "series": 1
    },
    "status": "published"
  }'
```

## Автоматическая обработка

Связи обрабатываются автоматически:

1. При создании записи - связи создаются в таблице `ContentRelation`
2. При обновлении записи - связи обновляются
3. При удалении записи - связи удаляются (cascade)

## Ограничения

- Связи работают только между Content Types
- Целевой Content Type должен существовать
- Целевые записи должны существовать при создании связи

## Рекомендации

1. **Создавайте Content Types сначала** - создайте все необходимые Content Types
2. **Используйте визуальный редактор** - проще настраивать связи через UI
3. **Используйте populate** - для получения связанных данных используйте `populate=true`
4. **Проверяйте связи** - используйте API для проверки созданных связей

