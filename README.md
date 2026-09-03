# Cascade Pro API (Api Golang Server)

Backend API for the Cascade Pro platform, built with Go and designed around a feature-oriented architecture.

_\*README was made with AI_

## Features

- Authentication and authorization
- JWT-based security
- User management
- Session management with Redis
- Client management
- Request management
- Media upload and processing
- Image processing and WebP conversion
- S3-compatible object storage
- PostgreSQL persistence
- MongoDB integration
- Redis caching and session storage
- IP geolocation via [IPinfo](https://ipinfo.io/)
- Background workers
- Transactional Outbox pattern
- Multi-instance-safe Outbox event claiming
- Structured logging with Zap
- Request validation
- Swagger / OpenAPI documentation
- Graceful application shutdown
- Docker-based development environment
- Database migrations with `golang-migrate`

---

## Tech Stack

| Technology               | Purpose                     |
| ------------------------ | --------------------------- |
| **Go 1.25.10**           | Backend runtime             |
| **PostgreSQL 16**        | Primary relational database |
| **Redis 5**              | Cache and session storage   |
| **MongoDB 8**            | Document-oriented data      |
| **S3**                   | Object and media storage    |
| **pgx/v5**               | PostgreSQL driver           |
| **go-redis/v9**          | Redis client                |
| **MongoDB Go Driver v2** | MongoDB client              |
| **AWS SDK for Go v2**    | S3 integration              |
| **JWT**                  | Authentication              |
| **Zap**                  | Structured logging          |
| **Swaggo**               | Swagger/OpenAPI generation  |
| **Docker Compose**       | Local infrastructure        |
| **golang-migrate**       | Database migrations         |

---

## Architecture

The application follows a feature-oriented architecture with clear separation between business features, domain logic, infrastructure, and transport.

```text
.
├── cmd/
│   └── main.go
│
├── internal/
│   ├── app/
│   │
│   ├── core/
│   │   ├── config/
│   │   ├── context/
│   │   ├── domain/
│   │   ├── errors/
│   │   ├── infrastructure/
│   │   │   ├── postgres/
│   │   │   ├── redis/
│   │   │   ├── mongo/
│   │   │   ├── s3/
│   │   │   └── ...
│   │   ├── logger/
│   │   ├── security/
│   │   ├── transport/
│   │   └── ...
│   │
│   ├── features/
│   │   ├── auth/
│   │   ├── client/
│   │   ├── media/
│   │   ├── requests/
│   │   ├── sessions/
│   │   ├── settings/
│   │   ├── users/
│   │   └── ...
│   │
│   └── workers/
│       └── outbox/
│
├── migrations/
├── scripts/
├── docs/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
└── go.mod
```

## How to start

### Makefile targets

| Command                               | Description                          |
| ------------------------------------- | ------------------------------------ |
| `make env-up`                         | Start PostgreSQL, Redis and MongoDB  |
| `make env-down`                       | Stop infrastructure                  |
| `make env-cleanup`                    | Remove local database data           |
| `make env-port-forward`               | Expose PostgreSQL locally            |
| `make migrate-up`                     | Apply migrations                     |
| `make migrate-down`                   | Roll back migration                  |
| `make migrate-create name=<name>`     | Create migration                     |
| `make migrate-goto version=<version>` | Go to migration version              |
| `make mongo-init`                     | Initialize MongoDB                   |
| `make swagger-gen`                    | Generate Swagger documentation       |
| `make app-run`                        | Run application locally              |
| `make app-deploy`                     | Build and deploy with Docker Compose |

### Order to quickstart

1. Run `make env-up` to create & start images
2. Expose ports using `make env-port-forward`
3. For migrations run `make migrate-up` and `make mongo-init`
4. To create swagger run `make swagger-gen`

And **finally** run `make app-run` to start application

_P.S don't forget to configure .env file_
_P.S.S if you're using Windows, check Makefile for \*-windows targets_

## Graceful Shutdown

The application listens for:

```bash
SIGINT
SIGTERM
```

and propagates cancellation through the application context.

This allows the HTTP server and background workers to shut down gracefully instead of being terminated abruptly.

## Status

> Work in progress

Cascade Pro API is actively developed and its architecture and API may change.
