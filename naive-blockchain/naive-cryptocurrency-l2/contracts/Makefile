build-contracts:
	solcjs --base-path . --include-path node_modules/ --bin --abi --optimize -o builds/ contracts/*.sol

build-bridge:
	solcjs --base-path . --include-path node_modules/ --bin --abi --optimize -o builds/ contracts/bridge/*.sol

build-tokens:
	solcjs --base-path . --include-path node_modules/ --bin --abi --optimize -o builds/ contracts/tokens/erc20/*.sol
	solcjs --base-path . --include-path node_modules/ --bin --abi --optimize -o builds/ contracts/tokens/erc721/*.sol

build: build-contracts build-bridge build-tokens
	./scripts/go-gen.sh

deploy-l1-contracts:
	echo "Deploying contracts to L1 chain"
	SEQUENCER_ADDRESS=${SEQUENCER_ADDRESS} npm run deploy-l1-contracts
	SEQUENCER_ADDRESS=${SEQUENCER_ADDRESS} npm run deploy-l1-tokens

#TODO: Consider moving these up a level to the outer Makefile
SEQUENCER_CONFIG_FILE ?= ../configs/sequencer.config.json
SEQUENCER_DATA_DIR ?= ~/naive-sequencer-data
SEQUENCER_ADDRESS ?= $(shell cat ${SEQUENCER_DATA_DIR}/sequencer-l1-address.txt  | jq -r .address)

deploy-l2-contracts:
	SEQUENCER_ADDRESS=${SEQUENCER_ADDRESS} IPC_PATH=${SEQUENCER_DATA_DIR}/naive-sequencer.ipc npm run deploy-l2-contracts
	cp builds/l2-bridge-address.txt ${SEQUENCER_DATA_DIR}/contracts/l2-bridge-address.txt
	cp builds/l2-token-bridge-address.txt ${SEQUENCER_DATA_DIR}/contracts/l2-token-bridge-address.txt

L2_TOKEN_BRIDGE_ADDRESS ?= $(shell cat ${SEQUENCER_DATA_DIR}/contracts/l2-token-bridge-address.txt  | jq -r .address)

deploy-l2-tokens:
	SEQUENCER_ADDRESS=${SEQUENCER_ADDRESS} L2_TOKEN_BRIDGE_ADDRESS=${L2_TOKEN_BRIDGE_ADDRESS} IPC_PATH=${SEQUENCER_DATA_DIR}/naive-sequencer.ipc npm run deploy-l2-tokens
	cp builds/*-address.txt ${SEQUENCER_DATA_DIR}/contracts/

allow-tokens:
	ACCOUNT_PASS=password ../build/allow-tokens --sender ${SEQUENCER_ADDRESS} --config ${SEQUENCER_CONFIG_FILE}

#TODO: Test on goerli with latest contracts
deploy-goerli:
	SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-goerli

#TODO: Update to test with latest contracts
test-private:
	npm run test-private

test-goerli:
	npm run test-goerli

go-gen:
	./scripts/go-gen.sh

clean:
	rm -rf builds/ go/

# Define the default rule
.DEFAULT_GOAL := build
