// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-721 token interface
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract SpecialERC721 is ERC721 {
  constructor() ERC721("SpecialERC721", "SPEC") {
      _safeMint(msg.sender, 0, "This is a Special NFT token -- Woo!");
  }
}
