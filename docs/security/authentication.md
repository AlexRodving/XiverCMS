# Аутентификация

XiverCMS поддерживает два типа аутентификации: JWT токены для пользователей и API токены для программного доступа.

## JWT Токены

### Генерация токена

При логине пользователь получает JWT токен:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }'
```

### Использование токена

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

### Время жизни

Настраивается через `JWT_EXPIRATION` (по умолчанию: 24 часа).

### Безопасность

- Токены подписываются секретным ключом (`JWT_SECRET`)
- Измените `JWT_SECRET` в production
- Используйте HTTPS для передачи токенов

## API Токены

### Создание токена

```bash
curl -X POST http://localhost:8080/api/api-tokens \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Token",
    "type": "read-only",
    "expiresAt": "2025-12-31T23:59:59Z"
  }'
```

### Использование токена

```http
Authorization: xvc_YOUR_API_TOKEN
```

### Типы токенов

- `read-only` - только GET запросы
- `full-access` - все HTTP методы

### Безопасность

- Токены имеют срок действия
- Read-only токены ограничены чтением
- Отслеживается последнее использование

## Хранение токенов

⚠️ **ВАЖНО:** Не храните токены в открытом виде!

- Используйте переменные окружения
- Не коммитьте токены в репозиторий
- Используйте secrets manager в production

