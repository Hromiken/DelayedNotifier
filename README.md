# delayed-notifier

Сервис для **отложенной отправки уведомлений** (Telegram и Email).  
Принимает уведомления через REST API и отправляет их в заданное время.

---

## 🚀 Возможности

- Отложенная отправка уведомлений по времени `send_at`
- Поддержка временных зон (`ISO8601`)
- Каналы отправки:
    - Telegram Bot API
    - Email SMTP
- Очередь в PostgreSQL
- Фоновый воркер
- Docker Compose

---

## 🧱 Архитектура
POST /notify → PostgreSQL → Worker → Senders (Telegram / Email)



---

## ⚙️ Конфигурация `config.yaml`

```yaml
server:
  port: 8080

database:
  host: "postgres"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "notifier"

telegram:
  token: "123456789:ABC-EXAMPLE"

email:
  host: "smtp.yandex.ru"
  port: 465
  username: "yourmail@yandex.ru"
  password: "your_app_password"
  from: "yourmail@yandex.ru"

worker:
  interval_seconds: 3
  max_retries: 3
```

## 🐳 Docker Compose
```yaml
version: "3.8"

services:
  notifier:
    build: .
    container_name: delayed-notifier
    depends_on:
      - postgres
    volumes:
      - ./config.yaml:/app/config.yaml
    ports:
      - "8080:8080"

  postgres:
    image: postgres:15
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: notifier
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

Запуск:

`docker-compose up --build`

## 📡 API
Создать уведомление

POST /notify
Content-Type: application/json

### 🔔 Пример: Telegram уведомление
`{
"sender": "telegram",
"user_id": "123456789",
"message": "Test Moscow time",
"send_at": "2025-11-17T22:52:00+03:00"
}`

### ✉️ Пример: Email уведомление

`{
"sender": "email",
"user_id": "example@gmail.com",
"message": "Hello from delayed notifier!",
"send_at": "2025-11-17T22:52:00+03:00"
}`

### Ответ API
`{
"id": "009357a8-0e5e-4283-a544-b76c4c671c63"
}`

