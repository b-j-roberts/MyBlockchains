package ledger

func computeDifficulty(blockHeight uint32) uint32 {
  if blockHeight < 100 {
    return uint32(4026531840)
  }
  blockTime := GetBlock(blockHeight).BlockHeader.Timestamp - GetBlock(blockHeight - 100).BlockHeader.Timestamp
  expectedTime := uint32(100 * 30) // 100 blocks, 30 seconds each

  return uint32(float32(GetBlock(blockHeight - 1).BlockHeader.Difficulty * expectedTime) / float32(blockTime))
}
