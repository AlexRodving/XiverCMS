# Установка

XiverCRM можно установить несколькими способами: из исходного кода, через Docker или скомпилированный бинарник.

## Предварительные требования

- **Go 1.21+** - для сборки из исходного кода
- **Node.js 18+** - для frontend (если нужен UI)
- **PostgreSQL 12+** (опционально) - для production использования
- **SQLite** (встроен) - для разработки

## Установка из исходного кода

### 1. Клонирование репозитория

```bash
git clone https://github.com/xivercrm/xivercrm.git
cd xivercrm
```

### 2. Установка зависимостей Backend

```bash
# Установка Go модулей
go mod download
go mod tidy
```

### 3. Настройка конфигурации

Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

Отредактируйте `.env` файл:

```env
# Server Configuration
PORT=8080
GIN_MODE=debug

# Database Configuration
DB_DRIVER=sqlite
DB_PATH=./data/xivercrm.db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h

# CORS Configuration
CORS_ORIGIN=http://localhost:5173
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

### 4. Запуск Backend

```bash
# Разработка
go run main.go

# Или через Makefile
make run
```

Backend будет доступен по адресу `http://localhost:8080`

### 5. Установка Frontend (опционально)

```bash
cd frontend
npm install
npm run dev
```

Frontend будет доступен по адресу `http://localhost:5173`

## Docker установка

### Использование Docker Compose

```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

Это запустит:
- PostgreSQL базу данных
- Backend сервер
- Frontend приложение

## Сборка для Production

### Backend

```bash
# Сборка бинарника
go build -o bin/xivercrm main.go

# Или через Makefile
make build

# Запуск
./bin/xivercrm
```

### Frontend

```bash
cd frontend
npm run build
```

Собранные файлы будут в `frontend/dist/`

## Проверка установки

После запуска сервера проверьте:

```bash
curl http://localhost:8080/health
```

Должен вернуться ответ:

```json
{
  "status": "ok"
}
```

## Первоначальная настройка

После первого запуска создается администратор:

- **Email**: `admin@xivercrm.com`
- **Password**: `admin123`

⚠️ **ВАЖНО**: Смените пароль администратора в production!

## Следующие шаги

- [Быстрый старт](quickstart.md) - начните использовать XiverCRM
- [Конфигурация](../configuration/overview.md) - настройте систему под ваши нужды

