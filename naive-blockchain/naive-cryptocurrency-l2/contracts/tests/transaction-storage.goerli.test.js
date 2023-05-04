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
  it('Private: deploys a contract', () => {
    assert.ok(transactionStorage.options.address);
  });

  it('Private: has a default batch count of 0', async function() {
    const batchCount = await transactionStorage.methods.GetBatchCount().call();
    assert.equal(batchCount, 0);
  });

  it('Private: has a default last confirmed batch of 0', async function() {
    const lastConfirmedBatch = await transactionStorage.methods.GetLastConfirmedBatch().call();
    assert.equal(lastConfirmedBatch, 0);
  });

  it('Private: can store a batch', async function() {
    const bytes32Value = web3.utils.asciiToHex("test");
    const bytesValue = web3.utils.asciiToHex("test2");
    await transactionStorage.methods.StoreBatch(0, bytes32Value, bytesValue).send({from: accounts[0]});

    const batchCount = await transactionStorage.methods.GetBatchCount().call();
    assert.equal(batchCount, 1);
  });

  it('Private: can store more than one batch', async function() {
    const bytes32Value = web3.utils.asciiToHex("test");
    const bytesValue = web3.utils.asciiToHex("test2");
    await transactionStorage.methods.StoreBatch(1, bytes32Value, bytesValue).send({from: accounts[0]});
    await transactionStorage.methods.StoreBatch(2, bytes32Value, bytesValue).send({from: accounts[0]});

    const batchCount = await transactionStorage.methods.GetBatchCount().call();
    assert.equal(batchCount, 3);
  });

  // can get batch roots, can confirm blocks, can get last confirmed batch, can revert batches, cannot revert confirmed batches, cannot confirm too fa  r ahead, cannot store batch in wrong place, ...

});
