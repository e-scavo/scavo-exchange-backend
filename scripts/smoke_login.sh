#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${SCAVO_BASE_URL:-http://localhost:8080}"

echo "Smoke testing login against: $BASE_URL"

curl -sS -X POST "$BASE_URL/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@scavo.exchange","password":"dev"}'
echo