// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-20 token interface
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract BasicERC20 is ERC20 {
  constructor(address owner, uint256 initialSupply) ERC20("Basic", "BSC") {
    _mint(owner, initialSupply);
  }
}
