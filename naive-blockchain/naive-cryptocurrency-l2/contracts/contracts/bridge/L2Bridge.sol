pragma solidity >=0.8.2 <0.9.0;

//TODO: Features to add: upgradable, etc.
contract L2Bridge {
  uint256 ethDepositNonce = 0;
  uint256 ethWithdrawNonce = 0;

  event EthDeposited(uint256 nonce, address addr, uint256 amount);
  event EthWithdraw(uint256 nonce, address addr, uint256 amount);

  address public sequencer;

  modifier onlySequencer() {
    require(msg.sender == sequencer, "Only sequencer can call this function.");
    _;
  }

  constructor(address _sequencer) {
    sequencer = _sequencer;
  }

  function DepositEth(address addr, uint256 amount) external onlySequencer {
    ethDepositNonce++;
    emit EthDeposited(ethDepositNonce, addr, amount);
  }

  function GetEthDepositNonce() public view returns (uint256) {
    return ethDepositNonce;
  }

  function GetEthWithdrawNonce() public view returns (uint256) {
    return ethWithdrawNonce;
  }

  function WithdrawEth() external payable {
    ethWithdrawNonce++;
    // Burn eth
    payable(address(0x0)).transfer(msg.value);
    emit EthWithdraw(ethWithdrawNonce, msg.sender, msg.value);
  }

  function GetBurntBalance() public view returns (uint256) {
    return address(0x0).balance;
  }
}
