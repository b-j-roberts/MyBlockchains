#!/bin/bash
# sudo sh -c "ulimit -n 10000 && exec su $LOGNAME"
# ulimit -n 10000

# Unlock the account on node 2
geth --exec "web3.personal.unlockAccount(web3.personal.listAccounts[0], \"password\", 10000)" attach http://127.0.0.1:8545

FROM=$(cat data/keystore/* | jq -r '.address' | head -n 1)
JSOM='
{
    "jsonrpc": "2.0",
    "method": "eth_sendTransaction",
    "params": [{
        "from": "0x'$FROM'",
        "to": "0xc0ffee254729296a45a3885639AC7E10F9d54979",
        "value": "0x3e8",
        "gasPrice": "0xa"
    }],
    "id": 1
}
'

echo $JSOM > scripts/sendTransaction.json
ab -c 1500 -n 100000 -p scripts/sendTransaction.json -T application/json http://127.0.0.1:8545/
