import fs from 'fs'
import Web3 from 'web3'
import web3_contract from 'web3-eth-contract'
const { Contract, ContractSendMethod, Options } = web3_contract
import HDWalletProvider from "truffle-hdwallet-provider"
//import { Contract, ContractSendMethod, Options } from 'web3-eth-contract'

/**
 * Deploy the given contract
 * @param {string} contractName name of the contract to deploy
 * @param {Array<any>} args list of constructor' parameters
 * @param {string} from account used to send the transaction
 * @param {number} gas gas limit
 * @return {Options} deployed contract
 */
export const deploy = async (contractName, args, from, gas) => {

  const provider = new HDWalletProvider(
    process.env.MNEMONIC,
    'https://goerli.infura.io/v3/e1f81b43fa6e46a9a7ec9c48165732b1'
  )
  const web3 = new Web3(provider)
  console.log(`deploying ${contractName}`)

  const contractABI = JSON.parse(fs.readFileSync(`builds/${contractName}.abi`))
  const contractBytecode = fs.readFileSync(`builds/${contractName}.bin`).toString()

  const contract = new web3.eth.Contract(contractABI)

  try {
   const accounts = await web3.eth.getAccounts()
   console.log('accounts', accounts)
  
   const newContractInstance = await contract.deploy({
      data: contractBytecode,
      arguments: args
    })
    .send({
      from: from || accounts[0],
      gas: 10000000,
      //gasPrice: 30000000000
      value: 0,
      gasPrice: '200000'
      //gas: gas || 1500000,
      //gasPrice: '30000000000000'
    })

    console.log(`deployed ${contractName} at ${newContractInstance.options.address}`)
    return newContractInstance.options.address
  } catch (err) {
    console.log(err)
  }
}
