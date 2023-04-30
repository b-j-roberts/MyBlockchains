pragma solidity >=0.8.2 <0.9.0;

contract TransactionStorage {
    uint256 blockCount;

    function storeBlock(uint256 block) public {
        blockCount += 1;
    }

    function GetBlockCount() public view returns (uint256){
        return blockCount;
    }
}
