pragma solidity >=0.8.2 <0.9.0;

contract TransactionStorage {
  uint256 batchCount = 0;
  mapping(uint256 => bytes32) public batchRoots;
  uint256 lastConfirmedBatch = 0;
  mapping(uint256 => bool) public confirmedBatches;

  function StoreBatch(uint256 id, bytes32 root, bytes calldata batchData) public {
    if(id < 0 || id > batchCount) {
      revert("Invalid batch id");
    }
    batchRoots[id] = root;
    batchCount = id + 1;
  }

  function ConfirmBatch(uint256 id) public {
    if(id < 0 || id >= batchCount) {
      revert("Invalid batch id");
    }
    if(id > 0 && !confirmedBatches[id-1]) {
      revert("Previous batch not confirmed");
    }
    confirmedBatches[id] = true;
    lastConfirmedBatch = id;
  }

  function GetBatchCount() public view returns (uint256){
      return batchCount;
  }

  function GetBatchRoot(uint256 id) public view returns (bytes32){
      return batchRoots[id];
  }

  function GetBatchConfirmed(uint256 id) public view returns (bool){
      return confirmedBatches[id];
  }

  function GetLastConfirmedBatch() public view returns (uint256){
      return lastConfirmedBatch;
  }
}
