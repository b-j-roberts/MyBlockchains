#!/bin/bash
#
# This script starts the smart contract metrics exporter for l1 blockchain DA contract


display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -A, --l1-tx-storage-address  l1 contract address (Required)"
  echo "   -B, --l1-bridge-address    l1 bridge address (Required)"
  echo "   -M, --l2-bridge-address    l2 bridge address (Required)"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  #TODO: Examples for all
}

while getopts ":hA:B:M:o:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    A|--l1-tx-storage-address )
      L1_CONTRACT_ADDRESS=$OPTARG
      ;;
    B|--l1-bridge-address )
      L1_BRIDGE_ADDRESS=$OPTARG
      ;;
    M|--l2-bridge-address )
      L2_BRIDGE_ADDRESS=$OPTARG
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

if [ -z "${L1_CONTRACT_ADDRESS}" ]; then
  # Try and copy over address from build
  L1_CONTRACT_ADDRESS=$(cat $WORK_DIR/contracts/builds/tx-storage-address.txt | jq -r '.address')

  if [ -z "${L1_CONTRACT_ADDRESS}" ]; then
    echo "Missing required argument: -A" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L1_BRIDGE_ADDRESS}" ]; then
  # Try and copy over address from build
  L1_BRIDGE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l1-bridge-address.txt | jq -r '.address')

  if [ -z "${L1_BRIDGE_ADDRESS}" ]; then
    echo "Missing required argument: -B" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L2_BRIDGE_ADDRESS}" ]; then
  # Try and copy over address from build
  L2_BRIDGE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l2-bridge-address.txt | jq -r '.address')

  if [ -z "${L2_BRIDGE_ADDRESS}" ]; then
    echo "Missing required argument: -M" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${OUTPUT_FILE}" ]; then
  $WORK_DIR/build/smart-contract-metrics --l1-tx-storage-address ${L1_CONTRACT_ADDRESS} --l1-bridge-address ${L1_BRIDGE_ADDRESS} --l2-bridge-address ${L2_BRIDGE_ADDRESS}
else
  $WORK_DIR/build/smart-contract-metrics --l1-tx-storage-address ${L1_CONTRACT_ADDRESS} --l1-bridge-address ${L1_BRIDGE_ADDRESS} --l2-bridge-address ${L2_BRIDGE_ADDRESS} > ${OUTPUT_FILE} 2>&1 &
fi
