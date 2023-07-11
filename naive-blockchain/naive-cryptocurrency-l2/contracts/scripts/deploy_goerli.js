// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './goerli-deps.js'
import fs from 'fs'

(async () => {
  try {
      const result = await deploy('TransactionStorage', [process.env.SEQUENCER_ADDRESS])
      console.log(result)
      console.log("Deployed TransactionStorage contract to : ", result.address)
      var jsonOutput = "{\"address\": \"" + result.address + "\"}"
      fs.writeFileSync("./builds/tx-storage-address.txt", jsonOutput)

      const result2 = await deploy('L1Bridge', [process.env.SEQUENCER_ADDRESS])
      console.log(result2)
      console.log("Deployed L1Bridge contract to : ", result2.address)
      var jsonOutput2 = "{\"address\": \"" + result2.address + "\"}"
      fs.writeFileSync("./builds/l1-bridge-address.txt", jsonOutput2)
  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
