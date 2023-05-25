import fs from 'fs'
import Web3 from 'web3'
import web3_contract from 'web3-eth-contract'
const { Contract, ContractSendMethod, Options } = web3_contract
//import HDWalletProvider from "truffle-hdwallet-provider"
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
  console.log('deploying', contractName, 'with args', args, 'from', from, 'gas', gas)

  const web3 = new Web3('http://localhost:8545')
  console.log(`deploying ${contractName}`)

  const abiPath = `./builds/contracts_${contractName}_sol_${contractName}.abi`
  const abi = JSON.parse(fs.readFileSync(abiPath, 'utf8'))
  const bytecodePath = `./builds/contracts_${contractName}_sol_${contractName}.bin`
  const bytecode = fs.readFileSync(bytecodePath, 'utf8')

  const contract = new web3.eth.Contract(abi)
  const accounts = await web3.eth.getAccounts()

  const contractSend = contract.deploy({
    data: bytecode,
    arguments: args
  })

  const newContractInstance = await contractSend.send({
    from: from || accounts[0],
    gas: gas || 15000000,
  })

  return newContractInstance.options
}
