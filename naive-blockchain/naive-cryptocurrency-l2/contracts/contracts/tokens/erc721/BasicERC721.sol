// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-721 token interface
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract BasicERC721 is ERC721 {
  uint256 public tokenCounter;
  uint256 public totalSupply;

  constructor(uint256 ownerSupply, uint256 _totalSupply) ERC721("BasicERC721", "BNFT") {
    require(ownerSupply <= _totalSupply, "Owner supply cannot be greater than total supply");

    // Mint the owner supply
    for (uint256 i = 0; i < ownerSupply; i++) {
      bytes memory tokenURIBytes = abi.encodePacked("This is an NFT token w/ id ", i);
      _safeMint(msg.sender, tokenCounter, tokenURIBytes);
      tokenCounter++;
    }

    totalSupply = _totalSupply;
  }

  function mint() public {
    require(tokenCounter < totalSupply, "Cannot mint more tokens than the total supply");

    bytes memory tokenURI = abi.encodePacked("This is an NFT token w/ id ", tokenCounter);
    _safeMint(msg.sender, tokenCounter, tokenURI);
    tokenCounter++;
  }
}
