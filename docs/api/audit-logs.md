# API: Audit Logs

Просмотр логов аудита для отслеживания действий пользователей.

## Получить список логов

**Endpoint:** `GET /api/audit-logs`

**Параметры:**
- `page` - номер страницы
- `pageSize` - размер страницы
- `action` - фильтр по действию (create, update, delete, login, etc.)
- `subject` - фильтр по объекту (user, content-entry, etc.)
- `userId` - фильтр по пользователю

**Пример:**
```bash
# Все логи
curl http://localhost:8080/api/audit-logs \
  -H "Authorization: Bearer YOUR_TOKEN"

# Фильтр по действию
curl "http://localhost:8080/api/audit-logs?action=create" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Фильтр по объекту
curl "http://localhost:8080/api/audit-logs?subject=content-entry" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Комбинированный фильтр
curl "http://localhost:8080/api/audit-logs?action=create&subject=content-entry&userId=1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "action": "create",
      "subject": "content-entry",
      "subjectId": 1,
      "description": "Created content entry",
      "ipAddress": "192.168.1.1",
      "userAgent": "Mozilla/5.0...",
      "user": {
        "id": 1,
        "email": "user@example.com"
      },
      "metadata": {
        "contentType": "article",
        "status": "published"
      },
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "pageSize": 50,
      "total": 100
    }
  }
}
```

## Типы действий

- `create` - создание объекта
- `update` - обновление объекта
- `delete` - удаление объекта
- `login` - вход в систему
- `publish` - публикация контента
- `unpublish` - снятие с публикации

## Типы объектов

- `user` - пользователь
- `role` - роль
- `permission` - право доступа
- `content-type` - тип контента
- `content-entry` - запись контента
- `media-file` - медиа файл
- `api-token` - API токен

## Метаданные

Каждый лог может содержать дополнительные метаданные в поле `metadata`:

```json
{
  "metadata": {
    "contentType": "article",
    "status": "published",
    "field": "title",
    "oldValue": "Old Title",
    "newValue": "New Title"
  }
}
```

## Использование

Audit logs автоматически создаются при:
- Логине пользователей
- Создании/обновлении/удалении контента
- Изменении пользователей и ролей
- Других важных действиях

Логи помогают отслеживать:
- Кто и когда выполнил действие
- С какого IP адреса
- Какие изменения были внесены
- Историю изменений контента

