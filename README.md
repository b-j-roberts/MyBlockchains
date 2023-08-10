# My Blockchains

Various blockchain / cryptocurrency projects used for utility, learning, and testing blockchain systems.

### [Eth Private Network](https://github.com/b-j-roberts/MyBlockchains/tree/master/eth-private-network)

Scripts, Docker setup, and thing to launch an ethereum private network using the Clique POA consensus protocol.
Eth network should by default theoretically run ~120x faster than eth mainnet ( 10x from gas limit increase & 12x from blocktime decrease )

**Some useful features:**
- Large increase in transaction speeds, compared to mainnet ( & its all for you! )
- Launch a miner ( POA ) node, optionally launch peer rpc nodes and connect as peers
- Local, Docker, & Kubernetes setups + quick launch commands thru make
- Github Actions automatically tests Miner, RPC, transactions, and build + pushes docker images to [dockerhub](https://hub.docker.com/repositories/brandonjroberts).

![Grafana Dash Miner](https://github.com/b-j-roberts/MyBlockchains/blob/master/eth-private-network/media/grafana-dash-miner.png)


### Naive Blockchains

Various "naive" implementations of blockchains based on a specified concept ( eg. POW & L2 )

#### [Naive Cryptocurrency L2](https://github.com/b-j-roberts/MyBlockchains/tree/master/naive-blockchain/naive-cryptocurrency-l2)

EVM compatible L2 cryptocurrency based on the concept of using zk proofs to verify L2 transaction data and state transitions stored on the L1.

The L2 Blockchain uses geth functionality to run the execution layer thru a Sequencer node, which in turn batches transaction data, compresses them, and posts them to L1 for the Prover node to get and verify.

**Some cool features:**
- EVM compatibility ( Uses go-ethereum EVM implementation )
- Minimal Geth Fork
- Quickly launch an L2 Sequencer, RPC Node, Prover, and more
- Prometheus, Grafana, and Smart Contract Metrics Exporter to provide easy visibility into all parts.
- Local & Docker setups + quick launch commands thru make
- Github Actions automatically test Sequencer, Prover, Bridging, and builds + pushes docker images to [dockerhub](https://hub.docker.com/repositories/brandonjroberts).

![Transaction Storage](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/transaction-storage.png)
![Prover](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/prover.png)
![Token Bridge 1](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/token-bridge-1.png)

#### [Naive Cryptocurrency POW](https://github.com/b-j-roberts/MyBlockchains/tree/master/naive-blockchain/naive-cryptocurrency-pow)

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

#### [Naive Blockchain POW](https://github.com/b-j-roberts/MyBlockchains/tree/master/naive-blockchain/naive-blockchain-pow)

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
