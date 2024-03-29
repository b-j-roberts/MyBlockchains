#!/bin/bash
#
# This script launches an rpc node ( Normal Geth Node for Clique PoA network )

# Defaults give theoretical max throughput of 1800 TPS = ~15 tx/sec ( mainnet ) X 12 X 10
CHAIN_ID=505
PERIOD=1 # 1 second per block ( 12x faster than Ethereum mainnet )
GAS_LIMIT=300000000 # 300M gas limit ( 10x Ethereum mainnet )

PEER_PORT=30303
HTTP_PORT=8545
RPC_PORT=8551

display_help() {
  echo "Usage: $0 [option...] {arguments...}"
  echo "WARNING: Requires the datadir be set up prior w/ genesis file. Use ./scripts/setup-rpc.sh for this"
  echo
  echo "   -h, --help                 show help"
  echo "   -d, --data                 Geth node data dir (required)"

  echo "   -c, --chain-id             Chain ID"
  echo "   -p, --period               Block period ( in seconds )"
  echo "   -g, --gas-limit            Gas limit"

  echo "   -H, --http-port            Http port ( default: 8545 )"
  echo "   -m, --peer-port            Peer port ( default: 30303 )"
  echo "   -r, --rpc-port             RPC port ( default: 8551 )"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  echo "Example: $0 -m 30306 -r 8550 -H 8548"
}

# Parse command line arguments
while getopts ":hd:c:p:g:m:H:r:o:" opt; do
  case ${opt} in
    d|--data )
      DATA_DIR=$OPTARG
      ;;
    h|--help )
      display_help
      exit 0
      ;;
    c|--chain-id )
      CHAIN_ID=$OPTARG
      ;;
    p|--period )
      PERIOD=$OPTARG
      ;;
    g|--gas-limit )
      GAS_LIMIT=$OPTARG
      ;;
    m|--peer-port )
      PEER_PORT=$OPTARG
      ;;
    H|--http-port )
      HTTP_PORT=$OPTARG
      ;;
    r|--rpc-port )
      RPC_PORT=$OPTARG
      ;;
    o|--output )
      OUTPUT_FILE=$OPTARG
      ;; 
    \? )
      echo "Invalid Option: -$OPTARG" 1>&2
      display_help
      exit 1
      ;;
    : )
      echo "Invalid Option: -$OPTARG requires an argument" 1>&2
      display_help
      exit 1
      ;;
  esac
done

# Check if required arguments are set
if [[ -z "$DATA_DIR" ]]; then
  echo "Missing required arguments" 1>&2
  display_help
  exit 1
fi

# Check if data dir exists
if [ ! -d "$DATA_DIR" ]; then
  echo "Data dir does not exist: $DATA_DIR, it must be setup first"
  display_help
  exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR=$SCRIPT_DIR/..
PASSWORD_FILE=$DATA_DIR/password.txt

# Create Geth Account for RPC node
if [ -z "$(ls -A ${HOME}/.eth-accounts)" ]; then
  echo "No accounts found, creating new account"
  $WORK_DIR/scripts/generate-account.sh -d ${DATA_DIR} -x
fi

cp -r ${HOME}/.eth-accounts/* $DATA_DIR/keystore/
mv $DATA_DIR/keystore/password.txt $PASSWORD_FILE
ACCOUNT1=$(cat $DATA_DIR/keystore/* | jq -r '.address' | head -n 1)

if [ -z "$OUTPUT_FILE" ]; then
  ${WORK_DIR}/go-ethereum/build/bin/geth --networkid $CHAIN_ID --datadir $DATA_DIR --http --http.api "eth,net,web3,personal,txpool,admin" --http.port $HTTP_PORT --unlock "0x$ACCOUNT1" --allow-insecure-unlock --password $PASSWORD_FILE --port $PEER_PORT --authrpc.port $RPC_PORT --metrics --metrics.addr 127.0.0.1 --metrics.expensive --metrics.port 6061
else
  ${WORK_DIR}/go-ethereum/build/bin/geth --networkid $CHAIN_ID --datadir $DATA_DIR --http --http.api "eth,net,web3,personal,txpool,admin" --http.port $HTTP_PORT --unlock "0x$ACCOUNT1" --allow-insecure-unlock --password $PASSWORD_FILE --port $PEER_PORT --authrpc.port $RPC_PORT --metrics --metrics.addr 127.0.0.1 --metrics.expensive --metrics.port 6061 > $OUTPUT_FILE 2>&1 &
fi
