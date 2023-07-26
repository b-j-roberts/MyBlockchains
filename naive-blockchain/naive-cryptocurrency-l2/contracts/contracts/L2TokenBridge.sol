// SPDX-License-Identifier: MIT
pragma solidity >=0.8.2 <0.9.0;

// Import the ERC-20 token interface
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./L2ERC20Minter.sol";

contract L2TokenBridge {
    address public sequencer;
    mapping(address => address) public l1l2TokenMap; // Mapping to store the L1-L2 token address mapping
    uint256 public tokenDepositNonce; // Nonce to track deposits of ERC-20 tokens
    uint256 public tokenWithdrawNonce; // Nonce to track withdrawals of ERC-20 tokens

    event TokensDeposited(uint256 nonce, address to, address tokenAddress, uint256 amount);
    event TokensWithdrawn(uint256 nonce, address from, address tokenAddress, uint256 amount);

    constructor(address _sequencer) {
        sequencer = _sequencer;
    }

    // Modifier to ensure only the admin can call certain functions
    modifier onlySequencer() {
        require(msg.sender == sequencer, "Only sequencer can call this");
        _;
    }

    // Function to allow users to lock ERC-20 tokens in this contract
    function MintTokens(address tokenAddress, address to, uint256 amount) external onlySequencer {
      // TODO: Deployer should be able to set the allowedTokens after deploying
        require(l1l2TokenMap[tokenAddress] != address(0x0), "Token not allowed");

        // Transfer tokens from the user to this contract
        //TODO: Use ERC-165 to check if token is ERC-20 ( w/ l2 compat? )
        L2ERC20Minter tokenContract = L2ERC20Minter(l1l2TokenMap[tokenAddress]);
        //TODO: Mint, not transfer -- Create new l2 bridge erc 20 which allows minting & burning only by sequencer
        tokenContract.sequencerMint(to, amount);

        // Increment the token deposit nonce
        tokenDepositNonce++;
        emit TokensDeposited(tokenDepositNonce, to, tokenAddress, amount);
    }

    // Function to allow the admin to add ERC-20 tokens to the allowed list
    function addAllowedToken(address l1TokenAddress, address l2TokenAddress) external onlySequencer {
        l1l2TokenMap[l1TokenAddress] = l2TokenAddress;
    }

    // Function to allow the admin to remove ERC-20 tokens from the allowed list
    function removeAllowedToken(address tokenAddress) external onlySequencer {
        delete l1l2TokenMap[tokenAddress];
    }

    // Function to allow the admin to withdraw ERC-20 tokens from the contract
    function WithdrawTokens(address tokenAddress, uint256 amount) external {
        require(l1l2TokenMap[tokenAddress] != address(0x0), "Token not allowed");


        L2ERC20Minter tokenContract = L2ERC20Minter(l1l2TokenMap[tokenAddress]);
        //TODO: Burn, not transfer
        address burnAddress = address(0x0);
        tokenContract.sequencerBurn(msg.sender, amount);

        // Increment the token withdraw nonce
        tokenWithdrawNonce++;
        emit TokensWithdrawn(tokenWithdrawNonce, msg.sender, tokenAddress, amount);
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
