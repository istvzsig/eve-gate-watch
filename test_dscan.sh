#!/bin/bash

SYSTEM_NAME=${1:-"tama"}
curl -sS -X GET "http://localhost:8080/api/dscan?system=$SYSTEM_NAME"