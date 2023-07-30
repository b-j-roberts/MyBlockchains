// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/introspection/ERC165.sol";

contract L1TokenBridge {
    address public sequencer;
    enum TokenTypes { NONE, ERC20, ERC721 }
    mapping(address => TokenTypes) public allowedTokens; // Mapping to store allowed ERC-20 & 721 token contracts
    
    uint256 public TokenDepositNonce = 0;
    uint256 public TokenWithdrawNonce = 0;

    // Events to be emitted when tokens are deposited or withdrawn
    // Nonce : Unique identifier for the deposit or withdrawal
    // TokenAddress : Address of the ERC-20 or 721 token contract
    // Value : Amount of tokens deposited or withdrawn for ERC-20 tokens and token ID for ERC-721 tokens
    event TokensDeposited(uint256 nonce, address indexed from, address indexed tokenAddress, uint256 value);
    event TokensWithdrawn(uint256 nonce, address indexed to, address indexed tokenAddress, uint256 value);

    constructor(address _sequencer) {
        sequencer = _sequencer;
    }

    // Modifier to ensure only the admin can call certain functions
    modifier onlySequencer() {
        require(msg.sender == sequencer, "Only sequencer can call this");
        _;
    }

    // Function to allow users to lock ERC-20 & 721 tokens in this contract
    // Must allow the contract to spend the ERC-20 or 721 tokens first
    // Value : Amount of tokens to be locked for ERC-20 tokens and token ID for ERC-721 tokens
    function DepositTokens(address tokenAddress, uint256 value) external {
        require(allowedTokens[tokenAddress] != TokenTypes.NONE, "Token not allowed");

        //TODO: Think about using erc165 like below and seperate function for erc20?
        //if (ERC165(tokenAddress).supportsInterface(type(IERC721).interfaceId)) {

        // Transfer tokens from the user to this contract
        if (allowedTokens[tokenAddress] == TokenTypes.ERC721) {
            IERC721 tokenContract = IERC721(tokenAddress);
            tokenContract.transferFrom(msg.sender, address(this), value);
        } else if (allowedTokens[tokenAddress] == TokenTypes.ERC20) {
            // Sadly, ERC-165 doesn't work with ERC-20 tokens
            IERC20 tokenContract = IERC20(tokenAddress);
            require(tokenContract.transferFrom(msg.sender, address(this), value), "ERC20 transfer failed");
        } else {
            revert("Token not allowed");
        }

        TokenDepositNonce++;
        emit TokensDeposited(TokenDepositNonce, msg.sender, tokenAddress, value);
    }

    // Function to allow the admin to add ERC-20 & 721 tokens to the allowed list
    function addAllowedToken(address tokenAddress, TokenTypes tokenType) external onlySequencer {
        allowedTokens[tokenAddress] = tokenType;
    }

    // Function to allow the admin to remove ERC-20 & 721 tokens from the allowed list
    function removeAllowedToken(address tokenAddress) external onlySequencer {
        allowedTokens[tokenAddress] = TokenTypes.NONE;
    }

    // Function to allow the admin to withdraw ERC-20 & 721 tokens from the contract
    // Value : Amount of tokens to be withdrawn for ERC-20 tokens and token ID for ERC-721 tokens
    function WithdrawTokens(address tokenAddress, address to, uint256 value) external onlySequencer {
        require(allowedTokens[tokenAddress] != TokenTypes.NONE, "Token not allowed");

        if (allowedTokens[tokenAddress] == TokenTypes.ERC721) {
            IERC721 tokenContract = IERC721(tokenAddress);
            tokenContract.transferFrom(address(this), to, value);
        } else if (allowedTokens[tokenAddress] == TokenTypes.ERC20) {
            IERC20 tokenContract = IERC20(tokenAddress);
            require(tokenContract.transfer(to, value), "Token transfer failed");
        } else {
            revert("Token not allowed");
        }

        TokenWithdrawNonce++;
        emit TokensWithdrawn(TokenWithdrawNonce, to, tokenAddress, value);
    }

    function GetAllowedToken(address tokenAddress) external view returns (bool) {
        return allowedTokens[tokenAddress] != TokenTypes.NONE;
    }

    function GetTokenDepositNonce() external view returns (uint256) {
      return TokenDepositNonce;
    }

    function GetTokenWithdrawNonce() external view returns (uint256) {
      return TokenWithdrawNonce;
    }
}
