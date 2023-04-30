#!/bin/bash

# This script will call the Infura API to get the latest block number
INFURA_URL="https://goerli.infura.io/v3/e1f81b43fa6e46a9a7ec9c48165732b1"

echo "Calling Infura API to get latest block number..."
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  $INFURA_URL

echo ""


# Get Account Balance
ACCOUNT="0xA8b82FFe97BD9A7Ef74AB258814378E8Be590A35"

echo "Calling Infura API to get account balance..."
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["'$ACCOUNT'", "latest"],"id":1}' \
  $INFURA_URL

echo ""
echo "Converting balance to decimal eth..."

# Convert hex to decimal
BALANCE_HEX=$(curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["'$ACCOUNT'", "latest"],"id":1}' \
  $INFURA_URL | jq -r '.result')

BALANCE_DEC=$(printf "%d" $BALANCE_HEX)

echo "Balance: $BALANCE_DEC wei"
echo "Balance: $(bc <<< "scale=5; $BALANCE_DEC / 1000000000000000000") eth"
echo ""
