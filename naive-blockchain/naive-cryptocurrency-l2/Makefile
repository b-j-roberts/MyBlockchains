.PHONY: all contracts sequencer prover rpc verifier allow-tokens clean smart-contract-metrics bridge metrics-test
.DEFAULT_GOAL := all

# BUILD COMMANDS
all: geth contracts sequencer prover rpc verifier smart-contract-metrics allow-tokens bridge

clean:
	cd go-ethereum && make clean
	rm -rf build/
	cd contracts && make clean

geth:
	cd go-ethereum && make geth

contracts:
	cd contracts && make build

sequencer:
	go build -o build/sequencer ./cmd/sequencer

prover:
	go build -o build/prover ./cmd/prover

rpc:
	go build -o build/rpc ./cmd/rpc

verifier:
	go build -o build/verifier ./cmd/verifier

smart-contract-metrics:
	go build -o build/smart-contract-metrics ./cmd/smart-contract-metrics

allow-tokens:
	go build -o build/allow-tokens ./cmd/allow-tokens

bridge:
	go build -o build/bridge ./cmd/bridge

metrics-test:
	go build -o build/metrics-test ./cmd/metrics-test

docker-build:
	make clean
	docker build -f docker/Dockerfile.sequencer . -t brandonjroberts/naive-l2-sequencer:latest
	docker build -f docker/Dockerfile.prover . -t brandonjroberts/naive-l2-prover:latest
	docker build -f docker/Dockerfile.smart-contract-exporter . -t brandonjroberts/naive-l2-smart-contract-exporter:latest
	#TODO: naive-l2-rpc

docker-push:
	docker push brandonjroberts/naive-l2-sequencer:latest
	docker push brandonjroberts/naive-l2-prover:latest
	docker push brandonjroberts/naive-l2-smart-contract-exporter:latest


# RUN COMMANDS

#TODO: Do similar setup on eth private network
## Sequencer
SEQUENCER_CONFIG_FILE ?= configs/sequencer.config.json
SEQUENCER_XTRA_FLAGS ?= -x

run-sequencer:
ifeq ($(SEQUENCER_OUTPUT_FILE),)
	./scripts/start-sequencer.sh -f ${SEQUENCER_CONFIG_FILE} ${SEQUENCER_XTRA_FLAGS}
else
	./scripts/start-sequencer.sh -f ${SEQUENCER_CONFIG_FILE} -o ${SEQUENCER_OUTPUT_FILE} ${SEQUENCER_XTRA_FLAGS}
endif

run-sequencer-daemon:
	SEQUENCER_OUTPUT_FILE=${HOME}/naive-sequencer-data/sequencer.log make run-sequencer

run-sequencer-noclean:
	SEQUENCER_XTRA_FLAGS= make run-sequencer

docker-run-sequencer:
	# TODO: Allow docker config file path overwrite with env var, since hardcoded in dockerfile
	# TODO: Restart w/ no clean
	# GET CONTRACT, CONFIG, AND ACCOUNT DATA FROM ETH PRIVATE NETWORK
	rm -rf ~/tmp/l1-docker-data-for-sequencer/
	mkdir -p ~/tmp/l1-docker-data-for-sequencer/
	docker cp eth-private-miner:/root/.eth-accounts/ ~/tmp/l1-docker-data-for-sequencer/
	mkdir -p ~/tmp/l1-docker-data-for-sequencer/contracts/
	cp contracts/builds/*-address.txt ~/tmp/l1-docker-data-for-sequencer/contracts/
	cp configs/docker-sequencer.config.json ~/tmp/l1-docker-data-for-sequencer/sequencer.config.json
	# SETUP DATADIR & CONFIG
	mkdir -p ~/docker-sequencer/
	# RUN SEQUENCER
	docker run -it -p 5055:5055 -p 6160:6160 -p 30303:30303 --network host --rm --name naive-sequencer -e ACCOUNT_PASS=password \
	  --mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-sequencer/,target=/root/ \
	  -v ~/docker-sequencer/:/sequencer-data brandonjroberts/naive-l2-sequencer:latest

## Prover
PROVER_CONFIG_FILE ?= configs/sequencer.config.json

run-prover:
ifeq ($(PROVER_OUTPUT_FILE),)
	ACCOUNT_PASS=password ./scripts/start-prover.sh -c ${PROVER_CONFIG_FILE}
else
	ACCOUNT_PASS=password ./scripts/start-prover.sh -c ${PROVER_CONFIG_FILE} -o ${PROVER_OUTPUT_FILE}
endif

docker-run-prover:
	# GET CONTRACTS, CONFIG, AND ACCOUNT DATA FROM SEQUENCER
	rm -rf ~/tmp/l1-docker-data-for-prover/
	mkdir -p ~/tmp/l1-docker-data-for-prover/
	docker cp eth-private-miner:/root/.eth-accounts/ ~/tmp/l1-docker-data-for-prover/
	docker cp naive-sequencer:/sequencer-data/contracts/ ~/tmp/l1-docker-data-for-prover/
	cp configs/docker-prover.config.json ~/tmp/l1-docker-data-for-prover/sequencer.config.json
	# RUN PROVER
	docker run -it -p 6162:6162 --network host --rm --name naive-prover -e ACCOUNT_PASS=password \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-prover/,target=/root/ \
		brandonjroberts/naive-l2-prover:latest

## RPC
RPC_CONFIG_FILE ?= configs/rpc.config.json
RPC_OUTPUT_FILE ?= ~/naive-rpc-data/log.txt

run-rpc:
	./scripts/start-rpc.sh -f ${RPC_CONFIG_FILE} -x -o ${RPC_OUTPUT_FILE}

run-rpc-noclean:
	./scripts/start-rpc.sh -f ${RPC_CONFIG_FILE} -o ${RPC_OUTPUT_FILE}

## Smart Contract Exporter
SMART_CONTRACT_EXPORTER_CONFIG_FILE ?= configs/sequencer.config.json
SEQUENCER_ADDRESS ?= $(shell cat ~/l1-miner-data/keystore/UTC* | jq -r '.address')

run-smart-contract-metrics:
ifeq ($(METRICS_OUTPUT_FILE),)
	./scripts/start-smart-contract-metrics-exporter.sh -s ${SEQUENCER_ADDRESS} -c ${SMART_CONTRACT_EXPORTER_CONFIG_FILE}
else
	./scripts/start-smart-contract-metrics-exporter.sh -s ${SEQUENCER_ADDRESS} -c ${SMART_CONTRACT_EXPORTER_CONFIG_FILE} -o ${METRICS_OUTPUT_FILE}
endif

run-smart-contract-metrics-daemon:
	METRICS_OUTPUT_FILE=/dev/null make run-smart-contract-metrics

docker-run-smart-contract-metrics:
	# GET CONTRACTS & CONFIG FROM SEQUENCER
	rm -rf ~/tmp/l1-docker-data-for-smart-contract-metrics/
	mkdir -p ~/tmp/l1-docker-data-for-smart-contract-metrics/contracts/
	cp configs/docker-prover.config.json ~/tmp/l1-docker-data-for-smart-contract-metrics/sequencer.config.json
	docker cp eth-private-miner:/root/.eth-accounts/ ~/tmp/l1-docker-data-for-smart-contract-metrics/
	docker cp naive-sequencer:/sequencer-data/contracts/ ~/tmp/l1-docker-data-for-smart-contract-metrics/
	#	RUN SMART CONTRACT METRICS EXPORTER
	docker run -it -p 6169:6169 -v /home/brandon/docker-sequencer:/sequencer-data/ \
		--network host --rm --name naive-smart-contract-metrics-exporter \
		-e SEQUENCER_ADDRESS=$(shell cat ~/tmp/l1-docker-data-for-smart-contract-metrics/.eth-accounts/UTC* | jq -r '.address') \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-smart-contract-metrics/,target=/root/ \
		brandonjroberts/naive-l2-smart-contract-exporter:latest


# CONTRACT COMMANDS
deploy-l1-contracts:
	rm -f contracts/builds/*-address.txt
	cd contracts && SEQUENCER_ADDRESS=${SEQUENCER_ADDRESS} make deploy-l1-contracts

docker-deploy-l1-contracts:
	SEQUENCER_ADDRESS=$(shell cat ~/docker-l1-miner-data/keystore/UTC* | jq -r '.address') make deploy-l1-contracts

deploy-l2-contracts:
	cd contracts && make deploy-l2-contracts
	cd contracts && make deploy-l2-tokens
	cd contracts && make allow-tokens

docker-deploy-l2-contracts:
	ls contracts/builds/*-address.txt | xargs -I {} docker cp {} naive-sequencer:/app/contracts/builds/
	docker exec -it naive-sequencer sh -c "cd contracts && SEQUENCER_DATA_DIR=/sequencer-data/ SEQUENCER_CONFIG_FILE=/root/sequencer.config.json make deploy-l2-contracts"
	docker exec -it naive-sequencer sh -c "cd contracts && SEQUENCER_DATA_DIR=/sequencer-data/ SEQUENCER_CONFIG_FILE=/root/sequencer.config.json make deploy-l2-tokens"
	docker exec -it naive-sequencer sh -c "cd contracts && SEQUENCER_DATA_DIR=/sequencer-data/ SEQUENCER_CONFIG_FILE=/root/sequencer.config.json make allow-tokens"

allow-deployed-tokens:
	cd contracts && make allow-tokens


# BRIDGE COMMANDS
SEQUENCER_DATA_DIR ?= ~/naive-sequencer-data/

bridge-eth-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 1000000000000000 --config ${SEQUENCER_CONFIG_FILE}

docker-bridge-eth-to-l2:
	rm -rf ~/tmp/l1-docker-data-for-bridge/
	mkdir -p ~/tmp/l1-docker-data-for-bridge/
	docker cp eth-private-miner:/l1-miner-data/keystore/ ~/tmp/l1-docker-data-for-bridge/
	cp ~/tmp/l1-docker-data-for-bridge/keystore/* ~/.transactor
	#TODO: Must run this twice because cat is being used on a file that is being written to
	SEQUENCER_ADDRESS=$(shell cat ~/tmp/l1-docker-data-for-bridge/keystore/UTC*  | jq -r .address) SEQUENCER_CONFIG_FILE=configs/docker-sequencer.config.json make bridge-eth-to-l2
	rm -rf ~/tmp/l1-docker-data-for-bridge/

bridge-eth-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 10000000000000 --config ${SEQUENCER_CONFIG_FILE} --to-l1

docker-bridge-eth-to-l1:
	rm -rf ~/tmp/l1-docker-data-for-bridge/
	mkdir -p ~/tmp/l1-docker-data-for-bridge/
	docker cp naive-sequencer:/sequencer-data/keystore/ ~/tmp/l1-docker-data-for-bridge/
	cp ~/tmp/l1-docker-data-for-bridge/keystore/* ~/.transactor
	SEQUENCER_ADDRESS=$(shell cat ~/tmp/l1-docker-data-for-bridge/keystore/UTC*  | jq -r .address) SEQUENCER_CONFIG_FILE=configs/docker-sequencer.config.json make bridge-eth-to-l1
	rm -rf ~/tmp/l1-docker-data-for-bridge/

bridge-basic-erc20-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 100000 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc

docker-bridge-basic-erc20-to-l2:
	rm -rf ~/tmp/l1-docker-data-for-bridge/
	mkdir -p ~/tmp/l1-docker-data-for-bridge/
	docker cp eth-private-miner:/l1-miner-data/keystore/ ~/tmp/l1-docker-data-for-bridge/
	cp ~/tmp/l1-docker-data-for-bridge/keystore/* ~/.transactor
	SEQUENCER_ADDRESS=$(shell cat ~/tmp/l1-docker-data-for-bridge/keystore/UTC*  | jq -r .address) SEQUENCER_CONFIG_FILE=configs/bridge.config.json make bridge-basic-erc20-to-l2
	rm -rf ~/tmp/l1-docker-data-for-bridge/

bridge-basic-erc20-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 1000 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc --to-l1

docker-bridge-basic-erc20-to-l1:
	rm -rf ~/tmp/l1-docker-data-for-bridge/
	mkdir -p ~/tmp/l1-docker-data-for-bridge/
	docker cp naive-sequencer:/sequencer-data/keystore/ ~/tmp/l1-docker-data-for-bridge/
	cp ~/tmp/l1-docker-data-for-bridge/keystore/* ~/.transactor
	SEQUENCER_ADDRESS=$(shell cat ~/tmp/l1-docker-data-for-bridge/keystore/UTC*  | jq -r .address) SEQUENCER_CONFIG_FILE=configs/bridge.config.json make bridge-basic-erc20-to-l1
	rm -rf ~/tmp/l1-docker-data-for-bridge/

bridge-stable-erc20-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 100000 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/stable-erc20-address.txt | jq -r .address) --is-erc

bridge-stable-erc20-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 1000 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/stable-erc20-address.txt | jq -r .address) --is-erc --to-l1

bridge-basic-erc721-to-l2:
	for number in 0 1 2; do \
		ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value $$number --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/basic-erc721-address.txt | jq -r .address) --is-erc; \
	done

bridge-basic-erc721-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 0 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/basic-erc721-address.txt | jq -r .address) --is-erc --to-l1

bridge-special-erc721-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 0 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/special-erc721-address.txt | jq -r .address) --is-erc

bridge-special-erc721-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address ${SEQUENCER_ADDRESS} --value 0 --config ${SEQUENCER_CONFIG_FILE} --token $(shell cat contracts/builds/special-erc721-address.txt | jq -r .address) --is-erc --to-l1

## Quick launch
quick-launch:
	./scripts/quick-launch.sh -b

quick-launch-full:
	./scripts/quick-launch.sh -b -f

quick-clean:
	./scripts/quick-clean.sh

run-metrics-test:
	./build/metrics-test

#TODO: Clean this up and make it more effective
#
#watch-smart-contract:
#	./scripts/smart-contract-watch.sh -a $(shell cat contracts/builds/tx-storage-address.txt | jq -r .address)
#
#watch-l2-bridge-contract:
#	./scripts/smart-contract-watch.sh -a $(shell cat contracts/builds/l2-bridge-address.txt | jq -r .address) -p 5055 -d /home/brandon/workspace/blockchain/tools/smart-contract-watch
