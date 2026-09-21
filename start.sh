#!/bin/bash

set -euo pipefail

cp .env .env.bkp

source .env

export DISCORD_BOT_TOKEN=$DISCORD_BOT_TOKEN

go build -o gatewatch ./cmd/gatewatch

./gatewatch