# RBAC (Role-Based Access Control)

Детальная система управления правами доступа.

## Концепция

RBAC позволяет точно контролировать, кто и что может делать в системе.

## Компоненты

### Roles (Роли)
Группы прав доступа, назначаемые пользователям.

### Permissions (Права)
Конкретные разрешения на действия с объектами.

## Структура Permission

```json
{
  "action": "update",
  "subject": "content-type:article",
  "properties": {},
  "conditions": {}
}
```

### Actions (Действия)

- `create` - создание
- `read` - чтение
- `update` - обновление
- `delete` - удаление
- `publish` - публикация
- `all` - все действия

### Subjects (Объекты)

- `all` - все объекты
- `user` - пользователи
- `role` - роли
- `content-type` - все Content Types
- `content-type:article` - конкретный Content Type
- `media-file` - медиа файлы
- `api-token` - API токены

## Создание роли с правами

```bash
# 1. Создать Permission
curl -X POST http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "read",
    "subject": "content-type:article"
  }'

# 2. Создать Role с Permission
curl -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Editor",
    "description": "Can edit articles",
    "permissionIds": [1]
  }'

# 3. Назначить Role пользователю
curl -X PUT http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roleIds": [1]
  }'
```

## Примеры ролей

### Content Editor

```json
{
  "name": "Content Editor",
  "permissions": [
    {"action": "read", "subject": "content-type"},
    {"action": "create", "subject": "content-entry"},
    {"action": "update", "subject": "content-entry"},
    {"action": "publish", "subject": "content-entry"}
  ]
}
```

### Media Manager

```json
{
  "name": "Media Manager",
  "permissions": [
    {"action": "all", "subject": "media-file"}
  ]
}
```

## Проверка прав

Права проверяются автоматически через RBAC middleware при каждом запросе к API.

Пользователь должен иметь хотя бы одну роль с соответствующим permission.

## Super Admin

Пользователи с `isSuperAdmin: true` обходят все проверки прав и имеют полный доступ.

