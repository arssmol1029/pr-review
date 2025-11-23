# Сервис назначения ревьюеров для Pull Request’ов

Микросервис для автоматического назначения ревьюеров на Pull Request'ы и управления командами разработчиков.
Bзаимодействие — через HTTP API, спецификация лежит в openapi.yaml.

___

## Запуск проекта

Собрать и запустить сервис:
```bash
make run
```

Или через docker-compose:
```bash
docker-compose up -d --build
```

Сервис будет доступен по адресу: http://localhost:8080

## API Endpoints

### Команды

- POST /team/add - Создание команды
- GET /team/get?team_name={team_name} - Получение информации о команде

### Пользователи

- POST /users/setIsActive - Изменение активности пользователя
- GET /users/getReview?user_id={user_id} - Получение PR назначенных на пользователя

### Pull Requests

- POST /pullRequest/create - Создание PR
- POST /pullRequest/merge - Merge PR
- POST /pullRequest/reassign - Переназначение ревьюера

### Статистика

- GET /stats/total - Общая статистика
- GET /stats/team?team={teamName} - Статистика по команде
- GET /stats/user?user_id={user_id} - Статистика по пользователю

### База данных

Используется PostgreSQL со следующей схемой:

```sql
teams (name)
users (user_id, username, is_active, team_name)
pull_requests (id, name, author_id, status, created_at, merged_at)
pr_reviewers (pr_id, user_id)
```

## Команды

```bash
make run          # Сборка и запуск
make down         # Остановка сервисов
make logs         # Просмотр логов
make restart      # Перезапуск
make lint         # Запуск линтера
make clean        # Очистка
```

## Стек

- Язык: Go 1.24
- База данных: PostgreSQL 15
- Контейнеризация: Docker + Docker Compose
- Линтинг: golangci-lint

## Контакты

+ **Автор:** arssmol1029
+ **Telegram:** [@kepolb](https://t.me/kepolb)