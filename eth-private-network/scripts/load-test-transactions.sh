#!/bin/bash
# sudo --preserve-env bash
# ulimit -n 10000
#
# This script is used to load test the transaction throughput of a geth node.

RPC="http://localhost:8545/"
                                                           
TXN_COUNT=100000                                           
THREAD_COUNT=1500

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo "NOTE: Long form flags are not supported, but listed for reference." >&2
  echo "WARNING: You may need to alter the ulimit to allow for more open files. (ulimit -n 10000)" >&2
  echo
  echo "-h, --help            Display help"
  echo "-t, --txn-count       Number of transactions to send (default: 100000)"
  echo "-c, --thread-count    Number of threads to use (default: 1500)"
  echo "-r, --rpc             RPC endpoint (default: http://localhost:8545)"
  echo
  echo "Example: $0 -t 100000 -c 1500 -r http://localhost:8545"
}

while getopts ":ht:c:r:" opt; do
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
    r|rpc)
      RPC=$OPTARG
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
geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 10000)" attach $RPC

FROM=$(geth --exec "web3.personal.listAccounts[0]" attach $RPC)

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

cat $TMP_SEND_TX

ab -c $THREAD_COUNT -n $TXN_COUNT -p $TMP_SEND_TX -T application/json $RPC

rm -f $TMP_SEND_TX
