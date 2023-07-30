// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './private-deps.js'
import fs from 'fs'

(async () => {
  try {
      const result = await deploy('L2Bridge', [process.env.SEQUENCER_ADDRESS], process.env.IPC_PATH)
      console.log(result)
      console.log("Deployed L2Bridge contract to : ", result.address)
      var jsonOutput = "{\"address\": \"" + result.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-bridge-address.txt', jsonOutput)

      const result2 = await deploy('L2TokenBridge', [process.env.SEQUENCER_ADDRESS], process.env.IPC_PATH)
      console.log(result2)
      console.log("Deployed L2TokenBridge contract to : ", result2.address)
      var jsonOutput2 = "{\"address\": \"" + result2.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l2-token-bridge-address.txt', jsonOutput2)
  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
