#!/usr/bin/env bash

# ==========================================
# Testes da API
#
# Responsabilidades:
# - testar os endpoints da API;
# - facilitar testes durante o desenvolvimento.
# ==========================================

BASE_URL="http://localhost:8080"

echo
echo "====================================="
echo "TESTE DA API"
echo "====================================="

###########################################
# PRODUTOS
###########################################

echo
echo ">>> Criando Produto"

TIMESTAMP=$(date +%s)

PRODUCT_NAME="Notebook $TIMESTAMP"

CUSTOMER_EMAIL="bianca_${TIMESTAMP}@email.com"

USER_EMAIL="admin_${TIMESTAMP}@email.com"


PRODUCT_RESPONSE=$(curl -s \
-X POST "$BASE_URL/products" \
-H "Content-Type: application/json" \
-d '{
    "name":"'"$PRODUCT_NAME"'" , 
    "price":3500,
    "stock":10
}')

echo "$PRODUCT_RESPONSE"

PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | grep -oP '"id":"\K[^"]+')

echo
echo "Produto ID:"
echo "$PRODUCT_ID"

echo
echo ">>> Listando Produtos"

curl -s \
"$BASE_URL/products"

echo
echo
echo ">>> Buscando Produto"

curl -s \
"$BASE_URL/products/$PRODUCT_ID"

###########################################
# CLIENTES
###########################################

echo
echo
echo ">>> Criando Cliente"

CUSTOMER_RESPONSE=$(curl -s \
-X POST "$BASE_URL/customers" \
-H "Content-Type: application/json" \
-d '{
    "name":"Bianca",
    "email":"'"$CUSTOMER_EMAIL"'"
}')


if echo "$CUSTOMER_RESPONSE" | grep -q '"error"'; then
    echo "Erro ao criar cliente."
    exit 1
fi

echo "$CUSTOMER_RESPONSE"

CUSTOMER_ID=$(echo "$CUSTOMER_RESPONSE" | grep -oP '"id":"\K[^"]+')

echo
echo "Cliente ID:"
echo "$CUSTOMER_ID"

echo
echo ">>> Listando Clientes"

curl -s \
"$BASE_URL/customers"

echo
echo
echo ">>> Buscando Cliente"

curl -s \
"$BASE_URL/customers/$CUSTOMER_ID"

###########################################
# USUÁRIOS
###########################################

echo
echo
echo ">>> Criando Usuário"

USER_RESPONSE=$(curl -s \
-X POST "$BASE_URL/users" \
-H "Content-Type: application/json" \
-d '{
    "name":"Administrador",
    "email":"'"$USER_EMAIL"'",
    "password":"123456"
}')

echo "$USER_RESPONSE"

USER_ID=$(echo "$USER_RESPONSE" | grep -oP '"id":"\K[^"]+')

echo
echo "Usuário ID:"
echo "$USER_ID"

echo
echo ">>> Listando Usuários"

curl -s \
"$BASE_URL/users"

echo
echo
echo ">>> Buscando Usuário"

curl -s \
"$BASE_URL/users/$USER_ID"

echo
echo
echo "====================================="
echo "TESTES FINALIZADOS"
echo "====================================="

