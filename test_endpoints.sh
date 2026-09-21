#!/bin/bash

test_endpoint() {
	local path=$1
	curl -sS -X GET "http://localhost:8080/api/$path"
}

SYSTEM_NAME=${1:-"tama"}

test_endpoint "systems?q=tam"
test_endpoint "systems?q=ji"
test_endpoint "systems?q="
test_endpoint "systems?q=xyz"
test_endpoint "systems?q=JI"
test_endpoint "systems?q=t"
test_endpoint "dscan?system=$SYSTEM_NAME"
test_endpoint "pvp?system=$SYSTEM_NAME"
