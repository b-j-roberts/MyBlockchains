// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './private-deps.js'
import fs from 'fs'

(async () => {
  try {
      const result = await deploy('TransactionStorage', [])
      console.log(result)
      
      // Write the contract address to a file
      fs.writeFileSync('./builds/contract-address.txt', result)
  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
