# Production Deployment

Рекомендации по развертыванию XiverCRM в production.

## Подготовка

### 1. База данных

Используйте PostgreSQL:
- Настройте отдельного пользователя
- Ограничьте права доступа
- Настройте резервное копирование

### 2. Конфигурация

```env
GIN_MODE=release
DB_DRIVER=postgres
JWT_SECRET=your-strong-secret-key-min-32-chars
ALLOWED_ORIGINS=https://yourdomain.com
```

### 3. Сборка

```bash
# Backend
go build -o bin/xivercrm main.go

# Frontend
cd frontend && npm run build
```

## Развертывание

### Systemd Service

Создайте `/etc/systemd/system/xivercrm.service`:

```ini
[Unit]
Description=XiverCRM Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/xivercrm
ExecStart=/opt/xivercrm/bin/xivercrm
Restart=always
EnvironmentFile=/opt/xivercrm/.env

[Install]
WantedBy=multi-user.target
```

Запуск:

```bash
sudo systemctl enable xivercrm
sudo systemctl start xivercrm
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Мониторинг

- Логи приложения
- Health checks (`/health` endpoint)
- Мониторинг базы данных
- Мониторинг производительности

## Безопасность

1. Используйте HTTPS
2. Настройте firewall
3. Регулярно обновляйте зависимости
4. Используйте сильные пароли
5. Настройте rate limiting
6. Включите audit logs

## Резервное копирование

### База данных

```bash
# Ежедневный backup
pg_dump -U postgres xivercrm > backup_$(date +%Y%m%d).sql
```

### Медиа файлы

```bash
# Backup uploads
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz ./uploads
```

## Масштабирование

Для высокой нагрузки:
- Используйте load balancer
- Настройте несколько backend инстансов
- Используйте Redis для кеширования
- Оптимизируйте запросы к БД

