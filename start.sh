#!/bin/bash

set -euo pipefail

cp .env .env.bkp

source .env

go build -o gatewatch ./cmd/gatewatch

./gatewatch