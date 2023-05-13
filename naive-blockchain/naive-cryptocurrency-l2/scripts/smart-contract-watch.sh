#!/bin/bash

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR="${BASE_DIR}/.."

SMART_CONTRACT_WATCH_DIR=/home/b-j-roberts/workspace/blockchain/tools/smart-contract-watch

SMART_CONTRACT_ADDRESS=$(cat ${BASE_DIR}/contracts/builds/contract-address.txt | jq -r '.address')

# exit if no smart contract watch directory
if [ ! -d "$SMART_CONTRACT_WATCH_DIR" ]; then
  echo "Smart contract watch directory not found: $SMART_CONTRACT_WATCH_DIR"
  exit 1
fi

cd $SMART_CONTRACT_WATCH_DIR
rm -rf contracts/
mkdir contracts

echo "Smart contract address: $SMART_CONTRACT_ADDRESS"
# Copy contracts from the blockchain
cp ${BASE_DIR}/contracts/builds/TransactionStorage.abi contracts/${SMART_CONTRACT_ADDRESS}.json
yarn start -f 1 -a ${SMART_CONTRACT_ADDRESS} -n "http://localhost:8545" -l "info" -q
