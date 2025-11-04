# Конфигурация базы данных

XiverCMS поддерживает SQLite и PostgreSQL.

## SQLite (по умолчанию)

Для разработки и небольших проектов:

```env
DB_DRIVER=sqlite
DB_PATH=./data/xivercms.db
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
DB_NAME=xivercms
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
CREATE DATABASE xivercms;

# Создайте пользователя (опционально)
CREATE USER xivercms_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE xivercms TO xivercms_user;
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

XiverCMS автоматически подключается к базе данных при запуске. Убедитесь, что:

1. База данных существует
2. Пользователь имеет права доступа
3. Параметры подключения корректны в `.env`

## Резервное копирование

### PostgreSQL

```bash
# Backup
pg_dump -U postgres xivercms > backup.sql

# Restore
psql -U postgres xivercms < backup.sql
```

### SQLite

```bash
# Backup
cp ./data/xivercms.db ./data/xivercms.db.backup

# Restore
cp ./data/xivercms.db.backup ./data/xivercms.db
```

## Производительность

Для production рекомендуется:
- Использовать PostgreSQL
- Настроить connection pooling
- Регулярно делать резервные копии
- Мониторить производительность запросов

