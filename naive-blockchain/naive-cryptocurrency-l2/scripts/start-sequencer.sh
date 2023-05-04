#!/bin/bash

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"
NAIVE_RPC_NODE_DATA="$HOME/naive-rpc-data"

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR="${BASE_DIR}/.."

# Cleanup previous runs
rm -rf ${NAIVE_SEQUENCER_DATA}
mkdir -p ${NAIVE_SEQUENCER_DATA}

rm -rf ${NAIVE_RPC_NODE_DATA}
mkdir -p ${NAIVE_RPC_NODE_DATA}

# Start the sequencer
cp ${BASE_DIR}/scripts/config/password.txt ${NAIVE_SEQUENCER_DATA}/password.txt
${BASE_DIR}/go-ethereum/build/bin/geth account new --datadir ${NAIVE_SEQUENCER_DATA} --password ${NAIVE_SEQUENCER_DATA}/password.txt
ACCOUNT1=$(cat ${NAIVE_SEQUENCER_DATA}/keystore/* | jq -r '.address' | head -n 1)

cp ${BASE_DIR}/scripts/config/password.txt ${NAIVE_RPC_NODE_DATA}/password.txt
${BASE_DIR}/go-ethereum/build/bin/geth account new --datadir ${NAIVE_RPC_NODE_DATA} --password ${NAIVE_RPC_NODE_DATA}/password.txt
ACCOUNT2=$(cat ${NAIVE_RPC_NODE_DATA}/keystore/* | jq -r '.address' | tail -n 1)

${BASE_DIR}/scripts/generate-genesis.sh ${ACCOUNT1} ${ACCOUNT2}
cp ${BASE_DIR}/build/genesis.json ${NAIVE_SEQUENCER_DATA}/genesis.json
cp ${BASE_DIR}/build/genesis.json ${NAIVE_RPC_NODE_DATA}/genesis.json

#${BASE_DIR}/go-ethereum/build/bin/geth init --datadir ${NAIVE_SEQUENCER_DATA} ${BASE_DIR}/build/genesis.json
go run ${BASE_DIR}/cmd/sequencer/sequencer.go

#geth init --datadir ${NAIVE_RPC_NODE_DATA} ${BASE_DIR}/build/genesis.json
