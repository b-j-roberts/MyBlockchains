// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-20 token interface
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./L2TokenMinter.sol";

contract BasicL2ERC20 is ERC20, L2TokenMinter {
  address public bridge;

  constructor(uint256 initialSupply, address _bridge) ERC20("Basic", "BSC") {
    bridge = _bridge;
    _mint(msg.sender, initialSupply);
  }

  function sequencerMint(address account, uint256 amount) external {
    require(msg.sender == bridge, "Only the bridge can call this function");
    _mint(account, amount);
  }

  function sequencerBurn(address account, uint256 amount) external {
    require(msg.sender == bridge, "Only the bridge can call this function");
    _burn(account, amount);
  }
}
