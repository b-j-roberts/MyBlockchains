#!/bin/bash
# 
# This script is used to connect a node to a peer running in a docker container.

ENODE_IP="localhost"
ENODE_PORT="30303"

MINER_CONTAINER_NAME="eth-private-miner"
RPC_CONTAINER_NAME="eth-private-rpc"

display_help() {
  echo "Usage: connect-docker-peers.sh [OPTION]..."
  echo
  echo "    -h, --help        Display this help message."

  echo "    -e, --enode       Enode of the peer to connect to. (Required or use enode-ipc)"
  echo "    -P, --enode-ipc   Path to the peer's IPC file. (Required or use enode)"
  echo "    -I, --rpc-ipc     Path to the RPC's IPC file. (Required)"

  echo "    -m, --miner-cont  Name of the miner container. (Default: eth-private-miner )"
  echo "    -r, --rpc-cont    Name of the RPC container. (Default: eth-private-rpc )"
  echo "    -i, --ip          IP address of the peer to connect to. (Default: localhost)"
  echo "    -p, --port        Port of the peer to connect to. (Default: 30306)"
  echo
  echo "Example : ./scripts/connect-docker-peers.sh -P /l1-miner-data/geth.ipc -I /l1-rpc-data/geth.ipc"
  exit 1
}

while getopts ":he:i:I:m:r:p:P:" opt; do
  case $opt in
    h | help)
      display_help
      ;;
    e | enode)
      ENODE=$OPTARG
      ;;
    i | ip)
      ENODE_IP=$OPTARG
      ;;
    p | port)
      ENODE_PORT=$OPTARG
      ;;
    I | rpc-ipc)
      RPC_IPC_PATH=$OPTARG
      ;;
    r | rpc-cont)
      RPC_CONTAINER_NAME=$OPTARG
      ;;
    P | enode-ipc)
      ENODE_IPC_PATH=$OPTARG
      echo "Enode: $ENODE"
      ;;
    m | miner-cont)
      MINER_CONTAINER_NAME=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      display_help
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      display_help
      ;;
  esac
done

if [[ -z "$ENODE" && -z "$ENODE_IPC_PATH" ]]; then
  echo "Enode or Enode IPC is required."
  display_help
fi

if [ -z "$RPC_IPC_PATH" ]; then
  echo "RPC IPC is required."
  display_help
fi

if [ -z "$ENODE" ]; then
    ENODE=$(docker exec $MINER_CONTAINER_NAME geth attach --exec "admin.nodeInfo.enode" $ENODE_IPC_PATH | tr -d '"' | cut -d '@' -f 1 | cut -d '/' -f 3)
fi

echo "Connecting node at $RPC_CONTAINER_NAME to peer on container $MINER_CONTAINER_NAME at : $ENODE@$ENODE_IP:$ENODE_PORT"
docker exec $RPC_CONTAINER_NAME geth attach --exec "admin.addPeer(\"enode://${ENODE}@${ENODE_IP}:${ENODE_PORT}\")" $RPC_IPC_PATH
docker exec $RPC_CONTAINER_NAME geth attach --exec "admin.peers" $RPC_IPC_PATH
