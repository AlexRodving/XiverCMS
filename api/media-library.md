# API: Media Library

Управление медиа файлами (изображения, документы и т.д.).

## Загрузить файл

**Endpoint:** `POST /api/upload`

**Формат:** `multipart/form-data`

**Параметры:**
- `file` (required) - файл для загрузки
- `alternative` - альтернативный текст
- `caption` - подпись

**Ограничения:**
- Максимальный размер: 10MB
- Поддерживаемые форматы: любые

**Пример:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/image.jpg" \
  -F "alternative=Image description" \
  -F "caption=Image caption"
```

**Ответ:**
```json
{
  "id": 1,
  "name": "image.jpg",
  "alternative": "Image description",
  "caption": "Image caption",
  "width": 1920,
  "height": 1080,
  "format": "jpg",
  "mime": "image/jpeg",
  "size": 245678,
  "url": "/api/uploads/1234567890_image.jpg",
  "path": "./uploads/1234567890_image.jpg",
  "provider": "local",
  "createdById": 1,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

## Получить список файлов

**Endpoint:** `GET /api/media-files`

**Параметры запроса:**
- `page` - номер страницы
- `pageSize` - размер страницы
- `search` - поиск по имени или альтернативному тексту
- `mime` - фильтр по MIME типу (например: `image%`)

**Пример:**
```bash
# Все файлы
curl http://localhost:8080/api/media-files \
  -H "Authorization: Bearer YOUR_TOKEN"

# Только изображения
curl "http://localhost:8080/api/media-files?mime=image%" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Поиск
curl "http://localhost:8080/api/media-files?search=photo" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "image.jpg",
      "alternative": "Image description",
      "caption": "Image caption",
      "width": 1920,
      "height": 1080,
      "format": "jpg",
      "mime": "image/jpeg",
      "size": 245678,
      "url": "/api/uploads/1234567890_image.jpg",
      "provider": "local",
      "createdBy": {
        "id": 1,
        "email": "user@example.com"
      },
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "pageSize": 20,
      "total": 1
    }
  }
}
```

## Получить файл

**Endpoint:** `GET /api/media-files/:id`

**Пример:**
```bash
curl http://localhost:8080/api/media-files/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Обновить файл

**Endpoint:** `PUT /api/media-files/:id`

**Тело запроса:**
```json
{
  "name": "new-name.jpg",
  "alternative": "New alternative text",
  "caption": "New caption"
}
```

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/media-files/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "alternative": "Updated description",
    "caption": "Updated caption"
  }'
```

## Удалить файл

**Endpoint:** `DELETE /api/media-files/:id`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/media-files/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Доступ к файлам

Загруженные файлы доступны по URL:

```
http://localhost:8080/api/uploads/{filename}
```

Например:
```
http://localhost:8080/api/uploads/1234567890_image.jpg
```

## Использование в Content Entries

После загрузки файла, используйте его URL в данных Content Entry:

```json
{
  "data": {
    "title": "My Article",
    "image": "/api/uploads/1234567890_image.jpg"
  }
}
```

