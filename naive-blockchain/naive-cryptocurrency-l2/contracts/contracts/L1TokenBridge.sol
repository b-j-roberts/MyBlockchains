// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-20 token interface
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract L1TokenBridge {
    address public sequencer;
    mapping(address => bool) public allowedTokens; // Mapping to store allowed ERC-20 token contracts
    
    uint256 public TokenDepositNonce = 0;
    uint256 public TokenWithdrawNonce = 0;

    event TokensDeposited(uint256 nonce, address indexed from, address indexed tokenAddress, uint256 amount);
    event TokensWithdrawn(uint256 nonce, address indexed to, address indexed tokenAddress, uint256 amount);

    constructor(address _sequencer) {
        sequencer = _sequencer;
    }

    // Modifier to ensure only the admin can call certain functions
    modifier onlySequencer() {
        require(msg.sender == sequencer, "Only sequencer can call this");
        _;
    }

    // Function to allow users to lock ERC-20 tokens in this contract
    function DepositTokens(address tokenAddress, uint256 amount) external {
        require(allowedTokens[tokenAddress], "Token not allowed");

        //// Transfer tokens from the user to this contract
        IERC20 tokenContract = IERC20(tokenAddress);
        require(tokenContract.transferFrom(msg.sender, address(this), amount), "Token transfer failed");
        TokenDepositNonce++;
        emit TokensDeposited(TokenDepositNonce, msg.sender, tokenAddress, amount);
    }

    // Function to allow the admin to add ERC-20 tokens to the allowed list
    function addAllowedToken(address tokenAddress) external onlySequencer {
        allowedTokens[tokenAddress] = true;
    }

    // Function to allow the admin to remove ERC-20 tokens from the allowed list
    function removeAllowedToken(address tokenAddress) external onlySequencer {
        allowedTokens[tokenAddress] = false;
    }

    // Function to allow the admin to withdraw ERC-20 tokens from the contract
    function WithdrawTokens(address tokenAddress, address to, uint256 amount) external onlySequencer {
        require(allowedTokens[tokenAddress], "Token not allowed");

        IERC20 tokenContract = IERC20(tokenAddress);
        require(tokenContract.transfer(to, amount), "Token transfer failed");

        TokenWithdrawNonce++;
        emit TokensWithdrawn(TokenWithdrawNonce, to, tokenAddress, amount);
    }

    function GetAllowedToken(address tokenAddress) external view returns (bool) {
        return allowedTokens[tokenAddress];
    }

    function GetTokenDepositNonce() external view returns (uint256) {
      return TokenDepositNonce;
    }

    function GetTokenWithdrawNonce() external view returns (uint256) {
      return TokenWithdrawNonce;
    }
}
