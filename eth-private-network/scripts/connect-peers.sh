#!/bin/bash
# 
# This script is used to connect a node to a peer

ENODE_IP="localhost"
ENODE_PORT="30306"

display_help() {
  echo "Usage: connect-peers.sh [OPTION]..."
  echo
  echo "    -h, --help        Display this help message."
  echo "    -e, --enode       Enode of the peer to connect to. (Required or use enode-ipc)"
  echo "    -P, --enode-host  Path to the peer's http rpc endpoint. (Required or use enode)"
  echo "    -H, --host        Path to the host's http rpc endpint. (Required)"
  echo "    -i, --ip          IP address of the peer to connect to. (Default: localhost)"
  echo "    -p, --port        Port of the peer to connect to. (Default: 30306)"
  echo
  echo "Example: ./scripts/connect-peers.sh -P ~/l1-rpc-data/geth.ipc -H ~/l1-miner-data/geth.ipc"
  exit 1
}

while getopts ":he:i:p:H:P:" opt; do
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
    H | host-ipc)
      HOST=$OPTARG
      ;;
    P | enode-ipc)
      ENODE_HOST=$OPTARG
      ENODE=$(geth attach --exec "admin.nodeInfo.enode" $ENODE_HOST | tr -d '"' | cut -d '@' -f 1 | cut -d '/' -f 3)
      echo "Enode: $ENODE"
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

if [[ -z "$ENODE" && -z "$ENODE_HOST" ]]; then
  echo "Enode or Enode IPC is required."
  display_help
fi

if [ -z "$HOST" ]; then
  echo "Host IPC is required."
  display_help
fi

echo "Connecting node at $HOST to peer: $ENODE@$ENODE_IP:$ENODE_PORT"
geth attach --exec "admin.addPeer(\"enode://${ENODE}@${ENODE_IP}:${ENODE_PORT}\")" $HOST
geth attach --exec "admin.peers" $HOST
