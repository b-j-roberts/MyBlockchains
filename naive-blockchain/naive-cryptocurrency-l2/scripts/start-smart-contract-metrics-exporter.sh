#!/bin/bash
#
# This script starts the smart contract metrics exporter for l1 blockchain DA contract


display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -A, --l1-contract-address  l1 contract address (Required)"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  #TODO: Examples for all
}

while getopts ":hA:k:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    A|--l1-contract-address )
      L1_CONTRACT_ADDRESS=$OPTARG
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
  L1_CONTRACT_ADDRESS=$(cat $WORK_DIR/contracts/builds/contract-address.txt | jq -r '.address')

  if [ -z "${L1_CONTRACT_ADDRESS}" ]; then
    echo "Missing required argument: -A" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${OUTPUT_FILE}" ]; then
  $WORK_DIR/build/smart-contract-metrics --l1-contract-address ${L1_CONTRACT_ADDRESS}
else
  $WORK_DIR/build/smart-contract-metrics --l1-contract-address ${L1_CONTRACT_ADDRESS} > ${OUTPUT_FILE} 2>&1 &
fi
