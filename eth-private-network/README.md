# Geth Private Clique POA Network

Project containing scripts to launch private clique networks using geth.

**Cool Features**
- Dynamically generates `genesis.json`
- Launch a miner ( POA ) node, optionally launch peer rpc nodes and connect as peers
- Exposes unlocked accounts on nodes for testing
- Local, Docker, & Kubernetes setups + quick launch commands thru make
- Contains various load tests / other testing scripts ( transactions work )
- Github Actions automatically tests Miner, RPC, transactions, and build + pushes docker images to [dockerhub](https://hub.docker.com/repositories/brandonjroberts).

![Load Test Transactions](https://github.com/b-j-roberts/MyBlockchains/blob/master/eth-private-network/media/load-test-transactions.png)

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

cd eth-private-network
make build
```

The install may fail if missing dependencies : `go`, `make`, ...

**Prometheus Setup (optional)**
Geth is setup to expose certain useful metrics in prometheus. Add the following to your `scrape_configs` in `/etc/prometheus/prometheus.yml`
```
  - job_name: "geth-l1-miner"
    metrics_path: /debug/metrics/prometheus
    static_configs:
      - targets: ["localhost:6060"]
  - job_name: "geth-l1-rpc"
    metrics_path: /debug/metrics/prometheus
    static_configs:
      - targets: ["localhost:6061"]
```

Then reset or start prometheus using something like `sudo systemctl start prometheus.service`

**Grafana Setup (optional)**
This repo contains a grafana dashboard with useful views of the running Miner & RPC nodes. That is contained at `../dashboards/grafana/`. See the instructions there for more info.

![Grafana Dash Miner](https://github.com/b-j-roberts/MyBlockchains/blob/master/eth-private-network/media/grafana-dash-miner.png)

![Grafana Dash RPC](https://github.com/b-j-roberts/MyBlockchains/blob/master/eth-private-network/media/grafana-dash-rpc.png)

## How to Run

### Local

**Running**
Normally, you will just build & run the miner ( PoA Agent ). This will build the geth client, generate a new account and genesis.json, then launch a new Clique PoA network.
```
make build
make run-miner
```

Other options include running as a daemon ( background task outputting to a file ) or to run without reseting the network ( ie continue with previously created datadir / state ). Respectively run like :
```
OUTPUT_FILE=out.txt make run-miner-daemon
make run-miner-noclean
```
Use `ps aux | grep geth` to find running daemon tasks to kill them.


You can additionally run RPC nodes ( Don't participate in consensus ) and connect them as a peer to the Miner to load blocks from the Clique network. Run like :
```
make run-rpc
make connect-peers
```

RPC nodes contain similar daemon and noclean options : `run-rpc-daemon` and `run-rpc-noclean`

**Verify**
If using `make run-miner` or other commands, check output for new blocks being commited and finalized.
```
INFO [08-05|18:06:16.001] Successfully sealed new block            number=10 sealhash=0949cf..0d1556 hash=dd556f..3c7b09 elapsed=999.557ms
INFO [08-05|18:06:16.001] "🔗 block reached canonical chain"        number=3  hash=86cfc9..172005
INFO [08-05|18:06:16.001] "🔨 mined potential block"                number=10 hash=dd556f..3c7b09
INFO [08-05|18:06:16.001] Commit new sealing work                  number=11 sealhash=a36793..21f534 uncles=0 txs=0 gas=0 fees=0 elapsed="144.923µs"
INFO [08-05|18:06:16.001] Commit new sealing work                  number=11 sealhash=a36793..21f534 uncles=0 txs=0 gas=0 fees=0 elapsed="258.977µs"

```

This should indicate that things are working well, however there are some other options that could give more validation:
- Test scripts : Under the `scripts` dir there are various testing scripts which can be used to test gets & transactions. These include :
    - get-latest.sh : Get latest block from chain
    - load-test-gets.sh : Uses `ab` to load test sending getBalance's to RPC server
      - Install `ab` from `apache2-utils`, also make sure to setup `ulimit` as needed
    - send-transaction.sh : Send a single transaction to the Node
    - load-test-transactions.sh : Uses `ab` to load test sending basic transactions to the node

- Prometheus Metrics and/or Grafana Dashboards: You can check the value of internals after doing various tests using PromQL & grafana dashboards. Things to check include :
    - `chain_head_block` : Latest head block number
    - `txpool_valid` : # of valid txs
    - `p2p_peers` : # of peers

**Cleanup**
The main artifacts from running locally include the data-dirs that store the chain state, and build files. Use these commands to clean up.
```
rm -rf ~/l1-miner-data/
rm -rf ~/l1-rpc-data/
make clean
```

Running geth nodes also creates an account at `~/.eth-accounts`, which will be auto removed and replaced on clean restart ( like the data dir ). But this is where accounts are read from, so use it!

### Docker

**Building**
There is no need to build & push the docker image for a typical user, but if you are doing development, these make commands might be useful for building the used docker images.
```
make docker-build
make docker-push
```
When pushing your branch changes to Github, Github Actions will build and push a docker image to dockerhub. The tag will be something like `branchname-commitsha`. When merging into master, it will update the `latest` tag image in dockerhub.
https://hub.docker.com/repositories/brandonjroberts
- https://hub.docker.com/repository/docker/brandonjroberts/eth-private-miner/general
- https://hub.docker.com/repository/docker/brandonjroberts/eth-private-rpc/general
- https://hub.docker.com/repository/docker/brandonjroberts/eth-private-node-setup/general

**Running**
You can also use docker images pushed to dockerhub ( or built locally ) to run the miner & rpc nodes :
```
make docker-run-miner
```
The above command will get the images from docker hub and run them to generate a datadir and new account, which will be used in the container as volumes. Then runs the miner POA agent like normal. If you want to run without resetting the state and accounts use `make docker-run-miner-noclean`.


To run the RPC node in docker & connect to the miner node use these commands.
```
make docker-run-rpc
make connect-peers
```

Lastly, you can noclean the rpc node aswell: `make docker-run-rpc-noclean`

**Verify**
Since docker logs to the console just like running locally, and since the network is set to host on the docker run commands, you can use the same validation steps as with a local setup.

**Cleanup**
Similarly, since volumes are used for the datadir and accounts, you can remove the docker volume datadirs to clear the state and cleanup artifacts.
```
rm -rf ~/docker-l1-miner-data/
rm -rf ~/docker-l1-rpc-data/
```

### Kubernetes

**Running**
Lastly, there is an existing kubernetes statefulset setup using the above docker images to deploy nodes into your kube cluster. This will deploy a statefulset, which will manage & setup your node pod, and deploy a service to expose ports for peers and rpc calls.
```
make kube-deploy-miner
```

For RPC nodes... ( Must wait for miner so the genesis.json file can be taken from there )
```
make kube-deploy-rpc
make kube-connect-peers
```

**Verify**
I typically use `k9s` to inspect running containers / pods in kubernetes. This allows a similar set of visual verification tools. Otherwise, `kubectl logs` should work well.

Using `make kube-connect-peers` also port-forwards the metrics and rpc ports, so you can also use other scripts and prometheus setup locally to view data in grafana...

**Cleanup**
The main options in kubernetes are to clean or reset.

Reset clears pods / state, but sts remains
```
# Reset the pods but keep the state
make kube-reset

# Reset pods and PVs to clear state
make kube-reset-all
```

Clear removes sts and pods ...
```
# Remove sts and deploys but keeps the state
make kube-clean

# Removes sts, deploys, and PVs to clear state completely
make kube-clean-all
```

---

## Other Info

### Useful links & Resources

- https://medium.com/swlh/how-to-set-up-a-private-ethereum-blockchain-c0e74260492c
- https://geth.ethereum.org/docs/fundamentals/private-network
- https://brodan.biz/blog/how-to-run-a-private-local-ethereum-network-with-geth/
- https://eips.ethereum.org/EIPS/eip-225
