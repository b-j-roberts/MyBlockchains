// This script can be used to deploy the "Storage" contract using Web3 library.
// Please make sure to compile "./contracts/1_Storage.sol" file before running this script.
// And use Right click -> "Run" from context menu of the file to run the script. Shortcut: Ctrl+Shift+S

import { deploy } from './deploy-deps.js'
import fs from 'fs'

(async () => {
  try {
      const result4 = await deploy('BasicERC20', [process.env.SEQUENCER_ADDRESS, 100000000000000], 'http://localhost:8545')
      console.log(result4)
      console.log("Deployed BasicERC20 contract to : ", result4.address)
      var jsonOutput4 = "{\"address\": \"" + result4.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/basic-erc20-address.txt', jsonOutput4)

      const result5 = await deploy('StableERC20', [process.env.SEQUENCER_ADDRESS, 100000000000000], 'http://localhost:8545')
      console.log(result5)
      console.log("Deployed StableERC20 contract to : ", result5.address)
      var jsonOutput5 = "{\"address\": \"" + result5.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/stable-erc20-address.txt', jsonOutput5)

      const result6 = await deploy('BasicERC721', [5, 10], 'http://localhost:8545')
      console.log(result6)
      console.log("Deployed BasicERC721 contract to : ", result6.address)
      var jsonOutput6 = "{\"address\": \"" + result6.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/basic-erc721-address.txt', jsonOutput6)

      const result7 = await deploy('SpecialERC721', [], 'http://localhost:8545')
      console.log(result7)
      console.log("Deployed SpecialERC721 contract to : ", result7.address)
      var jsonOutput7 = "{\"address\": \"" + result7.address + "\"}"
      // Write the contract address to a file
      fs.writeFileSync('./builds/special-erc721-address.txt', jsonOutput7)
  } catch (e) {
      console.log(e.message)
  }

  process.exit(0)
})()
