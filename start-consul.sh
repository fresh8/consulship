#!/bin/bash
CONTAINER="$(docker run -P -h consul -d consul agent -server -bootstrap -ui -client=0.0.0.0)"
PORT_OUTPUT="$(docker port ${CONTAINER} | grep "8500/tcp" | cut -d':' -f2)"
echo "Consul Port: ${PORT_OUTPUT}"
