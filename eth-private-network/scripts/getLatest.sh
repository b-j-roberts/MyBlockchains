#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
SCRIPT_DIR=$SCRIPT_DIR/..
geth attach --exec "eth.getBlock(\"latest\")" $SCRIPT_DIR/data/geth.ipc
