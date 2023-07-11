.PHONY: all contracts sequencer prover rpc verifier clean smart-contract-metrics

all: contracts sequencer prover rpc verifier smart-contract-metrics

contracts:
	cd contracts && make

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

deploy-private-l1:
	cd contracts && make deploy-private

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

docker-run-sequencer:
	rm -rf ~/tmp/l1-docker-data-for-sequencer/
	mkdir -p ~/tmp/l1-docker-data-for-sequencer/
	docker cp eth-private-miner:/l1-miner-data/keystore ~/tmp/l1-docker-data-for-sequencer/
	echo "Running sequencer w/ " ${L1_CONTRACT_ADDRESS}
	docker run -it -p 5055:5055 --network host --rm --name naive-sequencer -e L1_CONTRACT_ADDRESS=${L1_CONTRACT_ADDRESS} \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-sequencer/,target=/l1-data/ \
		brandonjroberts/naive-l2-sequencer:latest
	rm -rf ~/tmp/l1-docker-data-for-sequencer/

docker-run-prover:
	rm -rf ~/tmp/l1-docker-data-for-prover/
	mkdir -p ~/tmp/l1-docker-data-for-prover/
	docker cp eth-private-miner:/l1-miner-data/keystore ~/tmp/l1-docker-data-for-prover/
	echo "Running prover w/ " ${L1_CONTRACT_ADDRESS}
	docker run -it --network host --rm --name naive-prover -e L1_CONTRACT_ADDRESS=${L1_CONTRACT_ADDRESS} \
		--mount type=bind,source="$(shell echo ${HOME})"/tmp/l1-docker-data-for-prover/,target=/l1-data/ \
		brandonjroberts/naive-l2-prover:latest
	rm -rf ~/tmp/l1-docker-data-for-prover/

#TODO: Bridge address
docker-run-smart-contract-metrics:
	echo "Running smart contract metrics exporter w/ " ${L1_CONTRACT_ADDRESS}
	docker run -it --network host --rm --name naive-smart-contract-metrics-exporter -e L1_CONTRACT_ADDRESS=${L1_CONTRACT_ADDRESS} \
		brandonjroberts/naive-l2-smart-contract-exporter:latest

run-sequencer:
	./scripts/start-sequencer.sh -d ~/naive-sequencer-data/ -x -k ~/l1-miner-data/keystore/

run-prover:
	./scripts/start-prover.sh -k ~/l1-miner-data/keystore/

run-smart-contract-metrics:
	./scripts/start-smart-contract-metrics-exporter.sh

watch-smart-contract:
	./scripts/smart-contract-watch.sh -a $(shell cat contracts/builds/tx-storage-address.txt | jq -r .address)

bridge-eth-to-l2:
	go run ./scripts/bridge-eth-to-l2.go --address $(shell cat ~/l1-miner-data/keystore/*  | jq -r .address) --bridgeAddress $(shell cat contracts/builds/l1-bridge-address.txt | jq -r .address) --rpc http://localhost:8545 --value 100000 --keystore ~/l1-miner-data/keystore/

bridge-eth-to-l1:
	./scripts/bridge-eth-to-l1.sh -i ~/naive-sequencer-data/naive-sequencer.ipc -v 10000 -t 0x0000000000000000000000000000000000000505

# Define the default rule
.DEFAULT_GOAL := all
