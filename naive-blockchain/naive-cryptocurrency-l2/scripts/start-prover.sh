#!/bin/bash
#
# This script starts a prover for l1 block bathces posted by the sequencer.


display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -A, --l1-contract-address  l1 contract address (Required)"
  echo "   -k, --keystore             keystore directory for l1 prover address (Required)"

  echo
  #TODO: Examples for all
}

while getopts ":hA:k:" opt; do
  case ${opt} in
    h )
      display_help
      exit 0
      ;;
    A )
      L1_CONTRACT_ADDRESS=$OPTARG
      ;;
    k )
      keystore=$OPTARG
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

if [ -z "${keystore}" ]; then
  echo "Missing required argument: -k" 1>&2
  display_help
  exit 1
fi

proverAddress=$(cat ${keystore}/* | jq -r '.address')

$WORK_DIR/build/prover --l1-contract-address ${L1_CONTRACT_ADDRESS} --prover-address ${proverAddress} --prover-keystore $keystore
