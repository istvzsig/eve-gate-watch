#!/bin/bash

curl -vvv -X POST http://localhost:8080/api/dscan   -H 'Content-Type: application/json'   -d '{
    "systemName": "Tama",
    "ships": ["Sabre", "Ishtar", "Drake"]
  }'