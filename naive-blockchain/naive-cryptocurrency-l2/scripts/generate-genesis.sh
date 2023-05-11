#!/bin/bash

# Set the variables for the genesis block
chainId=505
period=1
ADDR1=$1
extradata="0x0000000000000000000000000000000000000000000000000000000000000000${ADDR1}0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
gasLimit=8000000
declare -A alloc=()
alloc[${ADDR1}]="200000000000000"
BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BASE_DIR=$BASE_DIR/..
mkdir -p $BASE_DIR/build
GENESIS_FILE=$BASE_DIR/build/genesis.json
touch $GENESIS_FILE

# Generate the genesis.json file with the variables
cat <<EOF > $GENESIS_FILE
{
  "config": {
    "chainId": $chainId,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
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
      echo "," >> $GENESIS_FILE
    fi
    echo -n "    \"$addr\": { \"balance\": \"${alloc[$addr]}\" }" >> $GENESIS_FILE
    count=$((count+1))
done

cat <<EOF >> $GENESIS_FILE
  }
}
EOF
