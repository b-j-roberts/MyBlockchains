#!/bin/bash
#TODO: Thru all help args

# This script launches 2 nodes & connects them to the each other
# This private network uses slot time of 1 second ( 12x faster ) & gas limit of 300 million ( 10x higher ) - So 15 tx/sec * 12 * 10 = 1500 tx/sec

# Launch node 1
rm private-network/geth-1.log
nohup ./build/bin/geth --networkid 505 --datadir private-network/data/ --nodiscover --http --http.api "eth,net,web3,personal,web3" --allow-insecure-unlock --unlock 0xc2a9bd81ef8af24f3baacfc6bf611cd8d20d0114 --password "private-network/password.txt" --mine --miner.etherbase 0xc2a9bd81ef8af24f3baacfc6bf611cd8d20d0114 --miner.gaslimit 300000000 >> "private-network/geth-1.log" 2>&1 &

echo "Waiting for node 1 to start..."
while true; do
  if grep -q "self=\"enode://" private-network/geth-1.log; then
    break
  fi
  sleep 1
done

NODE1_ENODE=$(./build/bin/geth attach --exec admin.nodeInfo.enode private-network/data/geth.ipc)
echo "Node 1 started with enode $NODE_ENODE"

# Launch node 2
rm private-network/geth-2.log
nohup ./build/bin/geth --networkid 505 --datadir private-network/data2/ --nodiscover --authrpc.port 8787 --port 30306 --http --http.api "eth,net,web3,personal,web3" --http.port 8488 --allow-insecure-unlock >> "private-network/geth-2.log" 2>&1 &

echo "Waiting for node 2 to start..."
while true; do
  if grep -q "self=\"enode://" private-network/geth-2.log; then
    break
  fi
  sleep 1
done

NODE2_ENODE=$(./build/bin/geth attach --exec admin.nodeInfo.enode private-network/data2/geth.ipc)
echo "Node 2 started with enode $NODE2_ENODE"

# Connect node 1 to node 2
./build/bin/geth attach --exec "admin.addPeer($NODE2_ENODE)" private-network/data/geth.ipc
./build/bin/geth attach --exec "admin.peers" private-network/data/geth.ipc
