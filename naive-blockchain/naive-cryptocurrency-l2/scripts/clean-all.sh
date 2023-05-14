#!/bin/bash

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"
NAIVE_RPC_DATA="$HOME/naive-rpc-data"

ps aux | grep "naive-sequencer" | grep -v grep | awk '{print $2}' | xargs kill -9
ps aux | grep "naive-rpc" | grep -v grep | awk '{print $2}' | xargs kill -9

rm -rf $NAIVE_SEQUENCER_DATA
rm -rf $NAIVE_RPC_DATA

make clean
