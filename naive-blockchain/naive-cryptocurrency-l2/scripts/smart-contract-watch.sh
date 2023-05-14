#!/bin/bash
#
# This script is used to watch the smart contract for changes and generates stdout logs

SMART_CONTRACT_WATCH_DIR=/home/b-j-roberts/workspace/blockchain/tools/smart-contract-watch
CONTRACT_HOST="http://localhost"
CONTRACT_PORT="8545"

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -a, --address           Smart contract address (required)"
  echo "   -d, --dir               Smart contract watch directory -- Neufund/smart-contract-watch.git repo (default: $SMART_CONTRACT_WATCH_DIR)"
  echo "   -H, --host              Blockchain host (default: $CONTRACT_HOST)"
  echo "   -p, --port              Blockchain port (default: $CONTRACT_PORT)"
  echo "   -h, --help              Show help"
  echo
  exit 1
}

while getopts ":a:d:H:p:h" opt; do
  case $opt in
    a)
      SMART_CONTRACT_ADDRESS=$OPTARG
      ;;
    d)
      SMART_CONTRACT_WATCH_DIR=$OPTARG
      ;;
    H)
      CONTRACT_HOST=$OPTARG
      ;;
    p)
      CONTRACT_PORT=$OPTARG
      ;;
    h)
      display_help
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      display_help
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      display_help
      ;;
  esac
done

if [ -z "$SMART_CONTRACT_ADDRESS" ]; then
  echo "Smart contract address is required"
  display_help
fi

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR="${BASE_DIR}/.."

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
yarn start -f 1 -a ${SMART_CONTRACT_ADDRESS} -n $CONTRACT_HOST:$CONTRACT_PORT -l "info" -q
