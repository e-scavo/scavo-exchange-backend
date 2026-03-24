#!/usr/bin/env bash
set -e

if [ -z "$SCAVO_POSTGRES_URL" ]; then
  echo "SCAVO_POSTGRES_URL not set"
  exit 1
fi

CMD=$1

if [ -z "$CMD" ]; then
  CMD="up"
fi

echo "Running migrations: $CMD"

goose -dir migrations postgres "$SCAVO_POSTGRES_URL" $CMD