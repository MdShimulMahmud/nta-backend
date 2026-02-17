#!/bin/bash

# health_check.sh - Health check script for NTA Backend

API_HOST="${API_HOST:-localhost}"
API_PORT="${API_PORT:-8080}"
URL="http://${API_HOST}:${API_PORT}/health"

echo "Checking NTA Backend health at $URL..."

response=$(curl -s -o /dev/null -w "%{http_code}" "$URL" 2>/dev/null)

if [ "$response" = "200" ]; then
  echo "✓ NTA Backend is healthy"
  exit 0
else
  echo "✗ NTA Backend is not responding (HTTP $response)"
  exit 1
fi
