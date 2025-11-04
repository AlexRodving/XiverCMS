# Конфигурация сервера

Настройка сервера XiverCMS.

## Порт

```env
PORT=8080
```

По умолчанию сервер запускается на порту 8080.

## Режим работы

```env
GIN_MODE=debug
```

**Возможные значения:**
- `debug` - режим разработки (подробные логи, паничные сообщения)
- `release` - production режим (минимальные логи, безопасность)

## Production рекомендации

1. Используйте `GIN_MODE=release`
2. Настройте reverse proxy (Nginx, Caddy)
3. Включите SSL/TLS
4. Настройте rate limiting
5. Используйте load balancer для высокой нагрузки

## Nginx пример конфигурации

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## SSL/TLS

Настройте SSL через Let's Encrypt или ваш сертификат:

```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        # ... остальные настройки
    }
}
```

