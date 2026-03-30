#!/bin/bash

set -e

# =========================
# 📌 CONFIG
# =========================

PROJECT_NAME="scavo-exchange-backend"

# Validar parámetro (etapa)
if [ -z "$1" ]; then
  echo "❌ Uso: ./backup.sh <etapa>"
  echo "Ejemplo: ./backup.sh 7.3.4"
  exit 1
fi

STAGE="$1"
ZIP_NAME="${PROJECT_NAME}-${STAGE}.zip"
OUTPUT_PATH="../${ZIP_NAME}"

# =========================
# 🧹 CLEAN PREVIOUS
# =========================

echo "🧹 Eliminando zip anterior..."
rm -f "$OUTPUT_PATH"

# =========================
# 📦 CREATE ZIP
# =========================

echo "📦 Generando backup: $ZIP_NAME"

zip -r "$OUTPUT_PATH" . \
  -x ".git/*" ".git" \
     "node_modules/*" \
     "*.log" \
     ".env" \
     ".DS_Store" \
     "build/*" \
     "cmd/scavo_geryon_be/scavo_geryon*" \
     "cmd/scavo_geryon_be_v2/scavo_geryon*" \
     "cmd/scavium_network_be/*" \
     "cmd/scavo_mercadopago/*" \
     "cmd/scavo_site_be/*" \
     "cmd/scavo_site_fe/*" \
     "cmd/scavo_tgbot/*" \
     "cmd/scavo_wallets/*" \
     "scavo-exchange-backend*.zip" \
     "*.exe" \
     "*.out"

# =========================
# 📂 COPY BACK TO PROJECT
# =========================

echo "📂 Copiando zip al proyecto..."
cp "$OUTPUT_PATH" .

# =========================
# ✅ DONE
# =========================

echo "✅ Backup generado: $ZIP_NAME"
