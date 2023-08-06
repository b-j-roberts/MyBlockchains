#!/bin/bash
#
# This script is used to get the latest block from an rpc endpoint

RPC="http://localhost:8545"

display_help() {
    echo "Usage: $0 [Options]... " >&2
    echo "NOTE: Long form flags are not supported, but listed for reference" >&2
    echo
    echo "   -h, --help                 show help"
    echo "   -r, --rpc                  rpc endpoint to use (default: $RPC)"
    echo "   -a, --address              address to get balance for"
    echo
    echo "Example: $0"
    exit 1
}

while getopts ":hr:a:" opt; do
    case $opt in
        h | help)
            display_help
            exit 0
            ;;
        r | rpc)
            RPC=$OPTARG
            ;;
        a | address)
            ADDRESS=$OPTARG
            ADDRESS=\"$ADDRESS\"
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

if [ -z "$ADDRESS" ]; then
  ADDRESS=$(geth --exec "web3.personal.listAccounts[0]" attach $RPC)

  if [ -z "$ADDRESS" ]; then
    echo "No address specified and no accounts found in geth"
    exit 1
  fi
fi

echo $ADDRESS
JsonRPCRequest='{"jsonrpc":"2.0","method":"eth_getBalance","params":['$ADDRESS',"latest"],"id":1}'
echo $JsonRPCRequest
curl -H "Content-Type: application/json" -X POST --data $JsonRPCRequest $RPC
