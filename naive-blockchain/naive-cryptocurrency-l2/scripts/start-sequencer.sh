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
  echo "  -A, --l1-contract-address  L1 contract address (Required)"

  echo "  -c, --chainid              Chain ID (default: 515)"
  echo "  -p, --period               Block period in seconds (default: 1)"
  echo "  -g, --gaslimit             Gas limit per block (default: 300000000)"

  echo "  -x, --clear                Clear state before starting"
  echo "  -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"
  echo
  echo "Example: $0 -d ~/naive-sequencer-data"
}

clear_data() {
  echo "Clearing data directory: ${NAIVE_SEQUENCER_DATA}"
  rm -rf ${NAIVE_SEQUENCER_DATA}
  mkdir -p ${NAIVE_SEQUENCER_DATA}
}

# Parse command line arguments
while getopts ":hd:k:A:c:p:g:xo:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    d|--datadir )
      NAIVE_SEQUENCER_DATA=$OPTARG
      ;;
    c|--chainid )
      CHAIN_ID=$OPTARG
      ;;
    p|--period )
      PERIOD=$OPTARG
      ;;
    g|--gaslimit )
      GAS_LIMIT=$OPTARG
      ;;
    k|--keystore )
      L1_KEYSTORE=$OPTARG
      ;;
    A|--l1-contract-address )
      L1_CONTRACT_ADDRESS=$OPTARG
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

# Check if required arguments are present
if [[ -z "${NAIVE_SEQUENCER_DATA}" || -z ${L1_KEYSTORE} ]]; then
  echo "Missing required argument: --datadir or --keystore" 1>&2
  display_help
  exit 1
fi

# Check if data directory exists
if [ ! -d "${NAIVE_SEQUENCER_DATA}" ]; then
  echo "Data directory does not exist: ${NAIVE_SEQUENCER_DATA}, setting mode to STATE_RESET=1" 1>&2
  clear_data
  STATE_RESET=1
fi

# Also check data dir has data
if [ ! -d "${NAIVE_SEQUENCER_DATA}/sequencer" ]; then
  echo "Data dir does not contain data: ${NAIVE_SEQUENCER_DATA}, setting mode to STATE_RESET=1"
  clear_data
  STATE_RESET=1
fi

if [ -z "${L1_CONTRACT_ADDRESS}" ]; then
  # Copy over the contract address
  cp ${WORK_DIR}/contracts/builds/contract-address.txt ${NAIVE_SEQUENCER_DATA}/l1-contract-address.txt

  if [ ! -f "${NAIVE_SEQUENCER_DATA}/l1-contract-address.txt" ]; then
    echo "Missing required argument: --l1-contract-address" 1>&2
    display_help
    exit 1
  fi

  # Start the sequencer
  L1_CONTRACT_ADDRESS=$(cat "${NAIVE_SEQUENCER_DATA}/l1-contract-address.txt" | jq -r '.address')
fi

PASSWORD_FILE="${NAIVE_SEQUENCER_DATA}/password.txt"

if [ $STATE_RESET -eq 1 ]; then
  # Create account for Sequencer
  ACCOUNT_PASS=${ACCOUNT_PASS:-password}
  touch ${PASSWORD_FILE}
  echo ${ACCOUNT_PASS} > ${PASSWORD_FILE}
  $WORK_DIR/go-ethereum/build/bin/geth account new --datadir ${NAIVE_SEQUENCER_DATA} --password ${PASSWORD_FILE}
fi
ACCOUNT1=$(cat ${NAIVE_SEQUENCER_DATA}/keystore/* | jq -r '.address' | head -n 1)

if [ $STATE_RESET -eq 1 ]; then
  # Create L2 Genesis & Init Chain
  GENESIS_FILE="${NAIVE_SEQUENCER_DATA}/genesis.json"

  echo "Creating L2 Genesis file: ${GENESIS_FILE} with account: ${ACCOUNT1} & balance: 10000000000000000000 wei (10 ETH)"
  $WORK_DIR/scripts/generate-genesis.sh -a ${ACCOUNT1} -b 10000000000000000000 -o ${GENESIS_FILE} -p ${PERIOD} -g ${GAS_LIMIT} -c ${CHAIN_ID}
  $WORK_DIR/go-ethereum/build/bin/geth init --datadir ${NAIVE_SEQUENCER_DATA} ${GENESIS_FILE}

  # Copy over the sequencer l1 address
  for p in  ${L1_KEYSTORE}/*; do cp $p ${NAIVE_SEQUENCER_DATA}/sequencer-l1-address.txt; break; done
fi

SEQUENCER_L1_ADDRESS=$(cat "${NAIVE_SEQUENCER_DATA}/sequencer-l1-address.txt" | jq -r '.address')

echo "Starting sequencer with L1 contract address: ${L1_CONTRACT_ADDRESS} & L1 sequencer address: ${SEQUENCER_L1_ADDRESS}"

if [ -z $OUTPUT_FILE ]; then
  $WORK_DIR/build/sequencer --datadir ${NAIVE_SEQUENCER_DATA} --l1contract ${L1_CONTRACT_ADDRESS} --sequencer ${SEQUENCER_L1_ADDRESS} --sequencerkeystore ${L1_KEYSTORE} --metrics --metrics.expensive
else
  $WORK_DIR/build/sequencer --datadir ${NAIVE_SEQUENCER_DATA} --l1contract ${L1_CONTRACT_ADDRESS} --sequencer ${SEQUENCER_L1_ADDRESS} --sequencerkeystore ${L1_KEYSTORE} --metrics --metrics.expensive > $OUTPUT_FILE 2>&1 &
fi
