// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-20 token interface
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract StableERC20 is ERC20 {
  address public owner;

  //TODO: indexed?
  event ReturnedTokens(address from, uint256 amount);

  modifier onlyOwner() {
    require(msg.sender == owner, "Only the owner can call this function.");
    _;
  }

  constructor(address _owner, uint256 initialSupply) public ERC20("Stable", "STBL") {
    owner = _owner;
    _mint(owner, initialSupply); //TODO: Initial supply should be 0
  }

  function mint(address account, uint256 amount) public onlyOwner {
    _mint(account, amount);
  }

  function returnToOwner(uint256 amount) public {
    _transfer(msg.sender, owner, amount);
    _burn(owner, amount);
    emit ReturnedTokens(msg.sender, amount);
  }
}
