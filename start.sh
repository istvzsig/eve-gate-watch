#!/bin/bash

set -euo pipefail

cp .env .env.bkp

source .env

export DISCORD_BOT_TOKEN="$DISCORD_BOT_TOKEN"
export DATABASE_URL="$DATABASE_URL"

BIN_DIR=".bin"
mkdir -p "$BIN_DIR"

go build -o "$BIN_DIR/migrate" ./cmd/migrate

echo "==> Running database migrations..."
"$BIN_DIR/migrate"

go build -o "$BIN_DIR/sdeimport" ./cmd/sdeimport

echo "==> Seeding EVE SDE..."
"$BIN_DIR/sdeimport"

go build -o "$BIN_DIR/gatewatch" ./cmd/gatewatch

echo "==> Starting Gatewatch..."
"$BIN_DIR/gatewatch"