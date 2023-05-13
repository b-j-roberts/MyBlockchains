.PHONY: all contracts sequencer prover rpc verifier clean

all: contracts sequencer prover rpc verifier

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

deploy-private-l1:
	cd contracts && make deploy-private

clean:
	rm -rf build/
	cd contracts && make clean

spotless: clean
	./scripts/clean-all.sh
	cd ../../eth-private-network/ && make clean

# Define the default rule
.DEFAULT_GOAL := all
