# Переменные окружения

Полный список переменных окружения для настройки XiverCRM.

## Server Configuration

### PORT
Порт для запуска сервера.

```env
PORT=8080
```

**По умолчанию:** `8080`

### GIN_MODE
Режим работы Gin framework.

```env
GIN_MODE=debug
```

**Возможные значения:**
- `debug` - режим разработки (подробные логи)
- `release` - production режим

**По умолчанию:** `debug`

## Database Configuration

### DB_DRIVER
Драйвер базы данных.

```env
DB_DRIVER=sqlite
```

**Возможные значения:**
- `sqlite` - SQLite (для разработки)
- `postgres` - PostgreSQL (для production)

**По умолчанию:** `sqlite`

### DB_HOST
Хост базы данных (для PostgreSQL).

```env
DB_HOST=localhost
```

**По умолчанию:** `localhost`

### DB_PORT
Порт базы данных (для PostgreSQL).

```env
DB_PORT=5432
```

**По умолчанию:** `5432`

### DB_USER
Пользователь базы данных (для PostgreSQL).

```env
DB_USER=postgres
```

**По умолчанию:** `postgres`

### DB_PASSWORD
Пароль базы данных (для PostgreSQL).

```env
DB_PASSWORD=postgres
```

**По умолчанию:** `postgres`

### DB_NAME
Имя базы данных (для PostgreSQL).

```env
DB_NAME=xivercrm
```

**По умолчанию:** `xivercrm`

### DB_PATH
Путь к файлу базы данных (для SQLite).

```env
DB_PATH=./data/xivercrm.db
```

**По умолчанию:** `./data/xivercrm.db`

## JWT Configuration

### JWT_SECRET
Секретный ключ для подписи JWT токенов.

```env
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

⚠️ **ВАЖНО:** Измените на уникальный секретный ключ в production!

**По умолчанию:** `change-this-secret-key-in-production`

### JWT_EXPIRATION
Время жизни JWT токенов.

```env
JWT_EXPIRATION=24h
```

**Формат:** Go duration (например: `24h`, `7d`, `30m`)

**По умолчанию:** `24h`

## CORS Configuration

### CORS_ORIGIN
Разрешенный origin для CORS.

```env
CORS_ORIGIN=http://localhost:5173
```

**По умолчанию:** `http://localhost:5173`

### ALLOWED_ORIGINS
Список разрешенных origins (через запятую).

```env
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000,https://example.com
```

**По умолчанию:** `http://localhost:5173,http://localhost:3000`

## Пример полного .env файла

```env
# Server Configuration
PORT=8080
GIN_MODE=release

# Database Configuration (PostgreSQL)
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=xivercrm

# Database Configuration (SQLite - альтернатива)
# DB_DRIVER=sqlite
# DB_PATH=./data/xivercrm.db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-min-32-characters
JWT_EXPIRATION=24h

# CORS Configuration
CORS_ORIGIN=https://yourdomain.com
ALLOWED_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com
```

## Production рекомендации

1. **Используйте PostgreSQL** вместо SQLite
2. **Измените JWT_SECRET** на уникальный ключ минимум 32 символа
3. **Установите GIN_MODE=release**
4. **Настройте ALLOWED_ORIGINS** только для ваших доменов
5. **Используйте переменные окружения** вместо .env файла в production
6. **Храните секреты** в безопасном месте (secrets manager)

