.PHONY: tidy proto up down migrate-user migrate-order run-user run-order run-gateway build run-built test test-gateway

USER_DB_NAME ?= online_microservice_user_db
ORDER_DB_NAME ?= online_microservice_order_db

tidy:
	go mod tidy

proto:
	./scripts/generate-proto.sh

up:
	docker-compose up -d

down:
	docker-compose down

migrate-user:
	docker exec user_db sh -c "psql -U postgres -d postgres -tAc \"SELECT 1 FROM pg_database WHERE datname='$(USER_DB_NAME)'\" | grep -q 1 || psql -U postgres -d postgres -c \"CREATE DATABASE \\\"$(USER_DB_NAME)\\\"\""
	docker exec -i user_db psql -U postgres -d $(USER_DB_NAME) < scripts/migration/user_db.sql

migrate-order:
	docker exec order_db sh -c "psql -U postgres -d postgres -tAc \"SELECT 1 FROM pg_database WHERE datname='$(ORDER_DB_NAME)'\" | grep -q 1 || psql -U postgres -d postgres -c \"CREATE DATABASE \\\"$(ORDER_DB_NAME)\\\"\""
	docker exec -i order_db psql -U postgres -d $(ORDER_DB_NAME) < scripts/migration/order_db.sql

run-user:
	cd user-service && go run main.go

run-order:
	cd order-service && go run main.go

run-gateway:
	cd api-gateway && go run main.go

build:
	mkdir -p bin
	go build -o bin/user-service ./user-service
	go build -o bin/order-service ./order-service
	go build -o bin/api-gateway ./api-gateway

run-built:
	./scripts/run-built.sh

test:
	go test ./...

test-gateway:
	go test ./api-gateway/tests -v
