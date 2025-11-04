# API: Роли и права доступа

Управление ролями и правами доступа (RBAC).

## Роли

### Получить список ролей

**Endpoint:** `GET /api/roles`

**Пример:**
```bash
curl http://localhost:8080/api/roles \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Получить роль

**Endpoint:** `GET /api/roles/:id`

**Пример:**
```bash
curl http://localhost:8080/api/roles/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Создать роль

**Endpoint:** `POST /api/roles`

**Тело запроса:**
```json
{
  "name": "Editor",
  "description": "Can edit content",
  "type": "custom",
  "permissionIds": [1, 2, 3]
}
```

**Параметры:**
- `name` (required) - имя роли
- `description` - описание
- `type` - тип: `public` или `custom` (по умолчанию: `custom`)
- `permissionIds` - массив ID прав доступа

**Пример:**
```bash
curl -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Editor",
    "description": "Can edit content",
    "permissionIds": [1, 2]
  }'
```

### Обновить роль

**Endpoint:** `PUT /api/roles/:id`

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/roles/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description",
    "permissionIds": [1, 2, 3, 4]
  }'
```

### Удалить роль

**Endpoint:** `DELETE /api/roles/:id`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/roles/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Права доступа

### Получить список прав

**Endpoint:** `GET /api/permissions`

**Пример:**
```bash
curl http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Получить право

**Endpoint:** `GET /api/permissions/:id`

**Пример:**
```bash
curl http://localhost:8080/api/permissions/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Создать право

**Endpoint:** `POST /api/permissions`

**Тело запроса:**
```json
{
  "action": "read",
  "subject": "content-type:article",
  "properties": {},
  "conditions": {}
}
```

**Параметры:**
- `action` (required) - действие: `create`, `read`, `update`, `delete`, `publish`, `all`
- `subject` (required) - объект: `all`, `user`, `role`, `content-type`, `content-type:article`, etc.
- `properties` - дополнительные свойства (JSON)
- `conditions` - условия (JSON)

**Примеры:**

```bash
# Полный доступ ко всему
curl -X POST http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "all",
    "subject": "all"
  }'

# Чтение всех Content Types
curl -X POST http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "read",
    "subject": "content-type"
  }'

# Редактирование конкретного Content Type
curl -X POST http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "subject": "content-type:article"
  }'
```

### Обновить право

**Endpoint:** `PUT /api/permissions/:id`

**Пример:**
```bash
curl -X PUT http://localhost:8080/api/permissions/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "subject": "content-type:article"
  }'
```

### Удалить право

**Endpoint:** `DELETE /api/permissions/:id`

**Пример:**
```bash
curl -X DELETE http://localhost:8080/api/permissions/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Система прав

### Действия (Actions)

- `create` - создание
- `read` - чтение
- `update` - обновление
- `delete` - удаление
- `publish` - публикация
- `all` - все действия

### Объекты (Subjects)

- `all` - все объекты
- `user` - пользователи
- `role` - роли
- `content-type` - все Content Types
- `content-type:article` - конкретный Content Type
- `media-file` - медиа файлы
- `api-token` - API токены

### Примеры комбинаций

```json
// Чтение всех Content Types
{
  "action": "read",
  "subject": "content-type"
}

// Создание и редактирование статей
{
  "action": "create",
  "subject": "content-type:article"
}
{
  "action": "update",
  "subject": "content-type:article"
}

// Управление пользователями
{
  "action": "all",
  "subject": "user"
}
```

