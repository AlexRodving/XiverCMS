# Быстрый старт с документацией

## Локальная разработка документации

### Установка Docsify CLI

```bash
npm install -g docsify-cli
```

### Запуск локального сервера

```bash
cd docs
docsify serve

# Документация будет доступна на http://localhost:3000
```

### Просмотр документации

Откройте в браузере: `http://localhost:3000`

## Публикация на GitHub Pages

### Автоматическая публикация

Документация автоматически публикуется при push в `docs/` через GitHub Actions.

Workflow файл: `.github/workflows/docs.yml`

### Ручная публикация

```bash
# Используя gh-pages
npm install -g gh-pages
gh-pages -d docs

# Или используя git subtree
git subtree push --prefix docs origin gh-pages
```

## Структура документации

```
docs/
├── index.html              # Docsify конфигурация
├── README.md               # Главная страница
├── _sidebar.md            # Боковое меню
├── getting-started/       # Начало работы
├── api/                   # API документация
├── configuration/         # Конфигурация
├── security/              # Безопасность
├── development/           # Разработка
└── deployment/            # Развертывание
```

## Настройка кастомного домена

1. Создайте файл `docs/CNAME`:
```
docs.yourdomain.com
```

2. Настройте DNS:
   - Тип: CNAME
   - Имя: docs
   - Значение: your-username.github.io

3. GitHub автоматически настроит SSL

## Полезные ссылки

- [Docsify Documentation](https://docsify.js.org/)
- [GitHub Pages](https://pages.github.com/)
- [GitHub Actions](https://docs.github.com/en/actions)

