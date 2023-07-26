build:
	npm run build
	npm run go-gen

deploy-private:
	SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-private

deploy-goerli:
	SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-goerli

deploy-l2:
	SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-l2

deploy-l2-erc20:
	L2_TOKEN_BRIDGE_ADDRESS=$(shell cat builds/l2-token-bridge-address.txt  | jq -r .address) SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-l2-erc20

allow-erc20:
	go run scripts/allow-erc20.go --sender $(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) --erc20 $(shell cat builds/basic-erc20-address.txt | jq -r .address) --l2erc20 $(shell cat builds/l2-basic-erc20-address.txt | jq -r .address) --tokenBridge $(shell cat builds/l1-token-bridge-address.txt | jq -r .address) --l2tokenBridge $(shell cat builds/l2-token-bridge-address.txt | jq -r .address)  --keystore ~/naive-sequencer-data/keystore/ --stableErc20 $(shell cat builds/stable-erc20-address.txt | jq -r .address) --l2stableErc20 $(shell cat builds/l2-stable-erc20-address.txt | jq -r .address)

test-private:
	npm run test-private

test-goerli:
	npm run test-goerli

go-gen:
	npm run go-gen

clean:
	rm -rf builds/ go/

# Define the default rule
.DEFAULT_GOAL := build
