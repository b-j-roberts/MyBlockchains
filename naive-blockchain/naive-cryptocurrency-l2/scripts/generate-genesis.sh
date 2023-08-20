#!/bin/bash
#
# This script generates a genesis.json file for a Clique PoA network.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR=$SCRIPT_DIR/..

# Setup the default values for the variables
chainId=505
period=1 # 1 second per block ( ~ 12x faster than mainnet )
gasLimit=300000000 # 300M gas per block ( ~10x mainnet )
output=${WORK_DIR}/genesis.json

display_help() {
  echo "Usage: generate-genesis.sh [options...]"
  echo "Arguments:"
  echo "  -h, --help: Display this help message"

  echo "  -c, --chain-id: The chain ID to use for the chain"
  echo "  -p, --period: The period to use for the chain (# of seconds between blocks)"
  echo "  -g, --gas-limit: The gas limit to use for the chain"

  echo "  -a, --addrs: The addresses to pre-fund with ether"
  echo "  -b, --balances: The balances to pre-fund the addresses with"

  echo "  -o, --output: The output file to write the genesis.json to"

  echo "Example: ./scripts/generate-genesis.sh -a 0x00000000001 -b 100000000000000"
}

# Parse the command line arguments
while getopts ":hc:p:g:a:b:o:" opt; do
  case ${opt} in
    h|--help)
      display_help
      exit 0
      ;;
    c|--chain-id)
      chainId=$OPTARG
      ;;
    p|--period)
      period=$OPTARG
      ;;
    g|--gas-limit)
      gasLimit=$OPTARG
      ;;
    a|--addrs)
      addrs=$OPTARG
      ;;
    b|--balances)
      balances=$OPTARG
      ;;
    o|--output)
      output=$OPTARG
      ;;
    \?)
      echo "Invalid option: $OPTARG" 1>&2
      display_help
      exit 1
      ;;
    :)
      echo "Invalid option: $OPTARG requires an argument" 1>&2
      display_help
      exit 1
      ;;
  esac
done

if [[ -z "$addrs" || -z "$balances" ]]; then
  echo "Missing required argument: -a and -b are required"
  display_help
  exit 1
fi

# Split the addresses and balances into arrays
IFS=',' read -r -a addrArray <<< "$addrs"
IFS=',' read -r -a balanceArray <<< "$balances"

# Check that the number of addresses and balances match
if [[ ${#addrArray[@]} != ${#balanceArray[@]} ]]; then
  echo "The number of addresses and balances do not match"
  exit 1
fi

# Ensure output is not a directory
if [[ -d "$output" ]]; then
  echo "Output must be a file, not a directory"
  exit 1
fi

rm -f $output
touch $output

# Generate the genesis.json file with the variables
extradata="0x0000000000000000000000000000000000000000000000000000000000000000${addrArray[0]}0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
declare -A alloc=()
for ((i=0;i<${#addrArray[@]};++i)); do
  alloc[${addrArray[$i]}]="${balanceArray[$i]}"
done

# Generate the genesis.json file with the variables
cat <<EOF > $output
{
  "config": {
    "chainId": $chainId,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": $period,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "$gasLimit",
  "extradata": "$extradata",
  "alloc": {
EOF

count=0
for addr in "${!alloc[@]}"; do
    if (($count > 0)); then
      echo "," >> $output
    fi
    echo -n "    \"$addr\": { \"balance\": \"${alloc[$addr]}\" }" >> $output
    count=$((count+1))
done

cat <<EOF >> $output
  }
}
EOF
