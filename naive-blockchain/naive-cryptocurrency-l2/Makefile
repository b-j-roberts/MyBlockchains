.PHONY: all contracts sequencer prover clean

all: contracts sequencer prover

contracts:
	cd contracts && make

sequencer:
	go build -o build/sequencer ./cmd/sequencer

prover:
	go build -o build/prover ./cmd/prover

deploy-private-l1:
	cd contracts && make deploy-private

clean:
	rm -rf build/
	cd contracts && make clean

# Define the default rule
.DEFAULT_GOAL := all
