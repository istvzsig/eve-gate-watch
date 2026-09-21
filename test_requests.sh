#!/bin/bash

SYSTEM_NAME="Tama"

curl -X POST http://localhost:8080/api/dscan \
  -H 'Content-Type: application/json' \
  -d "{
    \"systemName\": \"$SYSTEM_NAME\",
    \"ships\": [
      {\"name\": \"Rifter\"}
    ]
  }"

curl -X GET "http://localhost:8080/api/intel?system=$SYSTEM_NAME"