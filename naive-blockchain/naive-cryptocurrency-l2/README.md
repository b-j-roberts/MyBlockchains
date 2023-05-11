# Naive L2 Cryptocurrency

This project serves as a learning exercise, practice, and a simple example of an EVM compatible L2 cryptocurrency based on the concept of using zk proofs to verify L2 transaction data and state transitions on L1.

A minimal Geth fork is used to provide EVM compatibility and tons of useful cryptocurrency primitives, such as RPC, Blockchain, state DB, ...

The L2 Blockchain uses geth functionality to run the execution layer thru a Sequencer node. The network relies on a Clique base POA consensus, where the sequencer node(s) act as the authority agents. Other nodes can send transactions to RPC nodes or the seuqencer nodes to add them to the mempool.

## Types of Nodes

### Sequencer
Permissioned
POA Agent
Runs execution, posts batch tx to l1 contract, 

### Prover

### RPC

## Building & Running

### Contract / deploying?

### Running locally w/ go

### Running w/ Docker

### Running in k8s

### Command line Args

### Useful Scripts

## Code Layout

### Contract
contracts, go, remix, scripts, test, make/build?

### Cmd
Launch node

### Go-Ethereum
Fork, Minimal, Local ref in other code

### Utils

### Scripts
Starting, Account gen, ....





Full clean:
make clean
./scripts/clean-all.sh
cd ../../eth-private-network/ && make clean && cd ../naive-blockchain/naive-cryptocurrency-l2/


Full Run:
make all
./scripts/start-l1-private-network.sh
make deploy-private-l1
./scripts/start-sequencer-overwrite.sh

./scripts/start-prover.sh

send txs
cat ~/naive-sequencer-data/genesis.json
geth attach ~/naive-sequencer-data/naive-sequencer.ipc
eth.sendTransaction({from: "0x1afed87524e19ccae70f34517c328bac5f636e41", to: "06eba974246f46d6b8421fd5e0b1b5cafbeb0710", value: 100000})
