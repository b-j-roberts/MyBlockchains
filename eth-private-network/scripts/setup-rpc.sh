#!/bin/bash
#
# This script sets up an rpc datadir by getting and using the genesis file from another node in the network

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR=$SCRIPT_DIR/..

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"
  echo "   -d, --data                 Geth node data dir (required)"
  echo "   -g, --genesis              Genesis file (required)"

  echo
  echo "Example: $0 -d ~/l1-rpc-data/ -g ~/l1-miner-data/genesis.json"
}

# Parse command line arguments
while getopts ":hd:g:" opt; do
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
  echo "Data dir does not exist: $DATA_DIR, please use setupp command"
  display_help
  exit 1
fi

# Check if genesis file exists
if [ ! -f "$GENESIS_FILE" ]; then
  echo "Genesis file does not exist: $GENESIS_FILE" 1>&2
  display_help
  exit 1
fi

if [ -z "$(ls -A ${HOME}/.eth-accounts)" ]; then
  echo "No accounts found, creating a new one"
  $WORK_DIR/scripts/generate-account.sh -d ${DATA_DIR}
fi
cp -r ${HOME}/.eth-accounts/ $DATA_DIR/keystore
mv ${DATA_DIR}/keystore/password.txt $DATA_DIR/password.txt

# Create Geth Genesis & Init Chain
geth init --datadir $DATA_DIR $GENESIS_FILE
cp $GENESIS_FILE $DATA_DIR/genesis.json 2>/dev/null || :
