#!/bin/bash

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR="${BASE_DIR}/.."

l1ContractAddressFile=${NAIVE_SEQUENCER_DATA}/l1-contract-address.txt
proverL1AddressFile=${NAIVE_SEQUENCER_DATA}/sequencer-l1-address.txt

go run ${BASE_DIR}/cmd/prover/ --l1-contract-address-file ${l1ContractAddressFile} --prover-address-file ${proverL1AddressFile} --prover-keystore ${BASE_DIR}/../../eth-private-network/data/keystore/
