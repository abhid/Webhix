[![EN](https://img.shields.io/badge/lang-en-gray)](../CONTRIBUTING.md)

# Как контрибьютить в Webhix

Спасибо за интерес к проекту. Здесь всё что нужно чтобы начать.

## Что нужно установить

- Go 1.24+
- [sqlc](https://sqlc.dev) — для регенерации запросов к БД после изменений схемы
- [goose](https://github.com/pressly/goose) — для ручного запуска миграций
- Node.js 20+ — только если меняете UI

## Структура проекта

```
cmd/webhix/              точка входа (main.go)
internal/
  app/                   сборка приложения (зависимости, сервисы, HTTP)
  cli/                   CLI-команды (serve, forward)
  config/                конфигурация через переменные окружения
  core/                  бизнес-логика, интерфейсы репозиториев
  domain/                доменные типы и ошибки
  hub/                   SSE event hub
  server/                HTTP-хендлеры, роутинг, middleware
  store/                 SQLite-реализация (sqlc-сгенерированные запросы)
    migrations/          goose-миграции
    query/               сырые SQL-запросы
    sqlc/                сгенерированный Go-код — не редактировать вручную
  web/                   встроенные статические файлы и исходники UI
    ui/src/              TypeScript/CSS фронтенд
pkg/                     общие утилиты
docs/                    документация
config/                  примеры .env файлов
```

## Запуск локально

```sh
# Скопируйте пример конфига
cp config/.env.example config/.env

# Запустите сервер
go run ./cmd/webhix serve

# С кастомным адресом
go run ./cmd/webhix serve --addr :9090 --base-url http://localhost:9090
```

## Запуск всех проверок

```sh
make ci
```

Запускает форматирование, линтинг, vet и тесты. Все проверки должны проходить до открытия PR.

Отдельные команды:

```sh
make fmt        # форматирование Go кода
make lint       # запуск golangci-lint
make test       # запуск тестов
make web-check  # TypeScript проверка + ESLint + Prettier
```

## Ветки

Следуйте нотации из [docs/branch-patterns.md](branch-patterns.md).

```
feature/my-feature
fix/some-bug
refactor/cleanup-handler
```

- Ветвитесь от `main`
- Держите ветки короткими
- Одна фича или фикс на ветку

## Коммиты

Следуйте конвенции из [docs/commit-patterns.md](commit-patterns.md).

```
feat(server): add replay endpoint
fix(store): handle null body on insert
refactor(core): extract token generation
```

- Только английский язык
- Повелительное наклонение: `add`, `fix`, `remove`
- Без точки в конце
- Никаких `//nolint` комментариев — исправьте причину lint-ошибки

## Pull request-ы

- Заголовок PR следует тому же формату что и коммиты
- Держите PR маленькими и сфокусированными — одна вещь на PR
- Добавьте короткое описание что изменилось и почему
- Все CI-проверки должны проходить до мержа

## Архитектура

Проект следует строгой слоистой архитектуре. Зависимости идут только в одном направлении:

```
domain  ←  store
domain  ←  core
domain  ←  server
core    ←  server
```

- `domain` не зависит от других внутренних пакетов
- `core` определяет интерфейсы репозиториев — не импортирует `store`
- `server` зависит только от интерфейсов `core`, не от `store` напрямую
- SQL-ошибки не должны утекать в `server` — оборачивайте их в `domain`-ошибки на уровне `store`

## Изменения в БД

Если меняете схему:

1. Добавьте миграцию в `internal/store/migrations/` в формате goose
2. Обновите SQL-запросы в `internal/store/query/`
3. Регенерируйте: `sqlc generate`
4. Никогда не редактируйте файлы в `internal/store/sqlc/` вручную

## Изменения в UI

Фронтенд находится в `internal/web/ui/src/` — vanilla TypeScript, без фреймворков.

```sh
make web-dev    # запустить Vite dev server
make web-build  # собрать и встроить в бинарник
make web-check  # lint + type check + prettier
```

## Переменные окружения

| Переменная          | По умолчанию            | Описание                                   |
| ------------------- | ----------------------- | ------------------------------------------ |
| `WEBHIX_BASE_URL`   | `http://localhost:8080` | Публичный URL для генерации ссылок         |
| `WEBHIX_ADDR`       | `:8080`                 | Адрес для прослушивания                    |
| `WEBHIX_DB_PATH`    | `./data`                | Директория SQLite базы данных              |
| `WEBHIX_PASSWORD`   | —                       | Пароль Basic Auth                          |
| `WEBHIX_SECRET_KEY` | —                       | API секретный ключ (Bearer / X-Webhix-Key) |
