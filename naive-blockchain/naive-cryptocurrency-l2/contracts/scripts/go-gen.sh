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

mkdir -p ${GO_CONTRACT_DIR}/l1tokenbridge
abigen --bin=builds/contracts_L1TokenBridge_sol_L1TokenBridge.bin --abi=builds/contracts_L1TokenBridge_sol_L1TokenBridge.abi --pkg=l1tokenbridge --out=${GO_CONTRACT_DIR}/l1tokenbridge/l1-token-bridge.go

mkdir -p ${GO_CONTRACT_DIR}/l2tokenbridge
abigen --bin=builds/contracts_L2TokenBridge_sol_L2TokenBridge.bin --abi=builds/contracts_L2TokenBridge_sol_L2TokenBridge.abi --pkg=l2tokenbridge --out=${GO_CONTRACT_DIR}/l2tokenbridge/l2-token-bridge.go

mkdir -p ${GO_CONTRACT_DIR}/basicerc20
abigen --bin=builds/contracts_BasicERC20_sol_BasicERC20.bin --abi=builds/contracts_BasicERC20_sol_BasicERC20.abi --pkg=basicerc20 --out=${GO_CONTRACT_DIR}/basicerc20/basic-erc20.go

mkdir -p ${GO_CONTRACT_DIR}/basicl2erc20
abigen --bin=builds/contracts_BasicL2ERC20_sol_BasicL2ERC20.bin --abi=builds/contracts_BasicL2ERC20_sol_BasicL2ERC20.abi --pkg=basicl2erc20 --out=${GO_CONTRACT_DIR}/basicl2erc20/basic-l2-erc20.go

mkdir -p ${GO_CONTRACT_DIR}/stableerc20
abigen --bin=builds/contracts_StableERC20_sol_StableERC20.bin --abi=builds/contracts_StableERC20_sol_StableERC20.abi --pkg=stableerc20 --out=${GO_CONTRACT_DIR}/stableerc20/stable-erc20.go

mkdir -p ${GO_CONTRACT_DIR}/stablel2erc20
abigen --bin=builds/contracts_StableL2ERC20_sol_StableL2ERC20.bin --abi=builds/contracts_StableL2ERC20_sol_StableL2ERC20.abi --pkg=stablel2erc20 --out=${GO_CONTRACT_DIR}/stablel2erc20/stable-l2-erc20.go
