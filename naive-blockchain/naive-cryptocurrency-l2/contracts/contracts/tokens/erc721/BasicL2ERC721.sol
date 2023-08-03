// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-721 token interface
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "../../L2TokenMinter.sol";

contract BasicL2ERC721 is ERC721, L2TokenMinter {
  uint256 public totalSupply;

  address bridge;

  constructor(address _bridge, uint256 _totalSupply) ERC721("BasicERC721", "BNFT") {
    totalSupply = _totalSupply;
    bridge = _bridge;
  }

  function sequencerMint(address account, uint256 tokenId) external {
    require(msg.sender == bridge, "Only the bridge can mint");
    require(tokenId < totalSupply, "Cannot mint more than total supply");

    bytes memory tokenURI = abi.encodePacked("This is an NFT token w/ id ", tokenId);
    _safeMint(account, tokenId, tokenURI);
  }

  function sequencerBurn(address account, uint256 tokenId) external {
    require(msg.sender == bridge, "Only the bridge can burn");
    require(_exists(tokenId), "Token does not exist");

    _burn(tokenId);
  }
}
