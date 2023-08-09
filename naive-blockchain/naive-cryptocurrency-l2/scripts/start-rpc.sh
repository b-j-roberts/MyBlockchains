#!/bin/bash
#
# This script starts the sequencer in a docker container.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

STATE_RESET=0

# Defaults give theoretical max throughput of 1800 TPS = ~15 tx/sec ( mainnet ) X 12 X 10
CHAIN_ID=515
PERIOD=1 # 1 second per block ( 12x faster than Ethereum mainnet )
GAS_LIMIT=300000000 # 300M gas limit ( 10x Ethereum mainnet )

PEER_SERVER="http://localhost:5055"

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "  -h, --help                 Show help message"

  echo "  -f, --config               Config file (default: $WORK_DIR/configs/rpc.config.json)"
  echo "  -p, --peer                 Peer server (default: http://localhost:5055)"
  echo "  -x, --clear                Clear state before starting"
  echo "  -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"
  echo
  echo "Example: $0 -d ~/naive-sequencer-data"
}

clear_data() {
  echo "Clearing data directory: ${NAIVE_RPC_DATA}"
  rm -rf ${NAIVE_RPC_DATA}
  mkdir -p ${NAIVE_RPC_DATA}
}

# Parse command line arguments
while getopts ":hf::xo:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    f|--config )
      RPC_CONFIG_FILE=$OPTARG
      NAIVE_RPC_DATA=$(jq '."data-dir"' -r $RPC_CONFIG_FILE)
      NAIVE_SEQUENCER_DATA="${NAIVE_RPC_DATA}/../naive-sequencer-data"
      ;;
    p|--peer )
      PEER_SERVER=$OPTARG
      ;;
    x|--clear )
      clear_data
      STATE_RESET=1
      ;;
    o|--output )
      RPC_OUTPUT_FILE=$OPTARG
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


# Check if required arguments are present
if [[ -z "${NAIVE_RPC_DATA}" ]]; then
  echo "Missing required argument: --config" 1>&2
  display_help
  exit 1
fi

# Check if data directory exists
if [ ! -d "${NAIVE_RPC_DATA}" ]; then
  echo "Data directory does not exist: ${NAIVE_RPC_DATA}, setting mode to STATE_RESET=1" 1>&2
  clear_data
  STATE_RESET=1
fi

# Also check data dir has data
if [ ! -d "${NAIVE_RPC_DATA}/rpc" ]; then
  echo "Data dir does not contain data: ${NAIVE_RPC_DATA}, setting mode to STATE_RESET=1"
  clear_data
  STATE_RESET=1
fi

PASSWORD_FILE="${NAIVE_RPC_DATA}/password.txt"

if [ $STATE_RESET -eq 1 ]; then
  # Create account for Sequencer
  cp -r ${HOME}/.eth-accounts/ ${NAIVE_RPC_DATA}/keystore/
  mv ${NAIVE_RPC_DATA}/keystore/password.txt ${PASSWORD_FILE}
fi
ACCOUNT1=$(cat ${NAIVE_RPC_DATA}/keystore/* | jq -r '.address' | head -n 1)
echo "Using account: ${ACCOUNT1}"

# Copy over genesis
GENESIS_FILE="${NAIVE_RPC_DATA}/genesis.json"
cp ${NAIVE_SEQUENCER_DATA}/genesis.json ${GENESIS_FILE}
#$WORK_DIR/go-ethereum/build/bin/geth init --datadir ${NAIVE_RPC_DATA} ${GENESIS_FILE}

# Copy over the sequencer l1 address
SEQUENCER_L1_ADDRESS=$(cat "${NAIVE_SEQUENCER_DATA}/sequencer-l1-address.txt" | jq -r '.address')

echo "Starting RPC Node with L1 address: ${SEQUENCER_L1_ADDRESS}"

if [ -z $RPC_OUTPUT_FILE ]; then
  #$WORK_DIR/build/rpc --datadir ${NAIVE_RPC_DATA} --metrics
  $WORK_DIR/build/rpc --config ${RPC_CONFIG_FILE}
else
  $WORK_DIR/build/rpc --config ${RPC_CONFIG_FILE} > $RPC_OUTPUT_FILE 2>&1 &
  echo "Waiting for rpc to start..."
  while true; do
    if grep -q "self=enode://" "${RPC_OUTPUT_FILE}"; then
      break
    fi
    sleep 1
  done
  
  RPC_SERVER=http://$(cat ${RPC_CONFIG_FILE} | jq -r '.host'):$(cat ${RPC_CONFIG_FILE} | jq -r '.port')
  ENODE=$(geth attach --exec admin.nodeInfo.enode $RPC_SERVER)
  echo "ENODE: ${ENODE}"
  
  geth attach --exec "admin.addPeer(${ENODE})" ${PEER_SERVER}
  geth attach --exec "admin.peers" ${PEER_SERVER}
fi
