#!/bin/bash

SYSTEM_NAME="jita"

curl -sS -X GET "http://localhost:8080/api/dscan?system=$SYSTEM_NAME"