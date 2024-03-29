#!/bin/bash


# This script is used to quickly launch evenrything needed to run the L2, bridge, and test the system.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."
OUTPUT_DIR="${HOME}/blockchain-logs"
BRIDGE=false
FULL=false

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"
  echo "   -o, --output-dir           output directory for logs ( These will be cleaned )"
  echo "   -b, --bridge               bridge eth & tokens"
  echo "   -f, --full                 full launch ( including l1 & l2 rpc servers )"

  echo
  echo "Example: $0"
}

#TODO: Options including noclean launch, docker launch, etc.
while getopts ":h:o:bf" opt; do
  case ${opt} in
    h )
      display_help
      exit 0
      ;;
    o )
      OUTPUT_DIR=$OPTARG
      ;;
    b )
      BRIDGE=true
      ;;
    f )
      FULL=true
      ;;
    \? )
      echo "Invalid Option: -$OPTARG" 1>&2
      display_help
      exit 1
      ;;
  esac
done

rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

echo "Launching L1"
L1_MINER_LOGS=${OUTPUT_DIR}/l1-miner.logs
touch $L1_MINER_LOGS
cd ${WORK_DIR}/../../eth-private-network/ && OUTPUT_FILE=${L1_MINER_LOGS} make run-miner-daemon

# Wait for miner to be ready
sleep 5

if [ "$FULL" = true ]; then
  echo "Launching L1 RPC"
  L1_RPC_SERVER_LOGS=${OUTPUT_DIR}/l1-rpc.logs
  touch $L1_RPC_SERVER_LOGS
  cd ${WORK_DIR}/../../eth-private-network/ && OUTPUT_FILE=${L1_RPC_SERVER_LOGS} make run-rpc-daemon
  
  # Wait for rpc to be ready
  sleep 5
  
  echo "Connecting Miner & RPC"
  cd ${WORK_DIR}/../../eth-private-network/ && make connect-peers
fi

echo "Deploying Contracts to L1"
cd ${WORK_DIR} && make deploy-l1-contracts

echo "Launching L2"
L2_SEQUENCER_LOGS=${OUTPUT_DIR}/l2-sequencer.logs
touch $L2_SEQUENCER_LOGS

SEQUENCER_OUTPUT_FILE=${L2_SEQUENCER_LOGS} ${SCRIPT_DIR}/start-sequencer-retry.sh

if [ "$FULL" = true ]; then
  echo "Launching L2 RPC"
  L2_RPC_SERVER_LOGS=${OUTPUT_DIR}/l2-rpc.logs
  touch $L2_RPC_SERVER_LOGS
  cd ${WORK_DIR} && RPC_OUTPUT_FILE=${L2_RPC_SERVER_LOGS} make run-rpc
fi

echo "Deploying Contracts to L2"
cd ${WORK_DIR} && make deploy-l2-contracts

echo "Starting Prover"
L2_PROVER_LOGS=${OUTPUT_DIR}/l2-prover.logs
cd ${WORK_DIR} && PROVER_OUTPUT_FILE=${L2_PROVER_LOGS} make run-prover

echo "Starting metrics server"
METRICS_LOGS=${OUTPUT_DIR}/metrics.logs
cd ${WORK_DIR} && METRICS_OUTPUT_FILE=${METRICS_LOGS} make run-smart-contract-metrics

if [ "$BRIDGE" = true ]; then
  echo "Bridge things over..."
  echo "Bridge eth to l2"
  cd ${WORK_DIR} && make bridge-eth-to-l2
  cd ${WORK_DIR} && make bridge-eth-to-l2
  cd ${WORK_DIR} && make bridge-eth-to-l2
  cd ${WORK_DIR} && make bridge-eth-to-l2
  cd ${WORK_DIR} && make bridge-eth-to-l2

  sleep 5
  
  echo "Bridge eth to l1"
  cd ${WORK_DIR} && make bridge-eth-to-l1
  cd ${WORK_DIR} && make bridge-eth-to-l1
  
  echo "Bridge basic erc20 to l2"
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l2
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l2
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l2
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l2
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l2

  sleep 5
  
  echo "Bridge basic erc20 to l1"
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l1
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l1
  cd ${WORK_DIR} && make bridge-basic-erc20-to-l1
  
  echo "Bridge stable erc20 to l2"
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l2
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l2
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l2
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l2
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l2

  sleep 5
  
  echo "Bridge stable erc20 to l1"
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l1
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l1
  cd ${WORK_DIR} && make bridge-stable-erc20-to-l1

  echo "Bridging basic erc721 to l2"
  cd ${WORK_DIR} && make bridge-basic-erc721-to-l2

  sleep 5

  echo "Bridging basic erc721 to l1"
  cd ${WORK_DIR} && make bridge-basic-erc721-to-l1

  echo "Bridging special erc721 to l2"
  cd ${WORK_DIR} && make bridge-special-erc721-to-l2
fi
