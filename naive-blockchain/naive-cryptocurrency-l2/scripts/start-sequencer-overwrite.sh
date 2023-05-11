#!/bin/bash

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR="${BASE_DIR}/.."

# Cleanup previous runs
rm -rf ${NAIVE_SEQUENCER_DATA}
mkdir -p ${NAIVE_SEQUENCER_DATA}

# Start the sequencer
cp ${BASE_DIR}/scripts/config/password.txt ${NAIVE_SEQUENCER_DATA}/password.txt
${BASE_DIR}/go-ethereum/build/bin/geth account new --datadir ${NAIVE_SEQUENCER_DATA} --password ${NAIVE_SEQUENCER_DATA}/password.txt
ACCOUNT1=$(cat ${NAIVE_SEQUENCER_DATA}/keystore/* | jq -r '.address' | head -n 1)

${BASE_DIR}/scripts/generate-genesis.sh ${ACCOUNT1}
cp ${BASE_DIR}/build/genesis.json ${NAIVE_SEQUENCER_DATA}/genesis.json
cp ${BASE_DIR}/contracts/builds/contract-address.txt ${NAIVE_SEQUENCER_DATA}/l1-contract-address.txt
for p in  ${BASE_DIR}/../../eth-private-network/data/keystore/*; do cp $p ${NAIVE_SEQUENCER_DATA}/sequencer-l1-address.txt; break; done

L1_CONTRACT_ADDRESS=$(cat "${NAIVE_SEQUENCER_DATA}/l1-contract-address.txt" | jq -r '.address')
SEQUENCER_L1_ADDRESS=$(cat "${NAIVE_SEQUENCER_DATA}/sequencer-l1-address.txt" | jq -r '.address')

go run ${BASE_DIR}/cmd/sequencer/ --datadir ${NAIVE_SEQUENCER_DATA} --l1contract ${L1_CONTRACT_ADDRESS} --sequencer ${SEQUENCER_L1_ADDRESS} --sequencerkeystore ${BASE_DIR}/../../eth-private-network/data/keystore/
