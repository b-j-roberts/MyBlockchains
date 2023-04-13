# Naive Go CPU based POW Cryptocurrency Client

This cryptocurrency project serves as practice and a simple naive example of a POW blockchain backing an account based cryptocurrency.

libp2p is used for peer to peer connections, which will be established based on a beginner set of defined peers. Initial peers in the config/ directory should be setup.
Metrics are collected using prometheus for things like blockheight, peer count, and other useful stats.

## Different types of nodes

There are different types of node clients, which allow users to select varying ways to interact with the blockchain.

### Miner Node

Miner nodes are the nodes that run the POW consensus mechs using transactions coming onto the mempool. These nodes will require the most load on the system, and will use their cpu to mine blocks and connnect to other Miner & RPC nodes to transmit mined blocks and get rewards.

### RPC Node

RPC nodes are the nodes that are used to accept RPC requests to create Transactions and add them to the mempool or to inspect data on the blockchain. Essentially, the client will be storing the blockchain locally and obtain updates using a p2p mechanism, then they will expose an RPC endpoint to query that data or create transactions to post in the mempool.

### Transaction Node

This is a lightweight node which exposes a p2p connection only into the mempool network breifly, in this time it will be used to create and send a transaction. Once getting confirmed by a few nodes, it should shut down its connection to the mempool.

TODO: Future idea is to create a full RPC node which doesn't need to store data and queries other nodes using an incentive ( fee ).

## Building & Running

TODO: Download dependencies ( ex: btcutil, libp2p, ... ) (go get &/or install)
Go version 1.19

1. Ensure you have an account setup. If you do not follow these steps.
```
make account <account-name>
```
To see account info, like pub-private key pair and account address check accounts/account-<account-id>/

2. Run the client
### Running locally w/ go
1. Ensure you setup peers to trusted peers in the network at config/peer-list.txt
2. Build and run client
```
# <node-type> can be one of miner, rpc, or transaction
make <node-type>-node
./bin/<node-type>-node
```
3. Test it is working :
```
# Watch logs for errors and status of startup
# After node starts accepting RPC requests
```

### Running w/ docker
1. Ensure you setup peers to trusted peers in the network at config/peer-list.txt
2. Docker run
```
docker build --build-arg NODE_TYPE=<node-type> -t naive-<node-type> .
docker run -it naive-<node-type>
```
3. Test it is working :
```
# Watch logs for errors and status of startup
# After node starts accepting RPC requests
```

### Running in k8s
1. Ensure you setup peers to trusted peers in the network at config/peer-list.txt
2. Ensure you build and push the docker image
3. kubectl apply
```
kubectl apply -f kube/<node-type>-node-deploy-sts.yaml
```
3. Test it is working :
```
# Watch logs for errors and status of startup
# After node starts accepting RPC requests
```

### Command line arguments
Testing:
Check if all transactions go where they can with all possible verification steps
Check if all blocks go where they can with all possible verification steps
TODO: transaction node only reads p2p data on transaction interested in
Check if node fails based on different kinds of invalid data ( like nonce, transactions, blocks, ... )

mem-peer & mem-port : Start 2 transaction nodes, connect them p2p, send a transaction on one, wait to see it show on all: in all envs (-- local --, -- docker --,-- k8s --) & Other node types ( miner --- & rpc --- ) & cross all node types
```
# T1
## Clear out config/peer-list-mempool.txt
./bin/transaction-node

# T2
## Add above peer to config/peer-list-mempool.txt
./bin/transaction-node --rpc 18987 --port 8988

# T3
./src/scripts/rpc-transaction.sh -a accounts/account-121212/ -t 1 -am 1000 -f 100 -n 1 -p 18987 -o localhost
./src/scripts/rpc-transaction.sh -a accounts/account-121212/ -t 1 -am 1000 -f 100 -n 2 -p 8987 -o localhost



# T1
docker run -it brandonjroberts/naive-transaction-node

# T2
docker run -e PEER_TO_CALL="/ip4/172.17.0.2/tcp/8985/p2p/QmSzpK23zd5TDpnks3pn4oXwtUMXAZtGEifE6WSc74pWnS" -it brandonjroberts/naive-transaction-node

# T3
./src/scripts/rpc-transaction.sh -a accounts/account-121212/ -t 1 -am 1000 -f 100 -n 1 -p 8987 -o 172.17.0.2
./src/scripts/rpc-transaction.sh -a accounts/account-121212/ -t 1 -am 1000 -f 100 -n 2 -p 8987 -o 172.17.0.3


# T1
kubectl apply -f kube/transaction-node-deploy-sts.yaml
# k9s get logs for ip / id & add value to kube/transaction-node-deploy-sts-with-peer.yaml at PEER_TO_CALL env
kubectl apply -f kube/transaction-node-deploy-sts-with-peer.yaml

TODO: account not on node, send less info? - create transaction from go code and then send data
kubectl port-forward service/naive-cryptocurrency-pow-transaction-app-service 8987:8987
./src/scripts/rpc-transaction.sh -a accounts/account-121212/ -t 1 -am 1000 -f 100 -n 1 -p 8987 -o localhost
kubectl port-forward service/naive-cryptocurrency-pow-transaction-peer-app-service 8987:8987
./src/scripts/rpc-transaction.sh -a accounts/account-121212/ -t 1 -am 1000 -f 100 -n 2 -p 8987 -o localhost
```
git
multipeer mempool
clear mempool transactions when done / mined
                      Do the same w/ 3, 4, 5 transaction nodes in all envs & other node types
blk-peer & blk-port : Do the same as above but w/ miner nodes
Test the same as above connecting mulitple node types

account-id : create different accounts and share w/ volume in k8s. Ensure same # doesn't give same account

rpc : Try all rpc requests & change to rpc from http?
secio : figure it out
seed : Figure it out
air : test airdrop in all envs and make sure my account is saved somewhere private
snap : test snapshot in all envs and save volume works

TODO: Test incorrect inputs/try duping system/...

** Generic Arguments **
- mem-peer   : (Required) Peer to dial for mempool data #TODO: Change to config file and use same peer for blk & mem?
- blk-peer   : (Required) Peer to dial for blockchain / ledger data
- account-id : (Required) Identifier for account
- mem-port   : (Default=8985) Port to listen for mempool p2p connections
- blk-port   : (Default=8986) Port to listen for blockchain / ledger p2p connections
- rpc        : (Default=8987) Port to listen for RPC calls & Prom metrics
- secio      : (Default=false) Enable secio
- seed       : (Default=random) Seed for id generation
- air        : (Default=./config/airdrop_config.csv) Config used for airdrop #TODO: Optional and default
- snap       : (Optional) blockchain from snapshot file given
- accounts   : (Default=./accounts/) Account directory containing account data

** Transaction Node **
- peer       : (Required) Peer to dial for mempool data #TODO: Change to config file and use same peer for blk & mem?
- port       : (Default=8985) Port to listen for mempool p2p connections
- seed       : (Default=random) Seed for id generation
- rpc        : (Default=8987) Port to listen for RPC calls & Prom metrics
- secio      : (Default=false) Enable secio


### Debugging

Here are instructions to debug locally :



### Scripts

- generateAccount.go : Go script using the signer code to create an account based on the accound-id given. To run use something like this:
```
go run scripts/generateAccount.go --account-id 12345
```

- rpc-transaction.sh : script to send a transaction to the rpc endpoint. To run use something like this:
```
./rpc-transaction.sh -a <account> --to <to-address> -am <amount> -f <fee> -n <transaction-nonce>
```

- startup-miner-node.sh : Runs a miner node from miner-node.go with default configuration.
```
./startup-miner-node.sh
```

- startup-rpc-node.sh : Runs a rpc node from rpc-node.go with default configuration.
```
./startup-rpc-node.sh

```
- startup-transaction-node.sh : Runs a transaction node from transaction-node.go with default configuration.
```
./startup-transaction-node.sh
```

## Code Layout

All client code can be found in the `src/` directory.

### Ledger

This is the most important directory, containing all the code for the ledger, blocks, transactions, and accounts.

The object `ledger.TheLedger` stores the current clients overall ledger state. The Ledger state contains the following:
```
AccountBalances map[uint64]uint64
AccountNonces map[uint64]uint
Blockchain []Block
```

In order to store balances, nonces to prevent transaction double uses, and all the blocks from the blockchain.


The ledger can be created from block 1 using hardcoded genesis block data and an airdrop config. Or the ledger can be loaded from a snapshot.
Either way, the longest chain based POW consensus mechanism will take over, so if there are other peers with longer chain data, then that data will start to be
loaded as part of TheLedger.

Since airdrop and peers can be configured, a new network can be created by selecting certain peers. #TODO: Check this statement
Initial airdrop is setup in config/airdrop_config.csv based on the following format per line : <account-address>,<%-of-total-supply>
Note : Sum of all %s must add to 1 (100%). Otherwise there will be an error on startup. Also, airdrop config wont be used if loading chain state from a snapshot.


### Block

Each block is created in 2 steps. First, the `CreateUnminedBlock` function is called, containing a list of transactions, the previous block, and the new blockheight of the block.
This then generates a block with a null nonce, which is to be mined as part of the `miner.MineBlockNonce` method, which we will discuss later on.

The block is then added to `TheLedger` and the POW consensus mechanism takes place using the p2p connections after that.

A block contains this info :
```
type BlockHeader struct {
  Version int
  PrevBlockHash uint32
  TransactionHash uint32
  Timestamp int64
  Difficulty uint32
  Nonce int
}

type Block struct {
  BlockSize uint
  BlockHeader BlockHeader
  Transactions []Transaction
}
```

& is based on notes from this document : https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html

### Difficulty

The difficulty is determinalistically computed using an algorithm which changes the difficulty dynamically to make a block take about 30 seconds to mine.
So more validators and POW power causes difficulty to go up, since blocks will be mined faster and faster.

The algorithm takes toe blocktime over the last 100 blocks and scales the difficulty based on how long that took vs expected

### Transaction

Transactions can only store transfer data. There is no underlying VM. So a transaction is laid out like :
```
FromAddress uint64
ToAddress uint64
Amount uint64
Fee uint64
Nonce uint
Signature []byte
SignerPublicKey rsa.PublicKey
```

Which will update the balances of the addresses accordingly, give the transactions miner the allocated fee, change the account nonce to given value for nonce checks, and verify the tx with the signature.

### Mempool

The mempool module carries the object `mempool.TheMempool` which contains all the transactions to be processed into blocks. This is done using a p2p connection with other nodes.

### Miner

A miner node can be run using the `miner-node.go` code, and is used to run `CreateUnminedBlock` and `MineBlockNonce` functions to participate in consensus and block creation.
The miner will be awarded fees from the block's transactions, and will be incentivized to promote that new block to peers on the network & act kindly as to not waste time / energy.

### Metrics

Contains code to setup Prometheus metrics to be exposed at the /metrics endpoint. Existing prometheus metrics are the following :

- requests_total : # of requests to the node rpc
- request_duration_seconds : Duration for sending reqests ( total )
- block_height : Blockchain (TheLedger) tip height
- blocks_mined : # of blocks mined by miner node
- snapshot_loading : Status indicator for loading snapshot ( 0 = No snapshot, 1 = loading, 2 = verifying, 3 = done )
- inbound_peers : # of inbound peers

### P2P

Contains libp2p code for p2p connections on the ledger / blockchain & the mempool. The main functions are the `makeBasicHost` which sets up a libp2p identity & host with a rw buffer.
These buffers are continuously checked and used to read block information and transmit it to peers in the `readData` and `writeData` functions.

After receiving blocks in the `writeData` function it must be verified. Also, `writeData` contains code ( for now ) which can be used to input data requests or transactions directly into the client node using stdin.

### RPC

Contains `RpcSetup` functions for RPC and Transaction nodes individually. Each setup function contains a set of `http.HandleFunc` function calls each setting up an http endpoint with
a certain functionality. Here are the existing rpc methods :

- /help : Sends RPC help information for possibile rpc requests
- /block : Return information about the block indexs
- /add : send a transaction to the mempool

NOTE : For now http not rpc & /metrics is exposed here for prom

### Snapshot

TODO: Contains code to load the ledger from a snapshot

### Signer

This portion of the code contains things related to signing transactions and linking addresses to public keys. The Address <-> PublicKey interfacing can be found in signer/address.go and the code for generating public-private key pairs and loading them can be found in keygen.go
