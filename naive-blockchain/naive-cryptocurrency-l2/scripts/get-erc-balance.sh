#!/bin/bash
#
# This script is used to get the balance of an ERC20 token

RPC="http://localhost:8545"

display_help() {
    echo "Usage: $0 [Options]... " >&2
    echo "NOTE: Long form flags are not supported, but listed for reference" >&2
    echo
    echo "   -h, --help                 show help"
    echo "   -r, --rpc                  rpc endpoint to use (default: $RPC)"
    echo "   -a, --address              address to get balance for"
    echo "   -c, --contract             contract to get balance from (REQUIRED)"
    echo
    echo "Example: $0 -c 0x1234567890123456789012345678901234567890"
    exit 1
}

while getopts ":hr:a:c:" opt; do
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
            ;;
        c | contract)
            CONTRACT=$OPTARG
            CONTRACT=\"$CONTRACT\"
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

  # Strip quotes if they exist
  ADDRESS=$(echo $ADDRESS | sed -e 's/^"//' -e 's/"$//')
  # String leading 0x if it exists
  ADDRESS=$(echo $ADDRESS | sed -e 's/^0x//')
fi

if [ -z "$CONTRACT" ]; then
  echo "No contract specified"
  exit 1
fi

echo $ADDRESS $CONTRACT

JsonRPCRequest='{"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x70a08231000000000000000000000000'$ADDRESS'","to":'$CONTRACT'},"latest"],"id":67}'
echo $JsonRPCRequest
curl -H "Content-Type: application/json" -H "x-qn-api-version: 1" -X POST --data $JsonRPCRequest $RPC
