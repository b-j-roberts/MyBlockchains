#!/bin/bash
#
# This script starts the smart contract metrics exporter for l1 blockchain DA contract

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

CONFIG_FILE="${WORK_DIR}/configs/sequencer.config.json"

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -s, --sequencer           sequencer address (Required)"
  echo "   -c, --config              config file (default: $WORK_DIR/configs/sequencer.config.json)"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  #TODO: Examples for all
}

while getopts ":hs:c:o:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    s|--sequencer )
      SEQUENCER_ADDRESS=$OPTARG
      ;;
    c|--config )
      CONFIG_FILE=$OPTARG
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

if [ -z $SEQUENCER_ADDRESS} ]; then
    echo "Missing required argument: -s" 1>&2
    display_help
    exit 1
fi

echo "Starting smart contract metrics exporter w/ config: ${CONFIG_FILE} and sequencer: ${SEQUENCER_ADDRESS}"

if [ -z "${OUTPUT_FILE}" ]; then
  $WORK_DIR/build/smart-contract-metrics --config ${CONFIG_FILE} --sequencer ${SEQUENCER_ADDRESS}
else
  $WORK_DIR/build/smart-contract-metrics --config ${CONFIG_FILE} --sequencer ${SEQUENCER_ADDRESS} > ${OUTPUT_FILE} 2>&1 &
fi
