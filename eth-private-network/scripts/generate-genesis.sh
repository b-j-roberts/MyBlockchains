#!/bin/bash
#
# This script generates a genesis.json file for a Clique PoA network w/ one PoA agent ( 1st address ).

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR=$SCRIPT_DIR/..

# Setup the default values for the variables
chainId=505
period=1 # 1 second per block ( ~ 12x faster than mainnet )
gasLimit=300000000 # 300M gas per block ( ~10x mainnet )
output=${WORK_DIR}/genesis.json
fees=false

display_help() {
  echo "Usage: $0 [options...]" >&2
  echo "NOTE: Long form flags are not supported, but they are listed for reference"
  echo
  echo "Arguments:"
  echo
  echo "  -h, --help               Show this help message and exit"

  echo "  -c, --chain-id           The chain ID for the new private eth network (default: 505)"
  echo "  -p, --period             The period - ie # of seconds per block - for the new private eth network (default: 1)"
  echo "  -g, --gas-limit          The gas - ie most gas per block  (default: 300000000)"

  echo "  -f, --fees               Include london fork block fees (default: false)"

  echo "  -a, --addrs              A comma separated list of addresses to pre-fund"
  echo "  -b, --balances           A comma separated list of balances to pre-fund ( must match the number & order of addresses )"

  echo "  -o, --output             The output file for the genesis.json file (default: ${WORK_DIR}/genesis.json)"

  echo
  echo "Example: $0 -a 0x00000000001 -b 100000000000000"
}

# Parse the command line arguments
while getopts ":hc:p:g:fa:b:o:" opt; do
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
    f|--fees)
      fees=true
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
#"londonBlock": 0,
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
EOF

if [ "$fees" = true ]; then
  echo "    \"londonBlock\": 0," >> $output
fi

cat <<EOF >> $output
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
