#!/bin/bash

geth attach --exec "eth.getBlock(\"latest\")" data/geth.ipc
