pragma solidity >=0.8.2 <0.9.0;

//TODO: Features to add: Sequencer Only, Indirect, upgradable, etc.
contract L1Bridge {
  uint256 ethDepositNonce;

  event EthDeposited(uint256 nonce, address addr, uint256 amount);

  function DepositEth() public payable {
    ethDepositNonce++;
    emit EthDeposited(ethDepositNonce, msg.sender, msg.value);
  }

  function GetEthDepositNonce() public view returns (uint256) {
    return ethDepositNonce;
  }

  // TODO: Sig is a signature of the following data:
  // keccak256(abi.encodePacked(nonce, addr, amount))
  // where nonce is the nonce of the deposit tx
  // addr is the address to withdraw to
  // amount is the amount to withdraw
  // The signature is signed by the owner of the address
  // This is to prevent front-running attacks
  // which would allow anyone to withdraw the funds

  function WithdrawEth(
 //   uint256 nonce,
    address payable addr,
    uint256 amount
 //   bytes memory sig
  ) public {
  //  require(nonce == ethDepositNonce, "Invalid nonce");
    ethDepositNonce++;
    addr.send(amount);
  }

  function GetBridgeBalance() public view returns (uint256) {
    return address(this).balance;
  }
}
