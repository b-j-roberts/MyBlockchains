#!/bin/bash
#
# This script is used to send a transaction to L2 chain

HOST="localhost"                                           
PORT="8545"                                                
                                                           
display_help() {
  echo "Usage: send-transaction.sh [option...]"
  echo
  echo "-h, --help            Display help"
  echo "-H, --host            Host to connect to"
  echo "-p, --port            Port to connect to"
  echo "-v, --value           Value to send"
  echo "-t, --to              Address to send to"
  echo "-i, --ipc             IPC path to connect to"
  echo
  echo
}

while getopts ":ht:c:H:p:v:t:i:" opt; do
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
    v|value)
      VALUE=$OPTARG
      ;;
    t|to)
      TO=$OPTARG
      ;;
    i|ipc)
      IPC=$OPTARG
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

#TODO: Check for requireds

# Unlock the account
geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 10000)" attach http://$HOST:$PORT

#FROM=$(geth --exec "web3.personal.listAccounts[0]" attach $IPC)
FROM=$(cat ~/naive-sequencer-data/keystore/* | jq -r .address)

JSOM='
{
    "jsonrpc": "2.0",
    "method": "eth_sendTransaction",
    "params": [{
        "from": '$FROM',
        "to": '$TO',
        "value": '$VALUE',
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

#curl -H "Content-Type: application/json" -X POST --data @$TMP_SEND_TX http://$HOST:$PORT/
echo "Sending tx to $TO with value $VALUE from $FROM on IPC $IPC"
geth attach --exec "eth.sendTransaction({from: '$FROM', to: '$TO', value: '$VALUE'})" $IPC

rm -f $TMP_SEND_TX
