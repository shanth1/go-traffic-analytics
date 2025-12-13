#!/bin/bash
set -e

API_KEY=x-api-key

LOCAL_API_URL="http://localhost:8080"
SWAGGER_PATH_SUFFIX="/swagger/doc.json"

if [[ "$OSTYPE" == "darwin"* ]] || [[ "$OS" == "Windows_NT" ]]; then
    DOCKER_BASE_URL="http://host.docker.internal:8080"
else
    DOCKER_BASE_URL="http://localhost:8080"
fi

TEST_EMAIL="swagger_test_$(date +%s)@example.com"
TEST_PASS="someSecurePass123"

echo "--- 1. Ensuring Server is Up ---"
curl -s -o /dev/null --retry 5 --retry-connrefused "$LOCAL_API_URL/health" || \
    (echo "Server is not running at $LOCAL_API_URL" && exit 1)

echo "--- 2. Getting Auth Token ---"
curl -s -X POST "$LOCAL_API_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -H "X-API-KEY: ${API_KEY}" \
    -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"$TEST_PASS\"}" > /dev/null

TOKEN=$(curl -s -X POST "$LOCAL_API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"$TEST_PASS\"}" | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get token!"
    exit 1
fi
echo "Token received."

echo "--- 3. Running Contract Tests ---"
echo "Targeting Base URL: $DOCKER_BASE_URL"
echo "Swagger File: $DOCKER_BASE_URL$SWAGGER_PATH_SUFFIX"

docker run --rm \
    --network="host" \
    -e SCHEMATHESIS_BASE_URL="$DOCKER_BASE_URL" \
    schemathesis/schemathesis:stable \
    run "$DOCKER_BASE_URL$SWAGGER_PATH_SUFFIX" \
    --checks all \
    --header "Authorization: Bearer $TOKEN"

echo "--- SUCCESS ---"
