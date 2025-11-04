# GitHub Pages Deployment

Настройка автоматической публикации документации на GitHub Pages с помощью Docsify.

## Автоматическая публикация

### 1. Создать GitHub Actions workflow

Создайте файл `.github/workflows/docs.yml`:

```yaml
name: Deploy Docs

on:
  push:
    branches:
      - main
    paths:
      - 'docs/**'
      - '.github/workflows/docs.yml'

permissions:
  contents: write

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
          cname: docs.xivercrm.com  # Опционально, для кастомного домена
```

### 2. Настройка GitHub Pages

1. Перейдите в Settings → Pages вашего репозитория
2. Выберите Source: **GitHub Actions**
3. Документация будет автоматически публиковаться при каждом push в `docs/`

## Ручная публикация

Если вы хотите публиковать вручную:

```bash
# Установите gh-pages (если используете npm)
npm install -g gh-pages

# Или используйте git subtree
git subtree push --prefix docs origin gh-pages
```

## Кастомный домен

1. Добавьте файл `docs/CNAME`:
```
docs.xivercrm.com
```

2. Настройте DNS записи для вашего домена:
   - Тип: CNAME
   - Имя: docs
   - Значение: your-username.github.io

3. GitHub автоматически настроит SSL сертификат

## Локальная разработка

Для локальной разработки документации:

```bash
# Установите docsify-cli
npm install -g docsify-cli

# Запустите локальный сервер
cd docs
docsify serve

# Документация будет доступна на http://localhost:3000
```

## Структура файлов

```
docs/
├── index.html          # Docsify конфигурация
├── README.md           # Главная страница
├── _sidebar.md         # Боковое меню
└── ...                 # Остальная документация
```

## Проверка перед публикацией

Перед коммитом проверьте:

1. Все ссылки работают
2. Изображения загружаются
3. Примеры кода корректны
4. Навигация работает правильно

```bash
# Локальная проверка
docsify serve docs
```

