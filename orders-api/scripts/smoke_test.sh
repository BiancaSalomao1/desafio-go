#!/bin/bash

set -e

echo "==============================="
echo "Swagger"
echo "==============================="

curl -f http://localhost:8080/swagger/index.html >/dev/null

echo
echo "==============================="
echo "Login"
echo "==============================="

TOKEN=$(
curl -s \
-X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{
  "email":"admin@email.com",
  "password":"123456"
}' | jq -r '.accessToken'
)

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Falha no login"
    exit 1
fi

echo "Token obtido."

echo
echo "==============================="
echo "Produtos"
echo "==============================="

curl \
-H "Authorization: Bearer $TOKEN" \
http://localhost:8080/products

echo
echo
echo "Smoke Test OK"