# Webhix

[![EN](https://img.shields.io/badge/lang-en-gray)](../README.md)

Self-hosted инспектор вебхуков. Один бинарник, SQLite, никаких внешних зависимостей.

webhook.site удобен, но отправляет все данные на чужой сервер. Stripe payload-ы, OAuth токены, персональные данные - всё это уходит из вашей сети. Многие компании блокируют его по этой причине. Webhix работает на вашей инфраструктуре, хранит всё локально и не мешает работе.

## Что умеет

Создаёте endpoint, направляете на него вебхуки и смотрите что приходит. Каждый запрос сохраняется полностью - заголовки, тело, query параметры, IP, время, content-type, размер. Интерфейс обновляется без перезагрузки страницы.

Помимо просмотра:

- **Replay** - повторить любой запрос одним кликом, можно с изменениями
- **Кастомные ответы** - настроить статус код, заголовки и тело ответа вашего endpoint-а (удобно как лёгкий mock-сервер)
- **Forwarding через CLI** - проксировать входящие запросы на локальный порт: `webhix forward <token> --to localhost:3000`
- **Экспорт** - скопировать любой запрос как готовую команду curl или HTTPie

## Быстрый старт

### Бинарник

```sh
curl -fsSL https://webhix.dev/install.sh | sh
webhix serve --base-url https://hooks.yourdomain.com
```

### Docker

```sh
docker run -p 8080:8080 -v webhix-data:/data \
  -e WEBHIX_BASE_URL=https://hooks.yourdomain.com \
  ghcr.io/gaisbax/webhix
```

### Docker Compose

```yaml
services:
  webhix:
    image: ghcr.io/gaisbax/webhix
    ports: ["8080:8080"]
    volumes: ["./data:/data"]
    environment:
      WEBHIX_BASE_URL: https://hooks.yourdomain.com
```

### Локальная разработка (без домена)

```sh
webhix serve
# Listening on http://localhost:8080
```

URL endpoint-ов формируется по шаблону `https://<base-url>/r/<token>`.

## Авторизация

По умолчанию однопользовательский режим. Пароль задаётся через env:

```sh
WEBHIX_PASSWORD=yourpassword webhix serve
```

## Обратный прокси

Работает за Caddy, Nginx, Traefik. Автоматически читает заголовки `X-Forwarded-*`. Укажите `WEBHIX_BASE_URL` чтобы совпадал с вашим публичным доменом.

## Конфигурация

| Env переменная | По умолчанию | Описание |
| -------------- | ------------ | -------- |
| `WEBHIX_BASE_URL` | `http://localhost:8080` | Публичный URL для генерации ссылок на endpoint-ы |
| `WEBHIX_ADDR` | `:8080` | Адрес для прослушивания (например `0.0.0.0:9000`) |
| `WEBHIX_DB_PATH` | `./data` | Путь к директории с SQLite базой данных |

## Технические детали

- Написан на Go, собирается в один бинарник
- SQLite по умолчанию, внешняя база не нужна
- UI встроен в бинарник через `go:embed`
- Работает на Linux, macOS, Windows (amd64 + arm64)
- Потребление памяти в простое - менее 50 МБ

## Roadmap

### v0.2

- Мультипользовательский режим с базовым RBAC
- Верификация подписи вебхуков (в стиле Stripe, GitHub)
- Валидация схемы
- Уведомления о новых запросах (Slack, Telegram, Discord)
- Опциональная поддержка Postgres
- Авто-HTTPS через Let's Encrypt (без обратного прокси)

### v0.3+

- Tunnel режим - подключение к управляемому relay для получения публичного URL без сервера

## Лицензия

[AGPL-3.0](../LICENSE). Self-hosted использование всегда бесплатно и открыто.

Если хотите запустить Webhix как сетевой сервис с закрытыми изменениями - свяжитесь с нами по вопросу коммерческой лицензии.
