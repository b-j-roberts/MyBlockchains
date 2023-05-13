#!/bin/bash

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"
NAIVE_RPC_DATA="$HOME/naive-rpc-data"

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR="${BASE_DIR}/.."

# Cleanup previous runs
rm -rf ${NAIVE_RPC_DATA}
mkdir -p ${NAIVE_RPC_DATA}

ps aux | grep "naive-rpc" | grep -v grep | awk '{print $2}' | xargs kill -9

# Start the sequencer
cp ${BASE_DIR}/scripts/config/password.txt ${NAIVE_RPC_DATA}/password.txt
${BASE_DIR}/go-ethereum/build/bin/geth account new --datadir ${NAIVE_RPC_DATA} --password ${NAIVE_RPC_DATA}/password.txt
ACCOUNT1=$(cat ${NAIVE_RPC_DATA}/keystore/* | jq -r '.address' | head -n 1)

cp ${BASE_DIR}/build/genesis.json ${NAIVE_RPC_DATA}/genesis.json

go run ${BASE_DIR}/cmd/rpc/ --datadir ${NAIVE_RPC_DATA} > ${NAIVE_RPC_DATA}/rpc.log 2>&1 &

echo "Waiting for rpc to start..."
while true; do
  if grep -q "self=enode://" "${NAIVE_RPC_DATA}/rpc.log"; then
    break
  fi
  sleep 1
done

ENODE=$(geth attach --exec admin.nodeInfo.enode ${NAIVE_RPC_DATA}/naive-rpc.ipc)
echo "ENODE: ${ENODE}"

geth attach --exec "admin.addPeer(${ENODE})" ${NAIVE_SEQUENCER_DATA}/naive-sequencer.ipc
geth attach --exec "admin.peers" ${NAIVE_SEQUENCER_DATA}/naive-sequencer.ipc
