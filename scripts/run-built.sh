#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ ! -x "$ROOT_DIR/bin/user-service" || ! -x "$ROOT_DIR/bin/order-service" || ! -x "$ROOT_DIR/bin/api-gateway" ]]; then
  echo "binary belum ada. jalankan: make build"
  exit 1
fi

if [[ -f "$ROOT_DIR/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ROOT_DIR/.env"
  set +a
fi

export USER_SERVICE_GRPC_PORT="${USER_SERVICE_GRPC_PORT:-50051}"
export USER_DB_HOST="${USER_DB_HOST:-localhost}"
export USER_DB_PORT="${USER_DB_PORT:-5432}"
export USER_DB_USER="${USER_DB_USER:-postgres}"
export USER_DB_PASSWORD="${USER_DB_PASSWORD:-postgres}"
export USER_DB_NAME="${USER_DB_NAME:-user_db}"

export ORDER_SERVICE_GRPC_PORT="${ORDER_SERVICE_GRPC_PORT:-50052}"
export ORDER_DB_HOST="${ORDER_DB_HOST:-localhost}"
export ORDER_DB_PORT="${ORDER_DB_PORT:-5433}"
export ORDER_DB_USER="${ORDER_DB_USER:-postgres}"
export ORDER_DB_PASSWORD="${ORDER_DB_PASSWORD:-postgres}"
export ORDER_DB_NAME="${ORDER_DB_NAME:-order_db}"

export API_GATEWAY_PORT="${API_GATEWAY_PORT:-8080}"
export USER_SERVICE_URL="${USER_SERVICE_URL:-localhost:50051}"
export ORDER_SERVICE_URL="${ORDER_SERVICE_URL:-localhost:50052}"

"$ROOT_DIR/bin/user-service" &
PID_USER=$!
"$ROOT_DIR/bin/order-service" &
PID_ORDER=$!
"$ROOT_DIR/bin/api-gateway" &
PID_GATEWAY=$!

cleanup() {
  kill "$PID_GATEWAY" "$PID_ORDER" "$PID_USER" 2>/dev/null || true
  wait "$PID_GATEWAY" "$PID_ORDER" "$PID_USER" 2>/dev/null || true
}

trap cleanup INT TERM

echo "running binaries:"
echo "- user-service  : pid=$PID_USER"
echo "- order-service : pid=$PID_ORDER"
echo "- api-gateway   : pid=$PID_GATEWAY"

echo "press Ctrl+C to stop all"
wait "$PID_GATEWAY" "$PID_ORDER" "$PID_USER"
