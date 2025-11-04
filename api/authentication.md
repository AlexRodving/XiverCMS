# API: Аутентификация

## Логин

Вход в систему и получение JWT токена.

**Endpoint:** `POST /api/auth/login`

**Тело запроса:**
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

**Ответ:**
```json
{
  "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "user",
    "firstName": "John",
    "lastName": "Doe",
    "isSuperAdmin": false,
    "roles": [...]
  }
}
```

**Пример:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@xivercrm.com",
    "password": "admin123"
  }'
```

## Регистрация

Создание нового пользователя.

**Endpoint:** `POST /api/auth/register`

**Тело запроса:**
```json
{
  "email": "newuser@example.com",
  "username": "newuser",
  "password": "securepassword",
  "firstName": "Jane",
  "lastName": "Smith"
}
```

**Ответ:**
```json
{
  "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "email": "newuser@example.com",
    "username": "newuser",
    "firstName": "Jane",
    "lastName": "Smith"
  }
}
```

**Пример:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "username": "newuser",
    "password": "securepassword",
    "firstName": "Jane",
    "lastName": "Smith"
  }'
```

## Получение текущего пользователя

Получить информацию о текущем авторизованном пользователе.

**Endpoint:** `GET /api/auth/me`

**Заголовки:**
```http
Authorization: Bearer YOUR_JWT_TOKEN
```

**Ответ:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "user",
  "firstName": "John",
  "lastName": "Doe",
  "isSuperAdmin": false,
  "roles": [
    {
      "id": 1,
      "name": "Authenticated",
      "description": "Authenticated user role",
      "type": "public",
      "permissions": [...]
    }
  ]
}
```

**Пример:**
```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Использование JWT токена

После получения JWT токена, используйте его в заголовке `Authorization` для всех защищенных запросов:

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

Токен действителен в течение времени, указанного в `JWT_EXPIRATION` (по умолчанию 24 часа).

## Использование API Token

Альтернативный способ аутентификации - использование API токенов:

```http
Authorization: xvc_YOUR_API_TOKEN
```

API токены создаются через `/api/api-tokens` endpoint и могут иметь ограничения (read-only или full-access).

