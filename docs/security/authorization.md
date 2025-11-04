# Авторизация

Система авторизации XiverCMS основана на ролях и правах доступа.

## Super Admin

Пользователи с `isSuperAdmin: true` имеют полный доступ ко всем функциям без проверки прав.

## Роли

Каждый пользователь может иметь несколько ролей:
- Public - публичная роль
- Authenticated - для авторизованных пользователей
- Custom roles - созданные администратором

## Права доступа

Права определяют, что может делать пользователь:
- **Action** - действие (create, read, update, delete, publish, all)
- **Subject** - объект (user, role, content-type, media-file, etc.)

## Проверка прав

Права проверяются автоматически через RBAC middleware:
- При доступе к API endpoints
- В зависимости от ролей пользователя
- С учетом условий в permissions

## Примеры

### Разрешить чтение всех Content Types

```json
{
  "action": "read",
  "subject": "content-type"
}
```

### Разрешить редактирование конкретного типа

```json
{
  "action": "update",
  "subject": "content-type:article"
}
```

### Полный доступ

```json
{
  "action": "all",
  "subject": "all"
}
```

## Назначение прав

1. Создайте Permission через API
2. Создайте Role
3. Назначьте Permissions роли
4. Назначьте Role пользователю

Подробнее: [RBAC](rbac.md)

