#!/bin/bash


# This script is used to quickly launch evenrything needed to run the L2, bridge, and test the system.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

echo "Cleaning L1"
ps aux | grep eth-private-network/scripts/../go-ethereum/build/bin/geth | awk '{print $2}' | xargs kill

echo "Cleaning L2"
ps aux | grep scripts/../build/sequencer | awk '{print $2}' | xargs kill

echo "Cleaning smart contract metrics"
ps aux | grep scripts/../build/smart-contract-metrics | awk '{print $2}' | xargs kill

echo "Cleaning prover"
ps aux | grep scripts/../build/prover | awk '{print $2}' | xargs kill

echo "Cleaning l2 rpc"
ps aux | grep scripts/../build/rpc | awk '{print $2}' | xargs kill
