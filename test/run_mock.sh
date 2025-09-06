#!/bin/bash

echo "Starting mock server to capture Minecraft mod requests..."
echo "Configure your mod to point to localhost:8080"
echo "Press Ctrl+C to stop"
echo ""

cd "$(dirname "$0")"
go run mock_server.go