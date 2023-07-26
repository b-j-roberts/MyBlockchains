#!/bin/bash
#
# This script starts the smart contract metrics exporter for l1 blockchain DA contract


display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -s, --sequencer           sequencer address (Required)"
  echo "   -A, --l1-tx-storage-address  l1 contract address (Required)"
  echo "   -B, --l1-bridge-address    l1 bridge address (Required)"
  echo "   -M, --l2-bridge-address    l2 bridge address (Required)"
  echo "   -E, --l1-token-bridge-address   l1 token bridge address (Required)"
  echo "   -F, --l2-token-bridge-address   l2 token bridge address (Required)"
  echo "   -e, --erc20-address        erc20 address (Required)"
  echo "   -f, --l2-erc20-address     l2 erc20 address (Required)"
  echo "   -u, --stable              stable coin address (Required)"
  echo "   -U, --l2-stable           l2 stable coin address (Required)"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  #TODO: Examples for all
}

while getopts ":hs:A:B:M:E:F:e:f:u:U:o:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    s|--sequencer )
      SEQUENCER_ADDRESS=$OPTARG
      ;;
    A|--l1-tx-storage-address )
      L1_CONTRACT_ADDRESS=$OPTARG
      ;;
    B|--l1-bridge-address )
      L1_BRIDGE_ADDRESS=$OPTARG
      ;;
    M|--l2-bridge-address )
      L2_BRIDGE_ADDRESS=$OPTARG
      ;;
    E|--l1-token-bridge-address )
      L1_TOKEN_BRIDGE_ADDRESS=$OPTARG
      ;;
    F|--l2-token-bridge-address )
      L2_TOKEN_BRIDGE_ADDRESS=$OPTARG
      ;;
    e|--erc20-address )
      ERC20_ADDRESS=$OPTARG
      ;;
    f|--l2-erc20-address )
      L2_ERC20_ADDRESS=$OPTARG
      ;;
    u|--stable )
      STABLE_ADDRESS=$OPTARG
      ;;
    U|--l2-stable )
      L2_STABLE_ADDRESS=$OPTARG
      ;;
    o|--output )
      OUTPUT_FILE=$OPTARG
      ;;

    \? )
      echo "Invalid Option: -$OPTARG" 1>&2
      display_help
      exit 1
      ;;
    : )
      echo "Invalid Option: -$OPTARG requires an argument" 1>&2
      display_help
      exit 1
      ;;
  esac
done

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

if [ -z $SEQUENCER_ADDRESS} ]; then
    echo "Missing required argument: -s" 1>&2
    display_help
    exit 1
fi

if [ -z "${L1_CONTRACT_ADDRESS}" ]; then
  # Try and copy over address from build
  L1_CONTRACT_ADDRESS=$(cat $WORK_DIR/contracts/builds/tx-storage-address.txt | jq -r '.address')

  if [ -z "${L1_CONTRACT_ADDRESS}" ]; then
    echo "Missing required argument: -A" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L1_BRIDGE_ADDRESS}" ]; then
  # Try and copy over address from build
  L1_BRIDGE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l1-bridge-address.txt | jq -r '.address')

  if [ -z "${L1_BRIDGE_ADDRESS}" ]; then
    echo "Missing required argument: -B" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L2_BRIDGE_ADDRESS}" ]; then
  # Try and copy over address from build
  L2_BRIDGE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l2-bridge-address.txt | jq -r '.address')

  if [ -z "${L2_BRIDGE_ADDRESS}" ]; then
    echo "Missing required argument: -M" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L1_TOKEN_BRIDGE_ADDRESS}" ]; then
  # Try and copy over address from build
  L1_TOKEN_BRIDGE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l1-token-bridge-address.txt | jq -r '.address')

  if [ -z "${L1_TOKEN_BRIDGE_ADDRESS}" ]; then
    echo "Missing required argument: -E" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L2_TOKEN_BRIDGE_ADDRESS}" ]; then
  # Try and copy over address from build
  L2_TOKEN_BRIDGE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l2-token-bridge-address.txt | jq -r '.address')

  if [ -z "${L2_TOKEN_BRIDGE_ADDRESS}" ]; then
    echo "Missing required argument: -F" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${ERC20_ADDRESS}" ]; then
  # Try and copy over address from build
  ERC20_ADDRESS=$(cat $WORK_DIR/contracts/builds/basic-erc20-address.txt | jq -r '.address')

  if [ -z "${ERC20_ADDRESS}" ]; then
    echo "Missing required argument: -e" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L2_ERC20_ADDRESS}" ]; then
  # Try and copy over address from build
  L2_ERC20_ADDRESS=$(cat $WORK_DIR/contracts/builds/l2-basic-erc20-address.txt | jq -r '.address')

  if [ -z "${L2_ERC20_ADDRESS}" ]; then
    echo "Missing required argument: -f" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${STABLE_ADDRESS}" ]; then
  # Try and copy over address from build
  STABLE_ADDRESS=$(cat $WORK_DIR/contracts/builds/stable-erc20-address.txt | jq -r '.address')

  if [ -z "${STABLE_ADDRESS}" ]; then
    echo "Missing required argument: -u" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${L2_STABLE_ADDRESS}" ]; then
  # Try and copy over address from build
  L2_STABLE_ADDRESS=$(cat $WORK_DIR/contracts/builds/l2-stable-erc20-address.txt | jq -r '.address')

  if [ -z "${L2_STABLE_ADDRESS}" ]; then
    echo "Missing required argument: -U" 1>&2
    display_help
    exit 1
  fi
fi

if [ -z "${OUTPUT_FILE}" ]; then
  $WORK_DIR/build/smart-contract-metrics --l1-tx-storage-address ${L1_CONTRACT_ADDRESS} --l1-bridge-address ${L1_BRIDGE_ADDRESS} --l2-bridge-address ${L2_BRIDGE_ADDRESS} --l1-token-bridge-address ${L1_TOKEN_BRIDGE_ADDRESS} --l2-token-bridge-address ${L2_TOKEN_BRIDGE_ADDRESS} --erc20-address ${ERC20_ADDRESS} --l2-erc20-address ${L2_ERC20_ADDRESS} --stable-erc20-address ${STABLE_ADDRESS} --l2-stable-erc20-address ${L2_STABLE_ADDRESS} --sequencer ${SEQUENCER_ADDRESS}
else
  $WORK_DIR/build/smart-contract-metrics --l1-tx-storage-address ${L1_CONTRACT_ADDRESS} --l1-bridge-address ${L1_BRIDGE_ADDRESS} --l2-bridge-address ${L2_BRIDGE_ADDRESS} --l1-token-bridge-address ${L1_TOKEN_BRIDGE_ADDRESS} --l2-token-bridge-address ${L2_TOKEN_BRIDGE_ADDRESS} --erc20-address ${ERC20_ADDRESS} --l2-erc20-address ${L2_ERC20_ADDRESS} --stable-erc20-address ${STABLE_ADDRESS} --l2-stable-erc20-address ${L2_STABLE_ADDRESS} --sequencer ${SEQUENCER_ADDRESS} > ${OUTPUT_FILE} 2>&1 &
fi
