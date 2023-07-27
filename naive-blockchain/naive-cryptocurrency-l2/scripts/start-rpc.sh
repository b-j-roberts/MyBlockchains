#!/bin/bash
#
# This script starts the sequencer in a docker container.

STATE_RESET=0

# Defaults give theoretical max throughput of 1800 TPS = ~15 tx/sec ( mainnet ) X 12 X 10
CHAIN_ID=515
PERIOD=1 # 1 second per block ( 12x faster than Ethereum mainnet )
GAS_LIMIT=300000000 # 300M gas limit ( 10x Ethereum mainnet )

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "  -h, --help                 Show help message"

  echo "  -d, --datadir              Data directory (Required)"
  echo "  -k, --keystore             Keystore directory for l1 address (Required)"
  echo "  -B, --l1-bridge-address    L1 Bridge contract address (Required)"

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
while getopts ":hd:k:B:xo:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    d|--datadir )
      NAIVE_RPC_DATA=$OPTARG
      ;;
    k|--keystore )
      L1_KEYSTORE=$OPTARG
      ;;
    B|--l1-bridge-address )
      L1_BRIDGE_ADDRESS=$OPTARG
      ;;
    x|--clear )
      clear_data
      STATE_RESET=1
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

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."
NAIVE_SEQUENCER_DATA="${NAIVE_RPC_DATA}/../naive-sequencer-data"

# Check if required arguments are present
if [[ -z "${NAIVE_RPC_DATA}" || -z ${L1_KEYSTORE} ]]; then
  echo "Missing required argument: --datadir or --keystore" 1>&2
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

if [ -z "${L1_BRIDGE_ADDRESS}" ]; then
  # Copy over the contract address
  cp ${WORK_DIR}/contracts/builds/l1-bridge-address.txt ${NAIVE_RPC_DATA}/l1-bridge-address.txt

  if [ ! -f "${NAIVE_RPC_DATA}/l1-bridge-address.txt" ]; then
    echo "Missing required argument: --l1-bridge-address" 1>&2
    display_help
    exit 1
  fi

  L1_BRIDGE_ADDRESS=$(cat "${NAIVE_RPC_DATA}/l1-bridge-address.txt" | jq -r '.address')
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
for p in  ${L1_KEYSTORE}/*; do cp $p ${NAIVE_RPC_DATA}/sequencer-l1-address.txt; break; done
SEQUENCER_L1_ADDRESS=$(cat "${NAIVE_RPC_DATA}/sequencer-l1-address.txt" | jq -r '.address')

echo "Starting RPC Node with L1 address: ${SEQUENCER_L1_ADDRESS}"

if [ -z $OUTPUT_FILE ]; then
  #$WORK_DIR/build/rpc --datadir ${NAIVE_RPC_DATA} --metrics
  $WORK_DIR/build/rpc --datadir ${NAIVE_RPC_DATA} --addr ${SEQUENCER_L1_ADDRESS} --l1bridgecontract ${L1_BRIDGE_ADDRESS}
else
  $WORK_DIR/build/rpc --datadir ${NAIVE_RPC_DATA} --addr ${SEQUENCER_L1_ADDRESS} --l1bridgecontract ${L1_BRIDGE_ADDRESS} > $OUTPUT_FILE 2>&1 &
  echo "Waiting for rpc to start..."
  while true; do
    if grep -q "self=enode://" "${OUTPUT_FILE}"; then
      break
    fi
    sleep 1
  done
  
  ENODE=$(geth attach --exec admin.nodeInfo.enode ${NAIVE_RPC_DATA}/naive-rpc.ipc)
  echo "ENODE: ${ENODE}"
  
  geth attach --exec "admin.addPeer(${ENODE})" ${NAIVE_SEQUENCER_DATA}/naive-sequencer.ipc
  geth attach --exec "admin.peers" ${NAIVE_SEQUENCER_DATA}/naive-sequencer.ipc
fi
