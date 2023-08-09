#!/bin/bash
#
# This script starts a prover for l1 block bathces posted by the sequencer.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

PROVER_CONFIG_FILE=$WORK_DIR/configs/sequencer.config.json

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 show help"

  echo "   -c, --config               Path to config file ( default: config.json )"
  echo "   -a, --address              Prover address directory ( if not passed, uses keystore from config )"
  echo "   -o, --output               Output file -- If outfile selected, run task as daemon ( default: console )"

  echo
}

echo "Starting prover..."

while getopts ":hc:a:o:" opt; do
  case ${opt} in
    h|--help )
      display_help
      exit 0
      ;;
    c|--config )
      PROVER_CONFIG_FILE=$OPTARG
      ;;
    a|--address )
      PROVER_ADDRESS_FILE=$OPTARG
      ;;
    o|--output )
      PROVER_OUTPUT_FILE=$OPTARG
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

if [ -z "$PROVER_ADDRESS_FILE" ]; then
  datadir=$(cat $PROVER_CONFIG_FILE | jq -r '."data-dir"')
  PROVER_ADDRESS=$(cat ${datadir}/keystore/* | jq -r '.address')
  mkdir -p ${HOME}/.transactor
  for p in  ${datadir}/keystore/*; do cp $p ${HOME}/.transactor; break; done
else
  PROVER_ADDRESS=$(cat ${PROVER_ADDRESS_FILE}/UTC* | jq -r '.address')
  mkdir -p ${HOME}/.transactor
  for p in  ${PROVER_ADDRESS_FILE}/*; do cp $p ${HOME}/.transactor; break; done
fi

if [ -z "$PROVER_OUTPUT_FILE" ]; then
  $WORK_DIR/build/prover --prover-address ${PROVER_ADDRESS} --config $PROVER_CONFIG_FILE
else
  $WORK_DIR/build/prover --prover-address ${PROVER_ADDRESS} --config $PROVER_CONFIG_FILE > $PROVER_OUTPUT_FILE 2>&1 &
fi
