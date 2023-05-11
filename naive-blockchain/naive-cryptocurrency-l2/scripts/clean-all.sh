#!/bin/bash

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"
NAIVE_RPC_DATA="$HOME/naive-rpc-data"

rm -rf $NAIVE_SEQUENCER_DATA
rm -rf $NAIVE_RPC_DATA

make clean

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

$BASEDIR/../../../eth-private-network/scripts/clean-networks.sh