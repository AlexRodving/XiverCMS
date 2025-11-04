# Конфигурация базы данных

XiverCRM поддерживает SQLite и PostgreSQL.

## SQLite (по умолчанию)

Для разработки и небольших проектов:

```env
DB_DRIVER=sqlite
DB_PATH=./data/xivercrm.db
```

### Преимущества
- Не требует отдельного сервера
- Быстрая настройка
- Подходит для разработки

### Недостатки
- Не подходит для production с высокой нагрузкой
- Ограниченная поддержка параллельных запросов

## PostgreSQL (рекомендуется для production)

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=xivercrm
```

### Установка PostgreSQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql

# Docker
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```

### Создание базы данных

```bash
# Подключитесь к PostgreSQL
psql -U postgres

# Создайте базу данных
CREATE DATABASE xivercrm;

# Создайте пользователя (опционально)
CREATE USER xivercrm_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE xivercrm TO xivercrm_user;
```

## Миграции

Миграции выполняются автоматически при запуске приложения через GORM AutoMigrate.

Все модели автоматически создаются в базе данных:
- users
- roles
- permissions
- content_types
- content_entries
- media_files
- api_tokens
- audit_logs
- content_histories

## Подключение

XiverCRM автоматически подключается к базе данных при запуске. Убедитесь, что:

1. База данных существует
2. Пользователь имеет права доступа
3. Параметры подключения корректны в `.env`

## Резервное копирование

### PostgreSQL

```bash
# Backup
pg_dump -U postgres xivercrm > backup.sql

# Restore
psql -U postgres xivercrm < backup.sql
```

### SQLite

```bash
# Backup
cp ./data/xivercrm.db ./data/xivercrm.db.backup

# Restore
cp ./data/xivercrm.db.backup ./data/xivercrm.db
```

## Производительность

Для production рекомендуется:
- Использовать PostgreSQL
- Настроить connection pooling
- Регулярно делать резервные копии
- Мониторить производительность запросов

