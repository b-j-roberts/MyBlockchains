#!/bin/bash
#
# This script is used to test the RPC transaction
#

# Display help text
function display_help() {
  echo "Usage: $0 [OPTIONS]"
  echo ""
  echo "OPTIONS:"
  echo "  -p, --port         Port to send the request on (required)"
  echo "  -o, --host         Host to send the request to (required)"
  echo "  -a, --account      Path to the account directory (required)"
  echo "  -t, --to           Address to send the transaction to (required)"
  echo "  -am, --amount      Amount to send (required)"
  echo "  -f, --fee          Transaction fee (required)"
  echo "  -n, --nonce        Transaction nonce (required)"
  echo "  -h, --help         Show help text"
}

# Define command line flags
while [[ $# -gt 0 ]]
do
key="$1"

case $key in -a|--account)
    account_directory="$2"
    shift # past argument
    shift # past value
    ;;
    -p|--port)
    port="$2"
    shift # past argument
    shift # past value
    ;;
    -o|--port)
    host="$2"
    shift # past argument
    shift # past value
    ;;
    -t|--to)
    to_address="$2"
    shift # past argument
    shift # past value
    ;;
    -am|--amount)
    amount="$2"
    shift # past argument
    shift # past value
    ;;
    -f|--fee)
    fee="$2"
    shift # past argument
    shift # past value
    ;;
    -n|--nonce)
    nonce="$2"
    shift # past argument
    shift # past value
    ;;
    -h|--help)
    display_help
    exit 0
    ;;
    *)    # unknown option
    echo "Unknown option: $1"
    display_help
    exit 1
    ;;
esac
done

# Check if required arguments are provided
if [ -z "$account_directory" ] || [ -z "$to_address" ] || [ -z "$amount" ] || [ -z "$fee" ] || [ -z "$nonce" ] || [ -z "$port" ] || [ -z "$host" ]; then
  echo "Missing required argument(s)."
  display_help
  exit 1
fi

account_address=$(cat $account_directory/account_address.txt)
echo "Creating transaction from $account_address to $to_address with amount $amount and fee $fee and nonce $nonce"

# Example
# ./scripts/rpc-transaction.sh --account accounts/account-12345/ --to 2 --amount 1000 --fee 100 --nonce 1
echo "http://$host:$port/add?value=$account_address,$to_address,$amount,$fee,$nonce,$account_directory/public_key.pem,$account_directory/private_key.pem"
curl -X POST "http://$host:$port/add?value=$account_address,$to_address,$amount,$fee,$nonce,$account_directory/public_key.pem,$account_directory/private_key.pem"
#TODO: Use different scheme for sending transaction
