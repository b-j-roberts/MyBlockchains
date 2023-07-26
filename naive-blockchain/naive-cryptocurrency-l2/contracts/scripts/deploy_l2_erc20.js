// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './private-deps.js'
import fs from 'fs'

(async () => {
  try {
      console.log("Deploying BasicL2ERC20... with L2 Token Bridge Address: ", process.env.L2_TOKEN_BRIDGE_ADDRESS)
      const result = await deploy('BasicL2ERC20', [0, process.env.L2_TOKEN_BRIDGE_ADDRESS], '/home/brandon/naive-sequencer-data/naive-sequencer.ipc')
      console.log(result)
      console.log("Deployed BasicL2ERC20 to : ", result.address)
      var jsonOutput = "{\"address\": \"" + result.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-basic-erc20-address.txt', jsonOutput)
      //Sleep for 3 seconds to allow the contract to be deployed

      //TODO: Check above succeeded
      //TODO: Call to allowToken on l1&2
      console.log("Deploying StableL2ERC20... with L2 Token Bridge Address: ", process.env.L2_TOKEN_BRIDGE_ADDRESS)
      const result2 = await deploy('StableL2ERC20', [process.env.SEQUENCER_ADDRESS, 0, process.env.L2_TOKEN_BRIDGE_ADDRESS], '/home/brandon/naive-sequencer-data/naive-sequencer.ipc')
      console.log(result2)
      console.log("Deployed StableL2ERC20 to : ", result2.address)
      var jsonOutput2 = "{\"address\": \"" + result2.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-stable-erc20-address.txt', jsonOutput2)
  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
