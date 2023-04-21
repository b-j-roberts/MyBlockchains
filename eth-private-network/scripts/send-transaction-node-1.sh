#!/bin/bash

# Unlock the account on node 2
geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 1000)" attach http://127.0.0.1:8545

# Send a transaction to the node 2
geth --exec "eth.sendTransaction({from: web3.personal.listAccounts[0], to: \"0xc0ffee254729296a45a3885639AC7E10F9d54979\", value: '1000', gasPrice: '10'})" attach http://127.0.0.1:8545

# Check txpool
geth --exec "txpool.status" attach data/geth.ipc
geth --exec "txpool.content" attach data/geth.ipc
