# API: Публичные эндпоинты

XiverCMS предоставляет публичные API эндпоинты для получения контента без аутентификации. Доступ к этим эндпоинтам контролируется через `accessType` в Content Type.

## Базовый URL

```
http://localhost:8080/api
```

## Типы доступа (AccessType)

Каждый Content Type может иметь один из следующих типов доступа:

- **`public`** - полностью публичный доступ, не требует аутентификации
- **`authenticated`** - требует авторизованного пользователя (любая роль)
- **`moderator`** - требует роль Moderator или Admin
- **`admin`** - требует роль Admin или Super Admin

## Получить список публичных Content Types

**Endpoint:** `GET /api/content-types`

⚠️ **Важно:** Этот эндпоинт возвращает только Content Types с `accessType = "public"` и `isVisible = true`.

Для получения всех Content Types (включая с ограниченным доступом) используйте защищенный эндпоинт `/api/content-types` с JWT токеном.

Возвращает только публичные (`accessType = "public"`) и видимые (`isVisible = true`) Content Types.

**Параметры:**
- `page` - номер страницы (по умолчанию: 1)
- `pageSize` - размер страницы (по умолчанию: 10)

**Пример:**
```bash
curl http://localhost:8080/api/content-types
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "uid": "article",
      "kind": "collectionType",
      "displayName": "Article",
      "description": "Blog articles",
      "isVisible": true,
      "accessType": "public",
      "schema": {
        "title": {"type": "string"},
        "content": {"type": "text"}
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
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

## Получить публичный Content Type

**Endpoint:** `GET /api/content-types/:uid`

⚠️ **Важно:** Этот эндпоинт возвращает только Content Types с `accessType = "public"` и `isVisible = true`.

Для получения Content Types с ограниченным доступом используйте защищенный эндпоинт `/api/content-types/:uid` с JWT токеном.

Возвращает Content Type только если он публичный и видимый.

**Пример:**
```bash
curl http://localhost:8080/api/content-types/article
```

**Ответ:**
```json
{
  "id": 1,
  "uid": "article",
  "kind": "collectionType",
  "displayName": "Article",
  "description": "Blog articles",
  "isVisible": true,
  "accessType": "public",
  "schema": {
    "title": {"type": "string", "required": true},
    "content": {"type": "text"}
  }
}
```

**Ошибки:**
- `404 Not Found` - если Content Type не найден или не является публичным

## Получить публичные записи Content Type

**Endpoint:** `GET /api/content-types/:uid/entries`

Возвращает только опубликованные (`status = "published"`) записи. Доступ контролируется через `accessType` Content Type.

**Параметры:**
- `page` - номер страницы (по умолчанию: 1)
- `pageSize` - размер страницы (по умолчанию: 10)
- `search` - поиск по содержимому записи
- `populate` - загрузить связанные записи (true/false)

**Доступ:**
- Если `accessType = "public"` - не требует аутентификации
- Если `accessType = "authenticated"` - требует JWT токен
- Если `accessType = "moderator"` - требует роль Moderator/Admin
- Если `accessType = "admin"` - требует роль Admin/Super Admin

**Пример (публичный доступ):**
```bash
curl http://localhost:8080/api/content-types/article/entries
```

**Пример (с аутентификацией):**
```bash
curl http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Пример (с поиском):**
```bash
curl "http://localhost:8080/api/content-types/article/entries?search=news"
```

**Пример (с populate):**
```bash
curl "http://localhost:8080/api/content-types/article/entries?populate=true"
```

**Ответ:**
```json
{
  "data": [
    {
      "id": 1,
      "contentTypeId": 1,
      "data": {
        "title": "My Article",
        "content": "Article content",
        "published": true
      },
      "status": "published",
      "publishedAt": "2024-01-01T00:00:00Z",
      "createdBy": {
        "id": 1,
        "email": "author@example.com"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
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

**Ошибки:**
- `404 Not Found` - если Content Type не найден
- `403 Forbidden` - если недостаточно прав доступа

## Получить публичную запись

**Endpoint:** `GET /api/content-types/:uid/entries/:id`

Возвращает запись только если она опубликована. Доступ контролируется через `accessType` Content Type.

**Параметры:**
- `populate` - загрузить связанные записи (true/false)

**Пример (публичный доступ):**
```bash
curl http://localhost:8080/api/content-types/article/entries/1
```

**Пример (с populate):**
```bash
curl "http://localhost:8080/api/content-types/article/entries/1?populate=true"
```

**Ответ:**
```json
{
  "id": 1,
  "contentTypeId": 1,
  "data": {
    "title": "My Article",
    "content": "Article content",
    "published": true
  },
  "status": "published",
  "publishedAt": "2024-01-01T00:00:00Z",
  "createdBy": {
    "id": 1,
    "email": "author@example.com"
  },
  "updatedBy": {
    "id": 1,
    "email": "author@example.com"
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Ошибки:**
- `404 Not Found` - если запись не найдена или не опубликована
- `403 Forbidden` - если недостаточно прав доступа

## Настройка доступа к Content Type

### В админке

1. Перейдите в **Content Types**
2. Выберите или создайте Content Type
3. В настройках выберите **Access Type**:
   - **Public** - доступно всем без аутентификации
   - **Authenticated** - требует авторизованного пользователя
   - **Moderator** - требует роль Moderator/Admin
   - **Admin** - требует роль Admin/Super Admin

### Через API

```bash
curl -X PUT http://localhost:8080/api/content-types/article \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "accessType": "public",
    "isVisible": true
  }'
```

## Примеры использования

### Получить все публичные статьи

```bash
curl http://localhost:8080/api/content-types/article/entries
```

### Получить статью по ID

```bash
curl http://localhost:8080/api/content-types/article/entries/1
```

### Поиск статей

```bash
curl "http://localhost:8080/api/content-types/article/entries?search=javascript"
```

### Пагинация

```bash
curl "http://localhost:8080/api/content-types/article/entries?page=2&pageSize=20"
```

### С аутентификацией (для authenticated/moderator/admin)

```bash
# Сначала получите JWT токен
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.jwt')

# Используйте токен для доступа
curl http://localhost:8080/api/content-types/article/entries \
  -H "Authorization: Bearer $TOKEN"
```

## Важные замечания

1. **Только опубликованные записи**: Публичные эндпоинты возвращают только записи со статусом `published`
2. **AccessType контроль**: Доступ контролируется на уровне Content Type через поле `accessType`
3. **JWT токены**: Для `authenticated`, `moderator`, `admin` требуется JWT токен в заголовке `Authorization: Bearer TOKEN`
4. **Роли**: Убедитесь, что пользователь имеет необходимую роль для доступа к Content Type с ограниченным доступом

## Сравнение с защищенными эндпоинтами

| Функция | Публичные эндпоинты | Защищенные эндпоинты |
|---------|---------------------|---------------------|
| Аутентификация | Опциональна (зависит от accessType) | Обязательна |
| Статус записей | Только `published` | Все статусы |
| HTTP методы | Только GET | GET, POST, PUT, DELETE |
| URL | `/api/content-types/:uid/entries` | `/api/content-types/:uid/entries` |
| Доступ | Контролируется accessType | Контролируется RBAC |

## Использование в JavaScript

```javascript
// Получить публичные статьи
fetch('http://localhost:8080/api/content-types/article/entries')
  .then(res => res.json())
  .then(data => {
    console.log(data.data); // Массив статей
    console.log(data.meta); // Метаданные пагинации
  });

// Получить статью по ID
fetch('http://localhost:8080/api/content-types/article/entries/1')
  .then(res => res.json())
  .then(article => {
    console.log(article);
  });

// С аутентификацией
const token = 'YOUR_JWT_TOKEN';
fetch('http://localhost:8080/api/content-types/article/entries', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
  .then(res => res.json())
  .then(data => console.log(data));
```

## Использование в Next.js / React

```jsx
// pages/articles.js или components/ArticlesList.js
import { useEffect, useState } from 'react';

export default function ArticlesList() {
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/content-types/article/entries')
      .then(res => res.json())
      .then(data => {
        setArticles(data.data);
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      {articles.map(article => (
        <div key={article.id}>
          <h2>{article.data.title}</h2>
          <p>{article.data.content}</p>
        </div>
      ))}
    </div>
  );
}
```
