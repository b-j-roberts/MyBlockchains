build:
	npm run build
	npm run go-gen

deploy-private:
	SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-private

deploy-goerli:
	SEQUENCER_ADDRESS=$(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) npm run deploy-goerli

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
