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
b-j-roberts( ~/workspace/blockchain/my-chains/eth-private-network @ master ) ☯ docker build -t l1-private -f Dockerfile .
docker run -it brandonjroberts/eth-private-miner:latest

### Running in k8s

### Command line Args

### Useful Scripts

## Code Layout

### Contract
contracts, go, remix, scripts, test, make/build?
npm install

### Cmd
Launch node

### Go-Ethereum
Fork, Minimal, Local ref in other code

### Utils

### Scripts
Starting, Account gen, ....



TODO:  --txpool.globalslots





Full clean:
make clean
./scripts/clean-all.sh
cd ../../eth-private-network/ && make clean && cd ../naive-blockchain/naive-cryptocurrency-l2/


Full Run:
./scripts/start-l1-private-network.sh
make all
make deploy-private-l1
./scripts/start-sequencer-overwrite.sh

./scripts/start-prover.sh

./scripts/start-rpc-overwrite.sh

Watch UI / state :
./scripts/smart-contract-watch.sh
cd ~/workspace/blockchain/tools/explorer && npm start

send txs
cat ~/naive-sequencer-data/genesis.json
geth attach ~/naive-sequencer-data/naive-sequencer.ipc
eth.sendTransaction({from: "0x1afed87524e19ccae70f34517c328bac5f636e41", to: "06eba974246f46d6b8421fd5e0b1b5cafbeb0710", value: 100000})




make launch-miner-local
make clean
make all
make deploy-private-l1
# COPY ADDR
L1_CONTRACT_ADDRESS=<addr> make run-sequencer
make watch-smart-contract
L1_CONTRACT_ADDRESS=<addr> make run-prover


make docker-build
make docker-run-miner
make docker-build
make contracts
make deploy-private-l1
# COPY ADDR
L1_CONTRACT_ADDRESS=<addr> make docker-run-sequencer
make watch-smart-contract
L1_CONTRACT_ADDRESS=<addr> make docker-run-prover
