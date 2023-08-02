#!/bin/bash
#
# This script starts a prover for l1 block bathces posted by the sequencer.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

CONFIG_FILE=$WORK_DIR/configs/sequencer.config.json

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -c, --config               Path to config file ( default: config.json )"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
  #TODO: Examples for all
}

while getopts ":hc:o:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
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

datadir=$(cat $CONFIG_FILE | jq -r '."data-dir"')
proverAddress=$(cat ${datadir}/keystore/* | jq -r '.address')

if [ -z "$OUTPUT_FILE" ]; then
  $WORK_DIR/build/prover --prover-address ${proverAddress} --config $CONFIG_FILE
else
  $WORK_DIR/build/prover --prover-address ${proverAddress} --config $CONFIG_FILE > $OUTPUT_FILE 2>&1 &
fi
