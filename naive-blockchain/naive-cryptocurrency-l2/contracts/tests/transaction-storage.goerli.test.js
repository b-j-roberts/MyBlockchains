import assert from 'assert';
import Web3 from 'web3';
import HDWalletProvider from "truffle-hdwallet-provider"
import fs from 'fs'
  
const contractName = 'TransactionStorage'
const contractABI = JSON.parse(fs.readFileSync(`builds/${contractName}.abi`))
const contractBytecode = fs.readFileSync(`builds/${contractName}.bin`).toString()

const provider = new HDWalletProvider(
  process.env.MNEMONIC,
  'https://goerli.infura.io/v3/e1f81b43fa6e46a9a7ec9c48165732b1'
)

const web3 = new Web3(provider)

var accounts = await web3.eth.getAccounts();
var transactionStorage = await new web3.eth.Contract(contractABI)
  .deploy({data: contractBytecode, arguments: []})
  .send({from: accounts[0]});
console.log(`Deployed ${contractName} to ${transactionStorage.options.address}`)

describe('TransactionStorage-Goerli', () => {
  it('Goerli: deploys a contract', () => {
    assert.ok(transactionStorage.options.address);
  });

  it('Goerli: has a default transaction count of 0', async function() {
    const blockCount = await transactionStorage.methods.GetBlockCount().call();
    assert.equal(blockCount, 0);
  });

  it('Goerli: can add a transaction', async function() {
    await transactionStorage.methods.storeBlock(42).send({from: accounts[0]});
    const blockCount = await transactionStorage.methods.GetBlockCount().call();
    assert.equal(blockCount, 1);
  });
});
