// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

//TODO: Use abstract contract instead of interface to allow sequencer hardcoding setup
interface L2TokenMinter {
  function sequencerMint(address account, uint256 value) external;
  function sequencerBurn(address account, uint256 value) external;
}
