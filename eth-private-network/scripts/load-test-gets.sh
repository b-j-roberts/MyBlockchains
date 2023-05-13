#!/bin/bash
# sudo --preserve-env bash
# ulimit -n 10000
#
# This script will load test basic get rpcs on a geth node

HOST="localhost"
PORT="8545"

TXN_COUNT=100000
THREAD_COUNT=1500

display_help() {
  echo "Usage: load-test-gets.sh [Options]"
  echo
  echo "-h, --help            Display help"
  echo "-t, --txn-count       Number of transactions to send"
  echo "-c, --thread-count    Number of threads to use"
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
    t|txn-count)
      TXN_COUNT=$OPTARG
      ;;
    c|thread-count)
      THREAD_COUNT=$OPTARG
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

FROM=$(geth --exec "web3.personal.listAccounts[0]" attach http://$HOST:$PORT)
JSOM='
{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x'$FROM'"],
    "id": 1
}
'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TMP_SEND_GET=$SCRIPT_DIR/sendGet.json

rm -f $TMP_SEND_GET
touch $TMP_SEND_GET
echo $JSOM > $TMP_SEND_GET

ab -c $THREAD_COUNT -n $TXN_COUNT -p $TMP_SEND_GET -T application/json http://$HOST:$PORT/

rm -f $TMP_SEND_GET
