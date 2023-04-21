#!/bin/bash
# sudo sh -c "ulimit -n 10000 && exec su $LOGNAME"
# ulimit -n 10000

# Unlock the account on node 2
geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 10000)" attach http://127.0.0.1:8545

FROM=$(cat data/keystore/* | jq -r '.address' | head -n 1)
JSOM='
{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x'$FROM'"],
    "id": 1
}
'

echo $JSOM > scripts/sendGet.json
ab -c 1500 -n 100000 -p scripts/sendGet.json -T application/json http://127.0.0.1:8545/
