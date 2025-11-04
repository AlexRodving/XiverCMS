# API: Пользователи

Управление пользователями системы.

## Получить список пользователей

**Endpoint:** `GET /api/users`

**Параметры:**
- `page` - номер страницы
- `pageSize` - размер страницы
- `search` - поиск по email, username, имени

**Пример:**
```bash
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Получить пользователя

**Endpoint:** `GET /api/users/:id`

**Пример:**
```bash
curl http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Обновить пользователя

**Endpoint:** `PUT /api/users/:id`

**Тело запроса:**
```json
{
  "email": "newemail@example.com",
  "username": "newusername",
  "firstName": "John",
  "lastName": "Doe",
  "isActive": true,
  "roleIds": [1, 2]
}
```

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Jane",
    "lastName": "Smith",
    "roleIds": [1]
  }'
```

## Изменить пароль

**Endpoint:** `PUT /api/users/:id/password`

**Тело запроса:**
```json
{
  "password": "newpassword"
}
```

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/users/2/password \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "newsecurepassword"
  }'
```

## Удалить пользователя

**Endpoint:** `DELETE /api/users/:id`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

