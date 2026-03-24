#!/usr/bin/env bash
set -euo pipefail

if [ -z "${SCAVO_POSTGRES_URL:-}" ]; then
  echo "SCAVO_POSTGRES_URL not set"
  exit 1
fi

CMD="${1:-up}"

case "$CMD" in
  up|up-by-one|down|down-to|status|version|redo|reset)
    ;;
  *)
    echo "Unsupported migration command: $CMD"
    echo "Supported: up, up-by-one, down, down-to, status, version, redo, reset"
    exit 1
    ;;
esac

echo "Running migrations: $CMD"

goose -dir migrations postgres "$SCAVO_POSTGRES_URL" "$CMD" "${@:2}"