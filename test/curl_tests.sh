#!/bin/bash

# Lightweight API testing script
API_URL="http://localhost:8080"
WEB_URL="http://localhost:8081"

echo "=== PokéFactory API Tests ==="

# 1. Health Check
echo "1. Testing health endpoint..."
curl -s $API_URL/health | jq '.'

# 2. Server Authentication
echo -e "\n2. Getting server token..."
SERVER_TOKEN=$(curl -s -X POST $API_URL/api/v1/server/auth \
  -H "Content-Type: application/json" \
  -d '{"server_id":"test-server-1","server_key":"your-jwt-secret"}' | jq -r '.token')

echo "Server token: ${SERVER_TOKEN:0:20}..."

# 3. Create Player
echo -e "\n3. Creating test player..."
curl -s -X POST $API_URL/api/v1/server/player/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SERVER_TOKEN" \
  -d '{"player_uuid":"550e8400-e29b-41d4-a716-446655440000","username":"TestPlayer"}' | jq '.'

# 4. Catch Pikachu
echo -e "\n4. Catching Pikachu (#25)..."
curl -s -X POST $API_URL/api/v1/server/pokedex/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SERVER_TOKEN" \
  -d '{"player_uuid":"550e8400-e29b-41d4-a716-446655440000","national_id":25,"action":"catch"}' | jq '.'

# 5. Catch Charizard
echo -e "\n5. Catching Charizard (#6)..."
curl -s -X POST $API_URL/api/v1/server/pokedex/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SERVER_TOKEN" \
  -d '{"player_uuid":"550e8400-e29b-41d4-a716-446655440000","national_id":6,"action":"catch"}' | jq '.'

# 6. Get Player Stats
echo -e "\n6. Getting player stats..."
curl -s -X POST $API_URL/api/v1/server/pokedex/summary \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $SERVER_TOKEN" \
  -d '{"player_uuid":"550e8400-e29b-41d4-a716-446655440000"}' | jq '.'

# 7. Web Dashboard Tests
echo -e "\n7. Testing web leaderboards..."
curl -s $WEB_URL/api/v1/web/leaderboards | jq '.'

echo -e "\n8. Testing web player stats..."
curl -s $WEB_URL/api/v1/web/player/TestPlayer/stats | jq '.'

echo -e "\n9. Testing server analytics..."
curl -s $WEB_URL/api/v1/web/server/analytics | jq '.'

echo -e "\n=== Tests Complete ==="