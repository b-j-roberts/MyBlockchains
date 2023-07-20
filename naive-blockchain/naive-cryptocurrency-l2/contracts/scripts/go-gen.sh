#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

GO_CONTRACT_DIR=${SCRIPT_DIR}/../go/

rm -rf ${GO_CONTRACT_DIR}
mkdir -p ${GO_CONTRACT_DIR}

mkdir -p ${GO_CONTRACT_DIR}/txstorage
abigen --bin=builds/contracts_TransactionStorage_sol_TransactionStorage.bin --abi=builds/contracts_TransactionStorage_sol_TransactionStorage.abi --pkg=txstorage --out=${GO_CONTRACT_DIR}/txstorage/transaction-storage.go

mkdir -p ${GO_CONTRACT_DIR}/l1bridge
abigen --bin=builds/contracts_L1Bridge_sol_L1Bridge.bin --abi=builds/contracts_L1Bridge_sol_L1Bridge.abi --pkg=l1bridge --out=${GO_CONTRACT_DIR}/l1bridge/l1-bridge.go

mkdir -p ${GO_CONTRACT_DIR}/l2bridge
abigen --bin=builds/contracts_L2Bridge_sol_L2Bridge.bin --abi=builds/contracts_L2Bridge_sol_L2Bridge.abi --pkg=l2bridge --out=${GO_CONTRACT_DIR}/l2bridge/l2-bridge.go
