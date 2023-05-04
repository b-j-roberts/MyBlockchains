# My Blockchains

Various blockchain / cryptocurrency projects used for utility, learning, and testing blockchain systems.

### Eth Private Network

Scripts & Things to launch an ethereum private network on the Clique POA consensus protocol.
Eth network should theoretically run ~120x faster than eth mainnet ( 10x from gas limit increase & 12x from blocktime decrease )

**Some useful features:**
- Creates accounts for POA "Miner" node & RPC node
- Dynamically generates `genesis.json` file with POA "Miner" account as consensus agent and configures blocktime
- Unlocks accounts for easy development
- Large increase in transaction speeds, compared to mainnet ( & its all for you! )
- Various basic testing + load testing scripts
- Launches Miner & RPC Node + connects them as peers

### Naive Blockchains

Various "naive" implementations of blockchains based on a specified concept ( eg. POW & L2 )

#### Naive Blockchain POW

Go based L1 blockchain implementation based on the same POW consensus as Bitcoin.
A block only contains a single integer of data, representing the transaction data, and the block header, containing the hash chain, difficulty, nonce, ...

Note: There are no accounts or smart contracts, this is not a cryptocurrency only a test blockchain.

Also, sorry the code here is messy. This project mainly served as a learning playground for libp2p, POW consensus, & blockchain development.

**Some cool features:**
- libp2p based p2p communication of blocks ( longest POW chain from genesis is considered canonical chain )
- Snapshot storage for quick blockchain reload
- Fairly barebones example and overview of POW consensus
- Prometheus metrics for tracking blockchain & client metrics
- Terminal interactive UI & http interface for sending block data & ... ( labeled as RPC in code incorrectly )

#### Naive Cryptocurrency POW

Go based L1 cryptocurrency implementation with an account database model like Ethereum & POW conesnsus like Bitcoin ( So not UTXO ).
Blocks store transactions similar to those in the Bitcoin network. Smart contracts or EVM is not supported.

Warning: This project is a "naive" implementation and does have security vulnerabilities.

**Some cool features.**
- 3 different node types : Miner, RPC, & Transaction ( Only connects to mempool to send transactions )
- 30 second blocktime determined by dynamic deterministic difficulty
- Transaction signing with private key
- Easy to understand account model where each account just contains an Address, Balance & Nonce.
- libp2p based p2p communication for blocks & mempool txs ( longest POW chain from genesis is considered canonical chain )
- Snapshot storage & Loading
- Docker setup to launch nodes in a container. Also kube statefulset to deploy in k8s ( exposes ports to connect peers )
- Scripts to generate accounts, send http transactions, and startup nodes


#### Naive Cryptocurrency L2

Go based L2 cryptocurrency implementation with L1 contracts for base layer settlement & data availability.

**Some cool features:**
- EVM compatibility ( Uses go-ethereum EVM implementation )
- Takes advantage of geth implementations for many things : txPool, RPC setup, Blockchain, Database, ...
- L2 chain is Clique based POA chain where the sequencer node is the key authority.
- Batcher relays L2 transactions to L1 for data availability & block root storage
- Minimalized geth fork ( Geth fork for L2 chain, minimally changed to prevent difficulty upgrading )
- L1 contract development setup w/ build, deploy & test scripts.
- Scripts to dynamically generate L2 genesis, ...

Future :
- Prover proves state root transitions & finalizes blocks
- 
