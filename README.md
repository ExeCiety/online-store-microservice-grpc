# Online Store Microservice (Go + gRPC + PostgreSQL)

Backend microservice untuk toko online sederhana dengan arsitektur monorepo.

## Arsitektur

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

## Struktur Project

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

## Prasyarat

- Go >= 1.26
- Docker + Docker Compose
- Optional: `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`

## Setup Cepat

1. Copy env:

```bash
cp .env.example .env
```

2. Install dependency:

```bash
make tidy
```

3. Jalankan infra (PostgreSQL + Swagger UI):

```bash
docker compose up -d
```

4. Jalankan migration (idempotent, tanpa `psql` lokal):

```bash
make migrate-user
make migrate-order
```

Target migration akan otomatis membuat database jika belum ada.

5. (Optional) Generate proto:

```bash
make proto
```

## Menjalankan Service

### Opsi A: Go run (3 terminal)

```bash
make run-user
make run-order
make run-gateway
```

### Opsi B: Build binary lalu jalankan semua

```bash
make build
make run-built
```

Output binary ada di folder `bin/`:
- `bin/user-service`
- `bin/order-service`
- `bin/api-gateway`

## Testing

- Jalankan semua test:

```bash
make test
```

- Jalankan test endpoint API Gateway:

```bash
make test-gateway
```

- Jalankan semua test package:

```bash
go test ./...
```

Test endpoint dipisah per file di folder `api-gateway/tests/`.

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

## Contoh Request

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

- Password disimpan sebagai hash bcrypt.
- Komunikasi antar service menggunakan gRPC.
- API Gateway stateless.
- Sudah ada logging, validation, error handling, graceful shutdown.
- JWT masih placeholder (dummy token pada login).
