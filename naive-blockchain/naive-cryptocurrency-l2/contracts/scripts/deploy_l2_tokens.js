// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './deploy-deps.js'
import fs from 'fs'

(async () => {
  try {
      const result = await deploy('BasicL2ERC20', [0, process.env.L2_TOKEN_BRIDGE_ADDRESS], process.env.IPC_PATH)
      console.log(result)
      console.log("Deployed BasicL2ERC20 to : ", result.address)
      var jsonOutput = "{\"address\": \"" + result.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-basic-erc20-address.txt', jsonOutput)

      console.log("Deploying StableL2ERC20... with L2 Token Bridge Address: ", process.env.L2_TOKEN_BRIDGE_ADDRESS)
      const result2 = await deploy('StableL2ERC20', [process.env.SEQUENCER_ADDRESS, 0, process.env.L2_TOKEN_BRIDGE_ADDRESS], process.env.IPC_PATH)
      console.log(result2)
      console.log("Deployed StableL2ERC20 to : ", result2.address)
      var jsonOutput2 = "{\"address\": \"" + result2.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-stable-erc20-address.txt', jsonOutput2)

      console.log("Deploying BasicL2ERC721... with L2 Token Bridge Address: ", process.env.L2_TOKEN_BRIDGE_ADDRESS)
      const result3 = await deploy('BasicL2ERC721', [process.env.L2_TOKEN_BRIDGE_ADDRESS, 10], process.env.IPC_PATH)
      console.log(result3)
      console.log("Deployed BasicL2ERC721 to : ", result3.address)
      var jsonOutput3 = "{\"address\": \"" + result3.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-basic-erc721-address.txt', jsonOutput3)

      console.log("Deploying SpecialL2ERC721... with L2 Token Bridge Address: ", process.env.L2_TOKEN_BRIDGE_ADDRESS)
      const result4 = await deploy('SpecialL2ERC721', [process.env.L2_TOKEN_BRIDGE_ADDRESS], process.env.IPC_PATH)
      console.log(result4)
      console.log("Deployed SpecialL2ERC721 to : ", result4.address)
      var jsonOutput4 = "{\"address\": \"" + result4.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-special-erc721-address.txt', jsonOutput4)
  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
