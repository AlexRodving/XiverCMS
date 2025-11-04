# API: API Tokens

Управление API токенами для программного доступа.

## Получить список токенов

**Endpoint:** `GET /api/api-tokens`

**Параметры:**
- `page` - номер страницы
- `pageSize` - размер страницы

**Пример:**
```bash
curl http://localhost:8080/api/api-tokens \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "My API Token",
      "description": "Token for external service",
      "type": "read-only",
      "token": "xvc_abc123...",  // Маскирован
      "lastUsedAt": "2024-01-01T12:00:00Z",
      "expiresAt": "2025-12-31T23:59:59Z",
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
      "pageSize": 10,
      "total": 1
    }
  }
}
```

## Получить токен

**Endpoint:** `GET /api/api-tokens/:id`

**Примечание:** Токен возвращается маскированным (первые 8 символов).

**Пример:**
```bash
curl http://localhost:8080/api/api-tokens/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Создать токен

**Endpoint:** `POST /api/api-tokens`

**Тело запроса:**
```json
{
  "name": "My API Token",
  "description": "Token for external service",
  "type": "read-only",
  "expiresAt": "2025-12-31T23:59:59Z"
}
```

**Параметры:**
- `name` (required) - имя токена
- `description` - описание
- `type` - тип: `read-only` или `full-access` (по умолчанию: `read-only`)
- `expiresAt` - дата истечения (опционально)

**Типы токенов:**
- `read-only` - только чтение (GET запросы)
- `full-access` - полный доступ (все методы)

**Пример:**
```bash
curl -X POST http://localhost:8080/api/api-tokens \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API Token",
    "description": "Token for external service",
    "type": "read-only",
    "expiresAt": "2025-12-31T23:59:59Z"
  }'
```

**Ответ:**
```json
{
  "id": 1,
  "name": "My API Token",
  "description": "Token for external service",
  "type": "read-only",
  "token": "xvc_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",  // Полный токен только при создании
  "expiresAt": "2025-12-31T23:59:59Z",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

⚠️ **ВАЖНО:** Сохраните токен сразу после создания! Он показывается только один раз.

## Обновить токен

**Endpoint:** `PUT /api/api-tokens/:id`

**Тело запроса:**
```json
{
  "name": "Updated Token Name",
  "description": "Updated description",
  "type": "full-access",
  "expiresAt": "2026-12-31T23:59:59Z"
}
```

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/api-tokens/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description",
    "expiresAt": "2026-12-31T23:59:59Z"
  }'
```

## Удалить токен

**Endpoint:** `DELETE /api/api-tokens/:id`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/api-tokens/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Использование API токена

После создания токена, используйте его в заголовке `Authorization`:

```http
Authorization: xvc_YOUR_API_TOKEN
```

**Пример:**
```bash
# Создание токена
TOKEN=$(curl -X POST http://localhost:8080/api/api-tokens \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Token",
    "type": "read-only"
  }' | jq -r '.token')

# Использование токена
curl http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: $TOKEN"
```

## Ограничения

- **Read-only токены** могут выполнять только GET запросы
- **Full-access токены** могут выполнять все HTTP методы
- Токены с истекшим сроком действия не работают
- При каждом использовании обновляется `lastUsedAt`

## Безопасность

- Токены должны храниться в безопасном месте
- Не коммитьте токены в репозиторий
- Используйте переменные окружения для хранения токенов
- Регулярно ротируйте токены
- Удаляйте неиспользуемые токены

