#!/bin/bash
#
# This script launches an rpc node ( Normal Geth Node for Clique PoA network )

STATE_RESET=0

# Defaults give theoretical max throughput of 1800 TPS = ~15 tx/sec ( mainnet ) X 12 X 10
CHAIN_ID=505
PERIOD=1 # 1 second per block ( 12x faster than Ethereum mainnet )
GAS_LIMIT=300000000 # 300M gas limit ( 10x Ethereum mainnet )

PEER_PORT=30303
HTTP_PORT=8545
RPC_PORT=8551

display_help() {
  echo "Usage: launch-rpc.sh -d <data-dir> -g <genesis-file> [options]"
  echo
  echo "   -h, --help                 show help"
  echo "   -d, --data                 Geth node data dir (required)"
  echo "   -g, --genesis              Genesis file (required)"
  echo "   -x, --clear                Clear & Reset data dir before starting"

  echo "   -c, --chain-id             Chain ID"
  echo "   -p, --period               Block period ( in seconds )"
  echo "   -g, --gas-limit            Gas limit"

  echo "   -H, --http-port            Http port ( default: 8545 )"
  echo "   -m, --peer-port            Peer port ( default: 30303 )"
  echo "   -r, --rpc-port             RPC port ( default: 8551 )"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  echo "Example: ./scripts/launch-rpc.sh -d ~/l1-rpc-data/ -g build/genesis.json -x -m 30306 -r 8550 -H 8548"
}

clear_data() {
  rm -rf $DATA_DIR
  mkdir -p $DATA_DIR
}

# Parse command line arguments
while getopts ":hd:g:xc:p:g:m:H:r:o:" opt; do
  case ${opt} in
    d|--data )
      DATA_DIR=$OPTARG
      ;;
    g|--genesis )
      GENESIS_FILE=$OPTARG
      ;;
    h|--help )
      display_help
      exit 0
      ;;
    x|--clear )
      clear_data
      STATE_RESET=1
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
if [[ -z "$DATA_DIR" || -z "$GENESIS_FILE" ]]; then
  echo "Missing required arguments" 1>&2
  display_help
  exit 1
fi

# Check if data dir exists
if [ ! -d "$DATA_DIR" ]; then
  echo "Data dir does not exist: $DATA_DIR, setting mode to STATE_RESET=1"
  clear_data
  STATE_RESET=1
fi

# Also check data dir has data
if [ ! -d "$DATA_DIR/geth" ]; then
  echo "Data dir does not contain data: $DATA_DIR, setting mode to STATE_RESET=1"
  clear_data
  STATE_RESET=1
fi


# Check if genesis file exists
if [ ! -f "$GENESIS_FILE" ]; then
  echo "Genesis file does not exist: $GENESIS_FILE" 1>&2
  display_help
  exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR=$SCRIPT_DIR/..
PASSWORD_FILE=$DATA_DIR/password.txt

if [ $STATE_RESET -eq 1 ]; then
  # Create Geth Account for RPC node
  ACCOUNT_PASS=${ACCOUNT_PASS:-password}
  touch $PASSWORD_FILE
  echo $ACCOUNT_PASS > $PASSWORD_FILE
  geth account new --datadir $DATA_DIR --password $PASSWORD_FILE
fi
ACCOUNT1=$(cat $DATA_DIR/keystore/* | jq -r '.address' | head -n 1)

if [ $STATE_RESET -eq 1 ]; then
  # Create Geth Genesis & Init Chain
  geth init --datadir $DATA_DIR $GENESIS_FILE
fi

if [ -z "$OUTPUT_FILE" ]; then
  ${WORK_DIR}/go-ethereum/build/bin/geth --networkid $CHAIN_ID --datadir $DATA_DIR --http --http.api "eth,net,web3,personal,txpool" --http.port $HTTP_PORT --unlock "0x$ACCOUNT1" --allow-insecure-unlock --password $PASSWORD_FILE --port $PEER_PORT --authrpc.port $RPC_PORT --metrics --metrics.addr 127.0.0.1 --metrics.expensive --metrics.port 6061
else
  ${WORK_DIR}/go-ethereum/build/bin/geth --networkid $CHAIN_ID --datadir $DATA_DIR --http --http.api "eth,net,web3,personal,txpool" --http.port $HTTP_PORT --unlock "0x$ACCOUNT1" --allow-insecure-unlock --password $PASSWORD_FILE --port $PEER_PORT --authrpc.port $RPC_PORT --metrics --metrics.addr 127.0.0.1 --metrics.expensive --metrics.port 6061 > $OUTPUT_FILE 2>&1 &
fi
