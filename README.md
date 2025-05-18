# 📬 Messaging App (Insider Case Study)

This is a messaging system designed to send unsent messages to a webhook every 2 minutes.
Built using modern Go architecture (Fiber v3, PostgreSQL, Redis, Zap) and fully dockerized.

---

## 🚀 Features

* ✅ Periodic dispatcher (every 2 mins) for sending pending messages
* ✅ Idempotent behavior: no duplicate sends
* ✅ Redis caching for sent messageId and sentAt
* ✅ REST API for starting/stopping dispatcher and listing sent messages
* ✅ PostgreSQL persistence
* ✅ Clean code structure with interfaces (repository, service, cache)
* ✅ OpenAPI 3 documentation
* ✅ Fully dockerized

---

## 💠 Tech Stack

| Component | Technology              |
| --------- | ----------------------- |
| Language  | Go 1.23                 |
| Web       | Fiber v3                |
| Database  | PostgreSQL 14           |
| Cache     | Redis 7                 |
| Logging   | Uber Zap                |
| Docs      | Swagger (OpenAPI 3)     |
| Container | Docker + Docker Compose |

---

## 📂 Project Structure

```
messaging-app/
├── cmd/                # Entry point (main.go)
├── internal/           # All domain logic
│   ├── cache/          # Redis adapter with interface
│   ├── config/         # (Optional) configs
│   ├── errors/         # Centralized error definitions
│   ├── handler/        # HTTP endpoints using Fiber v3
│   ├── logger/         # Zap logger setup
│   ├── model/          # DTOs / Entities
│   ├── repository/     # DB access (PostgreSQL)
│   ├── scheduler/      # Dispatcher (runs every 2 mins)
│   └── service/        # Business logic (messaging)
├── docs/               # Swagger YAML
├── init.sql            # DB bootstrap script
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
└── README.md
```

---

## ⚙️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/oguzhancagliyan/messaging-app.git
cd messaging-app
```

### 2. Start the Stack

```bash
docker-compose up --build
```

> This launches:
>
> * Go web service on port `8080`
> * PostgreSQL on port `5432`
> * Redis on port `6379`

Database is automatically initialized via `init.sql`

---

## 🧰 API Endpoints

| Method | Endpoint         | Description                    |
| ------ | ---------------- | ------------------------------ |
| POST   | `/start`         | Starts the 2-minute dispatcher |
| POST   | `/stop`          | Stops the dispatcher           |
| GET    | `/messages/sent` | Returns sent messages          |

---

## 📁 API Documentation

Swagger (OpenAPI 3) file is available at [`docs/swagger.yaml`](./docs/swagger.yaml).

To visualize:

* Open [https://editor.swagger.io](https://editor.swagger.io)
* Import the `swagger.yaml` file to view and test the API

---

## 💼 Database Info

| Key     | Value       |
| ------- | ----------- |
| Host    | `localhost` |
| Port    | `5432`      |
| DB Name | `messaging` |
| User    | `user`      |
| Pass    | `pass`      |

Sample data is inserted automatically on first run.

---

## 🔢 Redis Keys

Messages that are successfully sent will be cached like:

```
message:{messageId} → 2025-05-17T12:34:56Z
```

---

## 📅 Unit Testing

The following components are interface-based and easily mockable:

* `repository.MessageRepository`
* `cache.Cache`
* `service.MessageService`

---

## 📉 Notes

* `init.sql` only runs if the volume is clean. Use `docker-compose down -v` to reset.
* You can use PgAdmin or DBeaver to inspect the DB at `localhost:5432`
* This project requires Go 1.23+ due to `fiber/v2`

---
