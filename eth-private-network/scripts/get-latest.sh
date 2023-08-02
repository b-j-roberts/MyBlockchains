#!/bin/bash
#
# This script is used to get the latest block from an rpc endpoint

HOST="localhost"
PORT="8545"

display_help() {
    echo "Usage: get-latest.sh [Options]... " >&2
    echo "NOTE: Long form flags are not supported, but listed for reference" >&2
    echo
    echo "   -h, --help                 show help"
    echo "   -H, --host                 host (default: localhost)"
    echo "   -p, --port                 port (default: 8545)"
    echo
    echo "Example: $0"
    exit 1
}

while getopts ":hH:p:" opt; do
    case $opt in
        h | help)
            display_help
            exit 0
            ;;
        H | host)
            HOST=$OPTARG
            ;;
        p | port)
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

curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}' http://$HOST:$PORT
