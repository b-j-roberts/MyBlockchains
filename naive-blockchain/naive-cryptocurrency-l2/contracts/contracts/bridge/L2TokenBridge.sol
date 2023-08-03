// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/introspection/ERC165.sol";
import "../L2TokenMinter.sol";

contract L2TokenBridge {
    address public sequencer;
    mapping(address => address) public l1l2TokenMap; // Mapping to store the L1-L2 token address mapping

    uint256 public tokenDepositNonce;
    uint256 public tokenWithdrawNonce;

    //TODO: add indexed for to, from, tokenAddress so that we can filter events?
    event TokensDeposited(uint256 nonce, address to, address tokenAddress, uint256 value);
    event TokensWithdrawn(uint256 nonce, address from, address tokenAddress, uint256 value);

    constructor(address _sequencer) {
        sequencer = _sequencer;
    }

    // Modifier to ensure only the admin can call certain functions
    modifier onlySequencer() {
        require(msg.sender == sequencer, "Only sequencer can call this");
        _;
    }

    // Function to allow users to lock ERC-20 & 721 tokens in this contract
    function MintTokens(address tokenAddress, address to, uint256 value) external onlySequencer {
        require(l1l2TokenMap[tokenAddress] != address(0x0), "Token not allowed");

        L2TokenMinter tokenContract = L2TokenMinter(l1l2TokenMap[tokenAddress]);
        tokenContract.sequencerMint(to, value);

        // Increment the token deposit nonce
        tokenDepositNonce++;
        emit TokensDeposited(tokenDepositNonce, to, tokenAddress, value);
    }

    // Function to allow the admin to add ERC-20 & 721 tokens to the allowed list
    function addAllowedToken(address l1TokenAddress, address l2TokenAddress) external onlySequencer {
        //TODO: use ERC-165 to check if the contract is tokenminter?
        l1l2TokenMap[l1TokenAddress] = l2TokenAddress;
    }

    // Function to allow the admin to remove ERC-20 & 721 tokens from the allowed list
    function removeAllowedToken(address tokenAddress) external onlySequencer {
        delete l1l2TokenMap[tokenAddress];
    }

    // Function to allow the admin to withdraw ERC-20 & 721 tokens from the contract
    function WithdrawTokens(address tokenAddress, uint256 value) external {
        require(l1l2TokenMap[tokenAddress] != address(0x0), "Token not allowed");


        L2TokenMinter tokenContract = L2TokenMinter(l1l2TokenMap[tokenAddress]);
        tokenContract.sequencerBurn(msg.sender, value);

        // Increment the token withdraw nonce
        tokenWithdrawNonce++;
        emit TokensWithdrawn(tokenWithdrawNonce, msg.sender, tokenAddress, value);
    }

    function GetAllowedToken(address tokenAddress) external view returns (bool) {
        return l1l2TokenMap[tokenAddress] != address(0x0);
    }

    function GetL2TokenAddress(address tokenAddress) external view returns (address) {
        return l1l2TokenMap[tokenAddress];
    }

    function GetTokenDepositNonce() external view returns (uint256) {
        return tokenDepositNonce;
    }

    function GetTokenWithdrawNonce() external view returns (uint256) {
        return tokenWithdrawNonce;
    }
}
