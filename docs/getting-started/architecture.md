# Архитектура

XiverCRM построен на современной архитектуре с разделением на backend и frontend.

## Общая архитектура

```
┌─────────────────┐
│   Frontend      │  Vue3 + Tailwind CSS
│   (Vue3)        │  Порт: 5173
└────────┬────────┘
         │ HTTP/REST API
         │
┌────────▼────────┐
│   Backend       │  Go + Gin + GORM
│   (Go)          │  Порт: 8080
└────────┬────────┘
         │
┌────────▼────────┐
│   Database      │  PostgreSQL / SQLite
└─────────────────┘
```

## Backend архитектура

### Структура проекта

```
XiverCRM/
├── config/          # Конфигурация приложения
├── models/          # GORM модели данных
├── database/        # Подключение и миграции БД
├── handlers/        # HTTP handlers (контроллеры)
├── middleware/      # Middleware (auth, CORS, RBAC)
├── routes/          # Маршрутизация API
├── auth/            # Аутентификация и JWT
└── main.go          # Точка входа
```

### Поток запроса

```
HTTP Request
    │
    ├─► CORS Middleware
    │
    ├─► API Token Middleware (опционально)
    │
    ├─► Auth Middleware (JWT)
    │
    ├─► RBAC Middleware (проверка прав)
    │
    ├─► Handler (обработка запроса)
    │
    ├─► Database (GORM)
    │
    └─► HTTP Response
```

### Модели данных

- **User** - пользователи системы
- **Role** - роли пользователей
- **Permission** - права доступа
- **ContentType** - типы контента
- **ContentEntry** - записи контента
- **MediaFile** - медиа файлы
- **APIToken** - API токены
- **AuditLog** - логи аудита
- **ContentHistory** - история изменений

## Frontend архитектура

### Структура проекта

```
frontend/
├── src/
│   ├── api/         # API клиенты
│   ├── views/       # Страницы (Vue компоненты)
│   ├── layouts/     # Layout компоненты
│   ├── stores/      # Pinia stores (state management)
│   ├── router/      # Vue Router
│   └── App.vue      # Корневой компонент
├── public/          # Статические файлы
└── package.json
```

### State Management

Используется Pinia для управления состоянием:

- **auth store** - состояние аутентификации
- Другие stores для различных модулей

## База данных

### Поддерживаемые БД

- **SQLite** - для разработки (по умолчанию)
- **PostgreSQL** - для production

### Миграции

Автоматические миграции через GORM AutoMigrate при запуске приложения.

## Безопасность

### Аутентификация

- JWT токены для пользователей
- API токены для программного доступа
- bcrypt для хеширования паролей

### Авторизация

- Role-Based Access Control (RBAC)
- Проверка прав через middleware
- Super Admin имеет все права

## API

RESTful API со следующими возможностями:

- Пагинация
- Фильтрация
- Поиск
- Сортировка

## Deployment

### Docker

- Многостадийная сборка
- Docker Compose для оркестрации
- Nginx для frontend

### Production

- Компиляция Go в бинарник
- Сборка frontend через Vite
- Настройка через переменные окружения

