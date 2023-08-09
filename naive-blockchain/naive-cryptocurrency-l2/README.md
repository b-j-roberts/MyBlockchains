# Naive L2 Cryptocurrency

This project serves as a learning exercise, practice, and a simple example of an EVM compatible L2 cryptocurrency based on the concept of using zk proofs to verify L2 transaction data and state transitions on L1.

A minimal Geth fork is used to provide EVM compatibility and tons of useful cryptocurrency primitives, such as RPC, Blockchain, state DB, ...

The L2 Blockchain uses geth functionality to run the execution layer thru a Sequencer node. The network uses a modified Clique POA consensus, where the sequencer node(s) act as the authority agents. The primary modification to Clique is baking bridging into the consensus layer.

Users can send transactions to RPC nodes or the seuqencer nodes to add them to the mempool.





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

Note need to setup sequencer config


## Future work
      bake in testnet w/ 90% block space and incentivization to fill block on mainnet
      naive stable coin fixup
      think how to make debugging easier
      grep todo ( commands here )

