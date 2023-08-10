# Naive L2 Cryptocurrency

This project serves as a learning exercise, practice, and a simple example of an EVM compatible L2 cryptocurrency based on the concept of using zk proofs to verify L2 transaction data and state transitions stored on the L1.

**Cool Features**
- Minimal Geth Fork ( Used to provide EVM compatibility & Useful cryptocurrency primitives )
- Eth, ERC20, and ERC721 Bridge ( with 2 ERC20 tokens & 2 ERC721 tokens used to test )
- Quickly launch an L2 Sequencer, RPC Node, Prover, and more
- Prometheus, Grafana, and Smart Contract Metrics Exporter to provide easy visibility into all parts.
- Dynamically generates L2 network & genesis
- Exposes unlocked accounts on nodes for testing
- Local & Docker setups + quick launch commands thru make
- Github Actions automatically test Sequencer, Prover, Bridging, and builds + pushes docker images to [dockerhub](https://hub.docker.com/repositories/brandonjroberts).

The L2 Blockchain uses geth functionality to run the execution layer thru a Sequencer node, which in turn batches transaction data, compresses them, and posts them to L1 for the Prover node to get and verify.
The L2 network uses a modified Clique POA consensus, where the sequencer node(s) act as the authority agents. The primary modification to Clique is baking bridging into the consensus layer.

![Transaction Storage](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/transaction-storage.png)

---

## Setup & Install

**Install & Dependencies**
Currently, you must build from source ( or use docker images from dockerhub -- see below )
```
git clone git@github.com:b-j-roberts/MyBlockchains.git
cd MyBlockchains

git clean -xfdf
git submodule foreach --recursive git clean -xfdf
git reset --hard
git submodule foreach --recursive git reset --hard
git submodule update --init --recursive

cd naive-blockchain/naive-cryptocurrency-l2/contracts
npm install
cd ../
make
```

The install may fail if missing dependencies : `go`, `make`, `node`, `npm`, `solc`, `abigen`, `jq`, ...

- *solc* - `npm install -g solc@0.8.18`
- *abigen* - `cd go-ethereum && make devtools`

**Prometheus Setup (optional)**
Geth is setup to expose certain useful metrics in prometheus, and additional metrics can be obtained from the `smart-contract-metrics` metric exporter. Add the following to your `scrape_configs` in `/etc/prometheus/prometheus.yml`
```
  - job_name: "geth-l2-sequencer"
    metrics_path: /debug/metrics/prometheus
    static_configs:
      - targets: ["localhost:6160"]
  - job_name: "geth-l2-rpc"
    metrics_path: /debug/metrics/prometheus
    static_configs:
      - targets: ["localhost:6161"]
  - job_name: "l2-smart-contract-exporter"
    metrics_path: /metrics
    static_configs:
      - targets: ["localhost:6169"]
  - job_name: "l2-prover"
    metrics_path: /metrics
    static_configs:
      - targets: ["localhost:6162"]
```

Then reset or start prometheus using something like `sudo systemctl start prometheus.service`

**Grafana Setup (optional)**
This repo contains a grafana dashboard with useful views of the running Sequencer & Prover nodes and smart contract values. That is contained at `../dashboards/grafana/`. See the instructions there for more info.

![Transaction Storage](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/transaction-storage.png)
![Prover](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/prover.png)
![Eth Bridge](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/eth-bridge.png)
![Token Bridge 1](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/token-bridge-1.png)
![Token Bridge 2](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/token-bridge-2.png)

# How to Run

## Local

### Running

**Run L1 Node :**

Firstly, you will need to either launch a private L1 Network or gain access to an L1 RPC endpoint. This will be used to provide data availability, interfaces for Provers, Bridging contracts and native currency, and provide finality.

I recommend using the private network setup in this repo under `MyBlockchain/eth-private-network`, simply follow the instructions under the [README.md there](https://github.com/b-j-roberts/MyBlockchains/tree/master/eth-private-network). However, if you wish to use your own private network provider or RPC node, that is fine.

NOTE: Most of the below interfaces expect the L1 Network to be at `http://localhost:8545`, so some configurations may need to be changed for you.


**Deploy Contracts :**

Next is to deploy the smart contracts used for the L2 network onto the L1 chain. These include :
- `TransactionStorage` - Stores batched transaction data posted from L2 Sequencer onto L1. Also provides interfaces to Provers to "prove" batch validity.
- `L1Bridge` - L1 Eth Bridge contract used to lock & unlock bridged eth. Sends events to be watched for on L2.
- `L1TokenBridge` - L1 ERC20 & ERC721 Token Bridge contract used to lock & unlock bridged tokens. Sends events to be watched for on L2.
- Tokens - 4 ERC token contracts used for testing bridge. `BasicERC20`, `StableERC20`, `BasicERC721`, and `SpecialERC721`.

Use the following to deploy the contracts to L1:
```
make contracts
make deploy-l1-contracts
```


**Setup Configurations :**

The default configuration almost certainly wont work on your system. You will probably need to change the following configs : `sequencer.config.json`, `bridge.config.json`, and `rpc.config.json`, but feel free to look at and change other configs located at `./configs`.

The fields `data-dir`, `genesis`, and `contracts` almost certainly need to be changed to match your FS and configs like `l1Url` or `metrics` may need to changed based on your setup.


**Run Sequencer :**

Now's when things get more interesting. The sequencer node will be the main authority agent for the L2 network, and will build blocks, batch transactions, and run the bridge. To startup the L2 network and Sequencer, use the following commands:

```
make sequencer
make run-sequencer
```

Other options include running as a daemon ( background task outputting to a file ) or to run without reseting the network ( ie continue with previously created datadir / state). Respectively run like :
```
make run-sequencer-daemon
make run-sequencer-noclean
```


**Run L2 RPC Node ( optional ) :**

You can additionally run RPC nodes ( Don't participate in consensus ) and connect them as a peer to the Sequencer Node to load blocks from the L2 network. Run like :
```
make rpc
make run-rpc
```

Using the `run-rpc` command will launch the rpc node and connect it as a peer to the Sequencer using a script. ( RPC node always runs as a daemon due to this, so logs can be watched, however there is a noclean option `make run-rpc-noclean` )


**Deploy L2 Contracts :**

There are a few contracts that must be deployed on the L2 Network. These include :
- `L2Bridge` - L2 Eth Bridge contract used to send events by Sequencer to mint Eth on L2. Also allows withdrawals to burn Eth on L2 to unlock on L1. Cross chain calls taken care of by Sequencer.
- `L2TokenBridge` - L2 ERC20 & ERC721 Token Bridge contract. Similar to `L2Bridge` functionality but for Tokens.
- Tokens - 4 Cooresponding ERC token contracts used to test bridging. `BasicL2ERC20`, `StableL2ERC20`, `BasicL2ERC721`, and `SpecialL2ERC721`
    - NOTE: These L2 tokens must contain functionality to allow the Sequencer to mint / burn tokens for bridging. This is done by implementing `L2TokenMinter` interface functions.

Use the following to deploy contracts to L2:
```
make deploy-l2-contracts
```


**Run Prover :**

The last command needed to have a the full L2 network running to full capacity is to start the prover. Prover nodes will get posted batches at the L1 `TransactionStorage` contract, then submit proofs to L1 therefore finalizing L2 transactions.

WARNING: The currently generated proofs are placeholders, since building an EVM prover would be very complex. Also open source and available options all require extremely heavy computational power ( ~ 1TB Ram )

```
make run-prover
```


**Smart Contract Metric ( optional ) :**

It can be very useful to run the Smart Contract Metric Exporter job to gain visibility into the running `TransactionStorage`, Bridge, and token contracts. This job will continually query your L1 and L2 nodes contracts / RPC endpoints to keep up to date useful metrics in Prometheus. Some examples of useful metrics:
- `last_confirmed_batch` - Last batch number confirmed by prover on L1
- `bridge_balance` - Total value locked in Eth Bridge on L1
- `l2_stable_token_sequencer_balance` - Total stable token balance of sequencer on L2 ( test bridge )


**Bridge things ( optional ) :**

There are a series of make commands which make bridging easier aswell.
```
make bridge-eth-to-l2
make bridge-eth-to-l1
make bridge-basic-erc20-to-l2
make bridge-basic-erc20-to-l1
make bridge-stable-erc20-to-l2
make bridge-stable-erc20-to-l1
make bridge-basic-erc721-to-l2
make bridge-basic-erc721-to-l1
make bridge-special-erc721-to-l2
make bridge-special-erc721-to-l1
```

**Quick Launch**

To run all of the above steps with a single command to launch L1, L2, Prover, Deploy Contracts, Start smart contract metrics, and bridge a bunch of tokens and eth, you can use the `quick-launch` option. To do so use: `make quick-launch`.

You can also use `quick-launch-full` to launch RPC nodes aswell. Also use `quick-clean` to kill things launched by `quick-launch`. Lastly, after waiting some time for bridge transactions to settle, you can use `make run-metrics-test` to see if all metrics match expected values after quick-launch


### Verify
After each of the above run steps, it is best to verify the piece is working properly before moving on to avoid weird issues debugging / getting things working. This section will detail how to do that.

**Verify Outputs**

Verify L1 & L2 mining new blocks by checking `run-miner` and `run-sequencer` ( and run-rpc ) outputs for lines like this to verify they are creating new blocks :
```
INFO [08-05|18:06:16.001] Successfully sealed new block            number=10 sealhash=0949cf..0d1556 hash=dd556f..3c7b09 elapsed=999.557ms
INFO [08-05|18:06:16.001] "🔗 block reached canonical chain"        number=3  hash=86cfc9..172005
INFO [08-05|18:06:16.001] "🔨 mined potential block"                number=10 hash=dd556f..3c7b09
INFO [08-05|18:06:16.001] Commit new sealing work                  number=11 sealhash=a36793..21f534 uncles=0 txs=0 gas=0 fees=0 elapsed="144.923µs"
INFO [08-05|18:06:16.001] Commit new sealing work                  number=11 sealhash=a36793..21f534 uncles=0 txs=0 gas=0 fees=0 elapsed="258.977µs"
```

Verify Contracts Deploy Successfully by watching `miner` and `sequencer` outputs for lines like this while deploying :
```
INFO [08-09|16:19:47.519] Submitted contract creation              hash=0x1ddb07edf121ab4110083981dc5a95134cf3c82db789f5f4e127974a644b97dd from=0xc5d1F00d68e2c66e8B1b6137939B0B9ACf7C2Cc6 nonce=0 contract=0xdaac2C9AdF662E01114Ad2d26Fb1791ad3215b6D value=0
```
Then check the `deploy-l*-contracts` command gives lines like `Deployed SpecialERC721 contract to :  0xF273880756351bbA793B87053Fc3a10eC0c54440`

Verify Prover verifying proofs on L1 by looking for output lines like `Proof verified on L1!`

Lastly, Verify bridge commands by either checking tx receipts with `geth attach` and `eth.getTransactionReceipt` or check balances with `eth.getBalance` for Eth bridge & use the `./scripts/get-erc-balance.sh` script to check erc balances.


**Prometheus / Grafana Option**

The most visual option for verifying stuff is working as expected is to use the grafana dashboard or prometheus metrics to check the state of the system. This also allows insight on the state of the smart contracts that you can't get from logs.

You'll want to check for green/red panels in each section to note that the piece is actively collecting metrics, and check the other panel values respectively. Here is an example of working metrics :

![Prover](https://github.com/b-j-roberts/MyBlockchains/blob/master/naive-blockchain/naive-cryptocurrency-l2/media/prover.png)


**metrics-test**

The quickest, but most limited, option for testing things are working as expected is to use the `metrics-test` command with `make metrics-test`. This will give a list of "Pass" / "Fails" based on prometheus metrics collected after running `quick-launch`. After waiting a couple minutes after `quick-launch`, check metric values are correct with the metrics-test.


### Cleanup

The main artifacts from running locally include the data-dirs that store the chain state, and build files. Use these commands to clean up after running.
```
rm -rf ~/l1-miner-data/
rm -rf ~/l1-rpc-data/
rm -rf ~/naive-sequencer-data/
rm -rf ~/naive-rpc-data/
make clean
make quick-clean
```

NOTE: If you adjusted datadir locations in your `configs/*.config.json` files, make sure to remove to correct datadir.

Running geth nodes also creates an account at `~/.eth-accounts` and transactor at `~/.transactor`, which will be auto removed and replaced on clean restart ( like the data dirs ). But this is where are accounts are read from, so use it!


## Docker

### Building

There is no need to build & push the docker image for a typical user, but if you are doing development, these make commands might be useful for building the used docker images.
```
make docker-build
make docker-push
```
When pushing your branch changes to Github, Github Actions will build and push a docker image to dockerhub. The tag will be something like `branchname-commitsha`. When merging into master, it will update the `latest` tag image in dockerhub.
https://hub.docker.com/repositories/brandonjroberts
- https://hub.docker.com/repository/docker/brandonjroberts/naive-l2-sequencer/general
- https://hub.docker.com/repository/docker/brandonjroberts/naive-l2-prover/general
- https://hub.docker.com/repository/docker/brandonjroberts/naive-l2-smart-contract-exporter/general

### Running

You can also use docker images pushed to dockerhub ( or built locally ) to run the sequencer, prover, and smart contract metrics :
```
make docker-run-sequencer
make docker-run-prover
make docker-run-smart-contract-metrics
```
The above commands will get the images from docker hub and run them with `docker run` and some predefined setups. You will need to edit `configs/docker-*.config.json` configs to setup your docker environment how desired.

The steps for running the overall full setup, from L1 private node to Bridging to L2, it is mostly the same as running locally, but you would add `docker-` to the make commands. So for example, a quick full runthrough would be :
```
# In eth private network
make docker-run-miner
make docker-run-rpc
make connect-peers

# In Naive L2
make docker-deploy-l1-contracts
make docker-run-sequencer
make docker-deploy-l2-contracts
make docker-run-prover
make docker-run-smart-contract-metrics
make docker-bridge-eth-to-l2
```

The docker setup uses volumes to load in accounts, contracts, and the datadir. So containers can be noclean restart aswell and use volumes to store data between runs.

### Verify

Since docker logs to the console just like running locally, and since the network is set to host on the docker run commands, you can use the same validation steps as with a local setup.

### Cleanup

Similarly, since volumes are used for the datadir and accounts, you can remove the docker volume datadirs to clear the state and cleanup artifacts.
```
rm -rf ~/docker-l1-miner-data/
rm -rf ~/docker-l1-rpc-data/
rm -rf ~/docker-sequencer/
rm -rf ~/tmp/l1-docker-data*
```

## Useful Scripts

- `./scripts/quick-launch.sh` - Completely launch everything from L1 to Bridging Eth to L2 locally. ( You will need to edit `configs/` for this to work )
- `./scripts/quick-clean.sh` - Clean everything up from `quick-launch` & stop running nodes
- `./build/metrics-test` - Test prometheus metric values after `quick-launch` to ensure things are working
- `./build/bridge` - Bridge Eth or Tokens between L1 & L2
- `./scripts/get-erc-balance.sh` - Basic bash script to get erc token balance


# Code Layout

### Contracts

- NPM package setup for managing node dependencies for deploying, testing, and openzepplin.
- Solidity smart contracts for L1 & L2, as well as test Token contracts. ( under `contracts/`
- Using `abigen` contract compile into go module under `go/`
- `scripts/` to deploy and test contracts
- `tests/` for contract unit tests

### Cmd
Go `main` modules used to launch go commands. Contains quick scripts & full node commands.

Includes `sequencer`, `prover`, `bridge`, ...

### Go-Ethereum
Minimal geth fork containing code needed to launch a Naive L2 Node. Also used as local geth code reference for development and useful cryptocurrency primitives, such as P2P services, Mempool, RPC, State management, and much more.

### Src
Go source code and modules used throughout the repo. Contains :
- `consensus/` - L2 Clique PoA Fork for Sequencer network w/ bridging
- `core/` - Core classes for L2 mechanics `Node`, `Batcher`, `Prover`, `BridgeWatcher`, `Sequencer`, ...
- ...


---

## Future work
- bake in testnet w/ 90% block space and incentivization to fill block on mainnet
- naive stable coin fixup
- think how to make debugging easier
- grep todo
    - `grep -r TODO --exclude-dir={go-ethereum,node_modules}`
    - `cd go-ethereum && git diff ea9e62ca3db5c33aa7438ebf39c189afd53c6bf8 | grep TODO`
