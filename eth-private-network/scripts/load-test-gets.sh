#!/bin/bash
# sudo sh -c "ulimit -n 10000 && exec su $LOGNAME"
# ulimit -n 10000

# Unlock the account on node 2
./build/bin/geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 10000)" attach http://127.0.0.1:8545

ab -c 1500 -n 100000 -p private-network/sendTransaction.json -T application/json http://127.0.0.1:8545/
