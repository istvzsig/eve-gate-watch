#!/bin/bash

test_endpoint() {
	local path=$1
	curl -vvv -X GET "http://localhost:8080/api/$path?system=$SYSTEM_NAME"
}

SYSTEM_NAME=${1:-"tama"}

test_endpoint "dscan"
test_endpoint "pvp"