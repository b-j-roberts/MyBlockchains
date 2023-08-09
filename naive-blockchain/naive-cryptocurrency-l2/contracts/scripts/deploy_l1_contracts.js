// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './deploy-deps.js'
import fs from 'fs'

(async () => {
  try {
      console.log("Deploying contracts... w/ sequencer address: ", process.env.SEQUENCER_ADDRESS)
      const result = await deploy('TransactionStorage', [process.env.SEQUENCER_ADDRESS], 'http://localhost:8545')
      console.log(result)
      console.log("Deployed TransactionStorage contract to : ", result.address)
      var jsonOutput = "{\"address\": \"" + result.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/tx-storage-address.txt', jsonOutput)

      const result2 = await deploy('L1Bridge', [process.env.SEQUENCER_ADDRESS], 'http://localhost:8545')
      console.log(result2)
      console.log("Deployed L1Bridge contract to : ", result2.address)
      var jsonOutput2 = "{\"address\": \"" + result2.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l1-bridge-address.txt', jsonOutput2)

      const result3 = await deploy('L1TokenBridge', [process.env.SEQUENCER_ADDRESS], 'http://localhost:8545') 
      console.log(result3)
      console.log("Deployed L1TokenBridge contract to : ", result3.address)
      var jsonOutput3 = "{\"address\": \"" + result3.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/l1-token-bridge-address.txt', jsonOutput3)

  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
