.PHONY: all contracts sequencer prover rpc verifier allow-tokens clean smart-contract-metrics bridge metrics-test
.DEFAULT_GOAL := all

all: contracts sequencer prover rpc verifier smart-contract-metrics allow-tokens bridge

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

deploy-l1-contracts:
	rm -f contracts/builds/*-address.txt
	cd contracts && make deploy-l1-contracts

run-sequencer:
	./scripts/start-sequencer.sh -f configs/sequencer.config.json -x

run-sequencer-github-actions:
	./scripts/start-sequencer.sh -f configs/github-actions-sequencer.config.json -x -o ${OUTPUT_FILE}

run-sequencer-github-actions-retry:
	./scripts/start-sequencer-retry-github-actions.sh

deploy-l2-contracts:
	cd contracts && make deploy-l2-contracts
	cp contracts/builds/l2-bridge-address.txt ~/naive-sequencer-data/contracts/l2-bridge-address.txt
	cp contracts/builds/l2-token-bridge-address.txt ~/naive-sequencer-data/contracts/l2-token-bridge-address.txt
	cd contracts && make deploy-l2-tokens
	cp contracts/builds/*-address.txt ~/naive-sequencer-data/contracts/
	cd contracts && make allow-tokens

deploy-l2-contracts-github-actions:
	cd contracts && make deploy-l2-contracts-github-actions
	cp contracts/builds/l2-bridge-address.txt /home/runner/work/sequencer-data/contracts/l2-bridge-address.txt
	cp contracts/builds/l2-token-bridge-address.txt /home/runner/work/sequencer-data/contracts/l2-token-bridge-address.txt
	cd contracts && make deploy-l2-tokens-github-actions
	cp contracts/builds/*-address.txt /home/runner/work/sequencer-data/contracts/
	cd contracts && make allow-tokens-github-actions

run-prover:
	ACCOUNT_PASS=password ./scripts/start-prover.sh

run-prover-github-actions:
	ACCOUNT_PASS=password ./scripts/start-prover.sh -o ${OUTPUT_FILE} -c configs/github-actions-sequencer.config.json



docker-deploy-contracts-l2:
	ls contracts/builds/*-address.txt | xargs -I {} docker cp {} naive-sequencer:/app/contracts/builds/
	docker exec -it naive-sequencer sh -c "cd contracts && make docker-deploy-l2"
	docker exec -it naive-sequencer sh -c "cp contracts/builds/l2-bridge-address.txt /sequencer-data/l2-bridge-address.txt"
	docker exec -it naive-sequencer sh -c "cp contracts/builds/l2-token-bridge-address.txt /sequencer-data/l2-token-bridge-address.txt"
	docker exec -it naive-sequencer sh -c "cd contracts && make docker-deploy-l2-tokens"
	docker exec -it naive-sequencer sh -c "cd contracts && make docker-allow-tokens"

allow-deployed-tokens:
	cd contracts && make allow-tokens

clean:
	rm -rf build/
	cd contracts && make clean

docker-build:
	make clean
	docker build -f docker/Dockerfile.sequencer . -t brandonjroberts/naive-l2-sequencer:latest
	docker build -f docker/Dockerfile.prover . -t brandonjroberts/naive-l2-prover:latest
	docker build -f docker/Dockerfile.smart-contract-exporter . -t brandonjroberts/naive-l2-smart-contract-exporter:latest

docker-push:
	docker push brandonjroberts/naive-l2-sequencer:latest
	docker push brandonjroberts/naive-l2-prover:latest

local-docker-run-sequencer:
	L1_CONTRACT_ADDRESS=$(shell cat contracts/builds/l1-bridge-address.txt | jq -r .address) L1_BRIDGE_ADDRESS=$(shell cat contracts/builds/l1-bridge-address.txt | jq -r .address) L1_TOKEN_BRIDGE_ADDRESS=$(shell cat contracts/builds/l1-token-bridge-address.txt | jq -r .address) make docker-run-sequencer

docker-run-sequencer:
	rm -rf ~/tmp/l1-docker-data-for-sequencer/
	mkdir -p ~/tmp/l1-docker-data-for-sequencer/
	docker cp eth-private-miner:/root/.eth-accounts/ ~/tmp/l1-docker-data-for-sequencer/
	cp contracts/builds/*-address.txt ~/tmp/l1-docker-data-for-sequencer/
	mkdir -p ~/docker-sequencer/
	cp configs/docker-sequencer.config.json ~/docker-sequencer/sequencer.config.json
	docker run -it -p 5055:5055 -p 6160:6160 -v ~/docker-sequencer/:/sequencer-data --network host --rm --name naive-sequencer \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-sequencer/,target=/root/ \
		brandonjroberts/naive-l2-sequencer:latest
	rm -rf ~/tmp/l1-docker-data-for-sequencer/

docker-run-prover:
	rm -rf ~/tmp/l1-docker-data-for-prover/
	mkdir -p ~/tmp/l1-docker-data-for-prover/
	docker cp eth-private-miner:/root/.eth-accounts/ ~/tmp/l1-docker-data-for-prover/
	echo "Running prover w/ " ${L1_CONTRACT_ADDRESS}
	docker run -it -p 6162:6162 --network host --rm --name naive-prover -e L1_CONTRACT_ADDRESS=${L1_CONTRACT_ADDRESS} \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-prover/.eth-accounts/,target=/root/.eth-accounts/ \
		brandonjroberts/naive-l2-prover:latest
	rm -rf ~/tmp/l1-docker-data-for-prover/

local-docker-run-prover:
	L1_CONTRACT_ADDRESS=$(shell cat contracts/builds/tx-storage-address.txt | jq -r .address) make docker-run-prover

docker-run-smart-contract-metrics:
	echo "Running smart contract metrics exporter w/ " ${L1_CONTRACT_ADDRESS}
	ls contracts/builds/*-address.txt | xargs -I {} docker cp {} naive-sequencer:/app/contracts/builds/
	docker cp naive-sequencer:/app/contracts/builds/ ~/tmp/l1-docker-data-for-smart-contract-metrics/
	docker run -it -p 6169:6169 -v /home/brandon/docker-sequencer:/sequencer-data/ --network host --rm --name naive-smart-contract-metrics-exporter -e L1_CONTRACT_ADDRESS=${L1_CONTRACT_ADDRESS} -e SEQUENCER_ADDRESS=${SEQUENCER_ADDRESS} -e L2_IPC_PATH=${L2_IPC_PATH} \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-smart-contract-metrics/builds/,target=/app/contracts/builds/ \
		brandonjroberts/naive-l2-smart-contract-exporter:latest

local-docker-run-smart-contract-metrics:
	rm -rf ~/tmp/l1-docker-data-for-smart-contract-metrics/
	mkdir -p ~/tmp/l1-docker-data-for-smart-contract-metrics/
	docker cp naive-sequencer:/sequencer-data/keystore/ ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/
	L1_CONTRACT_ADDRESS=$(shell cat contracts/builds/tx-storage-address.txt | jq -r .address) SEQUENCER_ADDRESS=$(shell cat ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/* | jq -r .address) L2_IPC_PATH=/sequencer-data/naive-sequencer.ipc make docker-run-smart-contract-metrics

run-sequencer-daemon:
	./scripts/start-sequencer.sh -f configs/sequencer.config.json -x -o ${OUTPUT_FILE}

run-sequencer-noclean:
	./scripts/start-sequencer.sh -f configs/sequencer.config.json

run-rpc:
	./scripts/start-rpc.sh -f configs/rpc.config.json -x -o ~/naive-rpc-data/log.txt

run-rpc-noclean:
	./scripts/start-rpc.sh -d ~/naive-rpc-data/ -k ~/l1-miner-data/keystore/ -o ~/naive-rpc-data/rpc.log

run-prover-daemon:
	ACCOUNT_PASS=password ./scripts/start-prover.sh -o ${OUTPUT_FILE}

run-smart-contract-metrics:
	./scripts/start-smart-contract-metrics-exporter.sh -s $(shell cat ~/l1-miner-data/keystore/*	| jq -r .address)

run-smart-contract-metrics-daemon:
	./scripts/start-smart-contract-metrics-exporter.sh -s $(shell cat ~/l1-miner-data/keystore/*	| jq -r .address) -o /dev/null

watch-smart-contract:
	./scripts/smart-contract-watch.sh -a $(shell cat contracts/builds/tx-storage-address.txt | jq -r .address)

watch-l2-bridge-contract:
	./scripts/smart-contract-watch.sh -a $(shell cat contracts/builds/l2-bridge-address.txt | jq -r .address) -p 5055 -d /home/brandon/workspace/blockchain/tools/smart-contract-watch

docker-own:
	sudo chown brandon ~/docker-sequencer/*

bridge-eth-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) --value 1000000000000000 --config configs/sequencer.config.json

bridge-eth-to-l2-github-actions:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) --value 1000000000000000 --config configs/github-actions-sequencer.config.json
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) --value 1000000000000000 --config configs/github-actions-sequencer.config.json
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) --value 1000000000000000 --config configs/github-actions-sequencer.config.json

docker-bridge-eth-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/*  | jq -r .address) --bridgeAddress $(shell cat contracts/builds/l1-bridge-address.txt | jq -r .address) --rpc http://localhost:8545 --value 1000000000000000 --keystore ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/

bridge-eth-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/naive-sequencer-data/keystore/* | jq -r .address) --value 10000000000000 --config configs/sequencer.config.json --to-l1

bridge-eth-to-l1-github-actions:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat /home/runner/work/sequencer-data/keystore/* | jq -r .address) --value 10000000000000 --config configs/github-actions-sequencer.config.json
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat /home/runner/work/sequencer-data/keystore/* | jq -r .address) --value 10000000000000 --config configs/github-actions-sequencer.config.json --to-l1

docker-bridge-eth-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/* | jq -r .address) --value 10000000000000 --bridgeAddress $(shell cat ~/tmp/l1-docker-data-for-smart-contract-metrics/builds/l2-bridge-address.txt | jq -r .address) --keystore ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/ --chainId 515 --to-l1 --ipc ~/docker-sequencer/naive-sequencer.ipc
	#./scripts/bridge-eth-to-l1.sh -i ~/naive-sequencer-data/naive-sequencer.ipc -v 1000000000000 -t 0x0000000000000000000000000000000000000505

bridge-basic-erc20-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 100000 --config configs/sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc

bridge-basic-erc20-to-l2-github-actions:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 100000 --config configs/github-actions-sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 100000 --config configs/github-actions-sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 100000 --config configs/github-actions-sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc

docker-bridge-basic-erc20-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/* | jq -r .address) --value 100000 --bridgeAddress $(shell cat contracts/builds/l1-token-bridge-address.txt | jq -r .address) --keystore ~/tmp/l1-docker-data-for-smart-contract-metrics/keystore/ --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc

bridge-basic-erc20-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/naive-sequencer-data/keystore/* | jq -r .address) --value 1000 --config configs/sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc --to-l1

bridge-basic-erc20-to-l1-github-actions:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat /home/runner/work/sequencer-data/keystore/* | jq -r .address) --value 1000 --config configs/github-actions-sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc --to-l1
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat /home/runner/work/sequencer-data/keystore/* | jq -r .address) --value 1000 --config configs/github-actions-sequencer.config.json --token $(shell cat contracts/builds/basic-erc20-address.txt | jq -r .address) --is-erc --to-l1

bridge-stable-erc20-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 100000 --config configs/sequencer.config.json --token $(shell cat contracts/builds/stable-erc20-address.txt | jq -r .address) --is-erc

bridge-stable-erc20-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/naive-sequencer-data/keystore/* | jq -r .address) --value 1000 --config configs/sequencer.config.json --token $(shell cat contracts/builds/stable-erc20-address.txt | jq -r .address) --is-erc --to-l1

bridge-basic-erc721-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 0 --config configs/sequencer.config.json --token $(shell cat contracts/builds/basic-erc721-address.txt | jq -r .address) --is-erc
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 1 --config configs/sequencer.config.json --token $(shell cat contracts/builds/basic-erc721-address.txt | jq -r .address) --is-erc
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 2 --config configs/sequencer.config.json --token $(shell cat contracts/builds/basic-erc721-address.txt | jq -r .address) --is-erc

bridge-basic-erc721-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/naive-sequencer-data/keystore/* | jq -r .address) --value 0 --config configs/sequencer.config.json --token $(shell cat contracts/builds/basic-erc721-address.txt | jq -r .address) --is-erc --to-l1

bridge-special-erc721-to-l2:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/l1-miner-data/keystore/* | jq -r .address) --value 0 --config configs/sequencer.config.json --token $(shell cat contracts/builds/special-erc721-address.txt | jq -r .address) --is-erc

bridge-special-erc721-to-l1:
	ACCOUNT_PASS=password ./build/bridge --address $(shell cat ~/naive-sequencer-data/keystore/* | jq -r .address) --value 0 --config configs/sequencer.config.json --token $(shell cat contracts/builds/special-erc721-address.txt | jq -r .address) --is-erc --to-l1

quick-launch:
	./scripts/quick-launch.sh -b

quick-launch-full:
	./scripts/quick-launch.sh -b -f

quick-clean:
	./scripts/quick-clean.sh

# Define the default rule
.DEFAULT_GOAL := all
