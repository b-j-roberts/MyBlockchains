import assert from 'assert';
import Web3 from 'web3';
import fs from 'fs'

const contractName = 'TransactionStorage'
const contractABI = JSON.parse(fs.readFileSync(`builds/${contractName}.abi`))
const contractBytecode = fs.readFileSync(`builds/${contractName}.bin`).toString()

const web3 = new Web3('http://localhost:8545')

var accounts = await web3.eth.getAccounts();
var transactionStorage = await new web3.eth.Contract(contractABI)
  .deploy({data: contractBytecode, arguments: []})
  .send({from: accounts[0]});
console.log(`Deployed ${contractName} to ${transactionStorage.options.address}`)
console.log("Accounts: ", accounts)

describe('TransactionStorage-Private', () => {
  it('Private: deploys a contract', () => {
    assert.ok(transactionStorage.options.address);
  });

  it('Private: has a default transaction count of 0', async function() {
    const blockCount = await transactionStorage.methods.GetBlockCount().call();
    assert.equal(blockCount, 0);
  });

  it('Private: can add a transaction', async function() {
    await transactionStorage.methods.storeBlock(42).send({from: accounts[0]});
    const blockCount = await transactionStorage.methods.GetBlockCount().call();
    assert.equal(blockCount, 1);
  });
});
