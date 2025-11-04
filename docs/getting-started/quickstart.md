# Быстрый старт

Это руководство поможет вам быстро начать работу с XiverCRM.

## 1. Вход в систему

Откройте админ-панель в браузере: `http://localhost:5173`

Войдите с учетными данными администратора:
- Email: `admin@xivercrm.com`
- Password: `admin123`

## 2. Создание Content Type

### Через API

```bash
curl -X POST http://localhost:8080/api/content-types \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
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
      },
      "author": {
        "type": "string"
      }
    }
  }'
```

### Через админ-панель

1. Перейдите в раздел "Content Types"
2. Нажмите "Create Content Type"
3. Заполните форму:
   - UID: `article`
   - Display Name: `Article`
   - Schema: JSON с полями
4. Нажмите "Create"

## 3. Создание записи

```bash
curl -X POST http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "title": "My First Article",
      "slug": "my-first-article",
      "content": "This is the content of my first article.",
      "published": true,
      "author": "John Doe"
    },
    "status": "published"
  }'
```

## 4. Получение записей

```bash
# Получить все записи
curl http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Получить конкретную запись
curl http://localhost:8080/api/content-types/article/entries/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 5. Загрузка медиа файла

```bash
curl -X POST http://localhost:8080/api/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/image.jpg" \
  -F "alternative=Image description" \
  -F "caption=Image caption"
```

## 6. Создание API Token

```bash
curl -X POST http://localhost:8080/api/api-tokens \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API Token",
    "description": "Token for external service",
    "type": "read-only",
    "expiresAt": "2025-12-31T23:59:59Z"
  }'
```

Используйте полученный токен для доступа к API:

```bash
curl http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: xvc_YOUR_API_TOKEN"
```

## 7. Управление пользователями

### Создание пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "newuser",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Назначение роли

```bash
curl -X PUT http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roleIds": [1, 2]
  }'
```

## 8. Просмотр Audit Logs

```bash
curl http://localhost:8080/api/audit-logs \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# С фильтрами
curl "http://localhost:8080/api/audit-logs?action=create&subject=content-entry" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Следующие шаги

- [API Документация](../api/overview.md) - полная документация API
- [Конфигурация](../configuration/overview.md) - настройка системы
- [Безопасность](../security/overview.md) - безопасность и права доступа

