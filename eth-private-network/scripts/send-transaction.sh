#!/bin/bash
# sudo --preserve-env bash
# ulimit -n 10000
#
# This script is used to load test the transaction throughput of a geth node.

HOST="localhost"                                           
PORT="8545"                                                
                                                           
display_help() {
  echo "Usage: send-transaction.sh [option...]"
  echo
  echo "-h, --help            Display help"
  echo "-H, --host            Host to connect to"
  echo "-p, --port            Port to connect to"
  echo
}

while getopts ":ht:c:H:p:" opt; do
  case $opt in
    h|help)
      display_help
      exit 0
      ;;
    H|host)
      HOST=$OPTARG
      ;;
    p|port)
      PORT=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      display_help
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      display_help
      exit 1
      ;;
  esac
done

# Unlock the account
geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 10000)" attach http://$HOST:$PORT

FROM=$(geth --exec "web3.personal.listAccounts[0]" attach http://$HOST:$PORT)

JSOM='
{
    "jsonrpc": "2.0",
    "method": "eth_sendTransaction",
    "params": [{
        "from": '$FROM',
        "to": "0xc0ffee254729296a45a3885639AC7E10F9d54979",
        "value": "0x3e8",
        "gasPrice": "0xa"
    }],
    "id": 1
}
'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TMP_SEND_TX=$SCRIPT_DIR/sendTx.json

rm -f $TMP_SEND_TX
touch $TMP_SEND_TX
echo $JSOM > $TMP_SEND_TX

curl -H "Content-Type: application/json" -X POST --data @$TMP_SEND_TX http://$HOST:$PORT/

rm -f $TMP_SEND_TX
