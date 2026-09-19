#!/bin/bash

cp .env .env.bkp
source .env

go build -o gatewatch ./cmd/gatewatch
./gatewatch