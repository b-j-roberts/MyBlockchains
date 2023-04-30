#!/bin/bash
#TODO: Thru all help args

# This script launches 2 nodes & connects them to the each other
# This private network uses slot time of 1 second ( 12x faster ) & gas limit of 300 million ( 10x higher ) - So 15 tx/sec * 12 * 10 = 1500 tx/sec
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
SCRIPT_DIR=$SCRIPT_DIR/..
mkdir $SCRIPT_DIR/logs/

$SCRIPT_DIR/scripts/clean-networks.sh


# Launch node 1
GETH_LOG_DIR="$SCRIPT_DIR/logs/$(date '+%Y-%m-%d_%H:%M:%S')"
mkdir -p $GETH_LOG_DIR
ACCOUNT1=$(cat $SCRIPT_DIR/data/keystore/* | jq -r '.address' | head -n 1)
nohup geth --networkid 505 --datadir $SCRIPT_DIR/data/ --nodiscover --http --http.api "eth,net,web3,personal,web3" --unlock "0x$ACCOUNT1" --mine --allow-insecure-unlock --password $SCRIPT_DIR/password.txt >> "$GETH_LOG_DIR/geth-1.log" --miner.etherbase "0x$ACCOUNT1" --miner.gaslimit 300000000 2>&1 &

echo "Waiting for node 1 to start..."
while true; do
  if grep -q "self=\"enode://" $GETH_LOG_DIR/geth-1.log; then
    break
  fi
  sleep 1
done

NODE1_ENODE=$(geth attach --exec admin.nodeInfo.enode $SCRIPT_DIR/data/geth.ipc)
echo "Node 1 started with enode $NODE1_ENODE"

# Launch node 2
nohup geth --networkid 505 --datadir $SCRIPT_DIR/data2/ --nodiscover --authrpc.port 8787 --port 30306 --http --http.api "eth,net,web3,personal,web3" --http.port 8488 --allow-insecure-unlock >> "$GETH_LOG_DIR/geth-2.log" 2>&1 &

echo "Waiting for node 2 to start..."
while true; do
  if grep -q "self=\"enode://" $GETH_LOG_DIR/geth-2.log; then
    break
  fi
  sleep 1
done

NODE2_ENODE=$(geth attach --exec admin.nodeInfo.enode $SCRIPT_DIR/data2/geth.ipc)
echo "Node 2 started with enode $NODE2_ENODE"

# Connect node 1 to node 2
geth attach --exec "admin.addPeer($NODE2_ENODE)" $SCRIPT_DIR/data/geth.ipc
geth attach --exec "admin.peers" $SCRIPT_DIR/data/geth.ipc
