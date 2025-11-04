# Структура проекта

Подробное описание структуры проекта XiverCMS.

## Корневая структура

```
XiverCMS/
├── config/              # Конфигурация
├── models/              # Модели данных
├── database/            # Работа с БД
├── handlers/            # HTTP handlers
├── middleware/          # Middleware
├── routes/              # Маршруты
├── auth/                # Аутентификация
├── frontend/            # Frontend приложение
├── docs/                # Документация
├── main.go              # Точка входа
├── go.mod               # Go модули
├── Makefile             # Команды сборки
├── docker-compose.yml   # Docker Compose
└── README.md            # Основной README
```

## Backend структура

### config/

Конфигурация приложения.

```
config/
└── config.go            # Загрузка конфигурации из .env
```

### models/

GORM модели данных.

```
models/
├── user.go              # User, Role, Permission
├── content.go           # ContentType, ContentEntry, MediaFile
├── token.go             # APIToken
└── audit.go             # AuditLog, ContentHistory
```

### database/

Работа с базой данных.

```
database/
└── database.go          # Подключение, миграции, seed
```

### handlers/

HTTP handlers (контроллеры).

```
handlers/
├── auth_handler.go      # Аутентификация
├── user_handler.go     # Управление пользователями
├── content_handler.go   # Content Types и Entries
├── media_handler.go     # Media Library
├── role_handler.go      # Роли и права
├── token_handler.go     # API Tokens
└── audit_handler.go    # Audit Logs и History
```

### middleware/

Middleware для обработки запросов.

```
middleware/
├── auth.go              # JWT аутентификация
├── api_token.go         # API Token аутентификация
├── cors.go              # CORS настройки
└── rbac.go              # Проверка прав доступа
```

### routes/

Маршрутизация API.

```
routes/
└── routes.go            # Настройка всех маршрутов
```

### auth/

Аутентификация и JWT.

```
auth/
└── auth.go              # JWT генерация и валидация
```

## Frontend структура

### frontend/src/

```
src/
├── api/                 # API клиенты
│   ├── client.js        # Axios клиент
│   ├── auth.js          # API аутентификации
│   ├── users.js         # API пользователей
│   └── content.js       # API контента
├── views/               # Страницы
│   ├── Login.vue
│   ├── Dashboard.vue
│   ├── ContentTypes.vue
│   └── ...
├── layouts/             # Layout компоненты
│   └── DashboardLayout.vue
├── stores/              # Pinia stores
│   └── auth.js
├── router/              # Vue Router
│   └── index.js
├── App.vue              # Корневой компонент
└── main.js              # Точка входа
```

## Документация

### docs/

```
docs/
├── index.html           # Docsify конфигурация
├── README.md            # Главная страница
├── _sidebar.md          # Боковое меню
├── getting-started/     # Начало работы
├── api/                 # API документация
├── configuration/       # Конфигурация
├── security/            # Безопасность
├── development/         # Разработка
└── deployment/          # Развертывание
```

## Конфигурационные файлы

- `.env.example` - пример конфигурации
- `.gitignore` - игнорируемые файлы
- `go.mod` - Go зависимости
- `package.json` - Node.js зависимости
- `Makefile` - команды сборки
- `docker-compose.yml` - Docker Compose
- `Dockerfile.backend` - Docker образ backend
- `frontend/Dockerfile` - Docker образ frontend

