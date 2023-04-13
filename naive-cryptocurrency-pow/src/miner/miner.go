package miner

import (
	"log"
	"naive-cryptocurrency-pow/src/ledger"
	"naive-cryptocurrency-pow/src/metrics"
)

func MineBlockNonce(blockHeader ledger.BlockHeader) int {
  MaxUint := ^uint32(0)
  threshhold := MaxUint - blockHeader.Difficulty
  for {
    blockHash := ledger.CalculateBlockHeaderHash(blockHeader)
    if(blockHash < threshhold) {
      log.Printf("Nonce %d found with hash 0x%08x", blockHeader.Nonce, blockHash)
      break
    }
    blockHeader.Nonce += 1
  }
  metrics.BlocksMined.Inc()
  return blockHeader.Nonce
}
