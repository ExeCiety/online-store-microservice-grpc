# Online Store Microservice (Go + gRPC + PostgreSQL)

A simple online store backend microservice project built in a monorepo architecture.

## Architecture

```text
Client (HTTP/JSON)
      |
      v
+--------------------+
|    API Gateway     | :8080 (Gin REST)
+--------------------+
      | gRPC                           | gRPC
      v                                v
+--------------------+         +--------------------+
|    User Service    | :50051  |   Order Service    | :50052
+--------------------+         +--------------------+
      |                                |
      v                                v
PostgreSQL user_db (:5432)      PostgreSQL order_db (:5433)

Swagger UI: :8081
```

## Tech Stack

- Go 1.26+
- Gin (HTTP API)
- gRPC
- PostgreSQL 15
- GORM
- bcrypt

## Project Structure

```text
.
├── api-gateway/
├── user-service/
├── order-service/
├── proto/
├── docs/
│   └── swagger.yaml
├── scripts/
├── docker-compose.yml
├── Makefile
└── README.md
```

## Prerequisites

- Go >= 1.26
- Docker + Docker Compose
- Optional: `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`

## Quick Setup

1. Copy environment file:

```bash
cp .env.example .env
```

2. Install dependencies:

```bash
make tidy
```

3. Start infrastructure (PostgreSQL + Swagger UI):

```bash
docker compose up -d
```

4. Run migrations (idempotent, no local `psql` required):

```bash
make migrate-user
make migrate-order
```

Migration targets will automatically create databases if they do not exist.

5. (Optional) Generate proto files:

```bash
make proto
```

## Running Services

### Option A: Go run (3 terminals)

```bash
make run-user
make run-order
make run-gateway
```

### Option B: Build binaries and run all

```bash
make build
make run-built
```

Built binaries are stored in `bin/`:
- `bin/user-service`
- `bin/order-service`
- `bin/api-gateway`

## Testing

- Run all tests:

```bash
make test
```

- Run API Gateway endpoint tests only:

```bash
make test-gateway
```

- Run all package tests directly:

```bash
go test ./...
```

Endpoint tests are separated per file under `api-gateway/tests/`.

## Swagger

- OpenAPI spec: `docs/swagger.yaml`
- Swagger UI (via Compose): `http://localhost:8081`

## API Endpoints

- `POST /api/register`
- `POST /api/login`
- `GET /api/users/:id`
- `POST /api/orders`
- `GET /api/orders/:id`
- `GET /api/users/:userId/orders`
- `GET /health`

## Example Requests

### Register User

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","name":"John Doe"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

### Create Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":"4e427d78-58c5-4f78-bfc1-e2c196e0b506","product_name":"Laptop","quantity":1,"total_price":15000000}'
```

### Get Orders by User

```bash
curl -X GET http://localhost:8080/api/users/4e427d78-58c5-4f78-bfc1-e2c196e0b506/orders
```

## Notes

- Passwords are stored using bcrypt hashing.
- Inter-service communication uses gRPC.
- API Gateway is stateless.
- Logging, validation, error handling, and graceful shutdown are implemented.
- JWT is still a placeholder (dummy token in login response).
