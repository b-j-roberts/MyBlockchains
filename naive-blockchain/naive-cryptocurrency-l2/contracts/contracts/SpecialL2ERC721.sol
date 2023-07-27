// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-721 token interface
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "./L2TokenMinter.sol";

contract SpecialL2ERC721 is ERC721, L2TokenMinter {
  address bridge;

  constructor(address _bridge) ERC721("SpecialERC721", "SPEC") {
    bridge = _bridge;
  }

  function sequencerMint(address account, uint256 tokenId) external {
    require(msg.sender == bridge, "Only the bridge can mint");
    require(tokenId == 0, "Cannot mint more than one token");

    bytes memory tokenURI = abi.encodePacked("This is a Special NFT token -- Woo!");
    _safeMint(account, tokenId, tokenURI);
  }

  function sequencerBurn(address account, uint256 tokenId) external {
    require(msg.sender == bridge, "Only the bridge can burn");
    require(_exists(tokenId), "Token does not exist");

    _burn(tokenId);
  }
}
