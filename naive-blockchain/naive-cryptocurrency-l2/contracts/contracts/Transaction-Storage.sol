pragma solidity >=0.8.2 <0.9.0;

//TODO: Features to add: Sequencer Only/Indirect?, upgradeable
contract TransactionStorage {
  uint256 batchCount = 0;
  mapping(uint256 => bytes32) public batchRoots;

  uint256 lastConfirmedBatch = 0;
  mapping(uint256 => bool) public confirmedBatches;
  mapping(uint256 => uint64) public batchRewards;
  mapping(uint256 => uint256) public batchL1Block;
  mapping(uint256 => uint256) public proofL1Block;

  event BatchStored(uint256 id, uint256 l1Block, bytes32 root);
  event BatchConfirmed(uint256 id, uint256 l1Block, bytes32 root);

  //TODO: Use on deploy
  function StoreGenesisState(bytes32 root) public {
    if (batchCount > 0) {
      revert("Genesis state already stored");
    }

    batchRoots[0] = root;
    batchCount = 1;
    batchL1Block[0] = block.number;

    confirmedBatches[0] = true;
    batchRewards[0] = 0;
    batchL1Block[0] = block.number;
    lastConfirmedBatch = 0;
  }

  function StoreBatch(uint256 id, bytes32 root, bytes calldata batchData) public {
    //TODO: rewinding?
    if(id <= 0 || id > batchCount) {
      revert("Invalid batch id");
    }
    batchRoots[id] = root;
    batchCount = id + 1;
    batchL1Block[id] = block.number;

    emit BatchStored(id, block.number, root);
  }

  function SubmitProof(uint256 id, bytes calldata proof) public {
    if(id < 0 || id >= batchCount) {
      revert("Invalid batch id");
    }
    if(confirmedBatches[id]) {
      revert("Batch already confirmed");
    }

    if(proof.length > 0) {
      //TODO: COnfirm proof

      confirmedBatches[id] = true;
      lastConfirmedBatch = id;
      batchRewards[id] = 100;
    } else {
      revert("Invalid proof");
    }

    emit BatchConfirmed(id, block.number, batchRoots[id]);
  }

  function GetBatchCount() public view returns (uint256){
      return batchCount;
  }

  function GetBatchRoot(uint256 id) public view returns (bytes32){
      return batchRoots[id];
  }

  function GetBatchPostStateRoot(uint256 id) public view returns (bytes32){
      return batchRoots[id];
  }

  function GetBatchPreStateRoot(uint256 id) public view returns (bytes32){
      require(id > 0, "Invalid batch id");
      return batchRoots[id-1];
  }

  function GetBatchConfirmed(uint256 id) public view returns (bool){
      return confirmedBatches[id];
  }

  function GetLastConfirmedBatch() public view returns (uint256){
      return lastConfirmedBatch;
  }

  function GetReward(uint256 id) public view returns (uint64){
      return batchRewards[id];
  }

  function GetBatchL1Block(uint256 id) public view returns (uint256){
      return batchL1Block[id];
  }

  function GetProofL1Block(uint256 id) public view returns (uint256){
      return proofL1Block[id];
  }
}
