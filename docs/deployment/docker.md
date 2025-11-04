# Docker Deployment

Развертывание XiverCMS с помощью Docker.

## Docker Compose

Используйте готовый `docker-compose.yml`:

```bash
docker-compose up -d
```

Это запустит:
- PostgreSQL базу данных
- Backend сервер
- Frontend приложение

## Отдельные контейнеры

### Backend

```bash
docker build -f Dockerfile.backend -t xivercms-backend .
docker run -p 8080:8080 xivercms-backend
```

### Frontend

```bash
cd frontend
docker build -t xivercms-frontend .
docker run -p 80:80 xivercms-frontend
```

## Переменные окружения

Используйте `.env` файл или передавайте через `-e`:

```bash
docker run -e PORT=8080 -e DB_DRIVER=postgres xivercms-backend
```

## Volumes

Для постоянного хранения данных:

```bash
docker run -v ./data:/app/data xivercms-backend
```

## Production рекомендации

1. Используйте Docker Compose для оркестрации
2. Настройте volumes для данных
3. Используйте secrets для паролей
4. Настройте health checks
5. Используйте reverse proxy (Nginx)

