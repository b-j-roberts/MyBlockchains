// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-20 token interface
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./L2TokenMinter.sol";

contract StableL2ERC20 is ERC20, L2TokenMinter {
  address public owner;
  address public bridge;

  event ReturnedTokens(address from, uint256 amount);

  modifier onlyOwner() {
    require(msg.sender == owner, "Only the owner can call this function.");
    _;
  }

  constructor(address _owner, uint256 initialSupply, address _bridge) public ERC20("Stable", "STBL") {
    owner = _owner;
    bridge = _bridge;
    _mint(owner, initialSupply);
  }

  function mint(address account, uint256 amount) public onlyOwner {
    _mint(account, amount);
  }

  //TODO: Should this be allowed on L2?
  function returnToOwner(uint256 amount) public {
    _transfer(msg.sender, owner, amount);
    _burn(owner, amount);
    emit ReturnedTokens(msg.sender, amount);
  }

  function sequencerMint(address account, uint256 amount) external {
    require(msg.sender == bridge, "Only the bridge can call this function.");
    _mint(account, amount);
  }

  function sequencerBurn(address account, uint256 amount) external {
    require(msg.sender == bridge, "Only the bridge can call this function.");
    _burn(account, amount);
  }
}
