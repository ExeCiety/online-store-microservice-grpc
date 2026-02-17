#!/usr/bin/env bash
set -euo pipefail

if ! command -v protoc >/dev/null 2>&1; then
  echo "error: protoc not found. install Protocol Buffers compiler first." >&2
  exit 1
fi

if ! command -v protoc-gen-go >/dev/null 2>&1 || ! command -v protoc-gen-go-grpc >/dev/null 2>&1; then
  echo "error: protoc-gen-go and protoc-gen-go-grpc are required." >&2
  echo "install with:" >&2
  echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" >&2
  echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" >&2
  exit 1
fi

protoc --go_out=. --go-grpc_out=. proto/user/user.proto
protoc --go_out=. --go-grpc_out=. proto/order/order.proto

echo "proto generation completed"
