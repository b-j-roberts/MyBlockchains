package core

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

type Batcher struct {
  BatcherConfig *BatcherConfig
  
  L1Comms *l2utils.L1Comms
  L2Blockchain *core.BlockChain
  BlockIdx     uint64

  PostedBlockIdx uint64
  LastPostTime time.Time

  TxBatch      []*types.Transaction
  BatchId int64
}  

type BatcherConfig struct {
  L1NodeUrl string
  L1ContractAddress common.Address
  PosterAddress common.Address
  BatchSize int
  MaxBatchTimeMinutes int
}
    
   
func NewBatcher(l2Blockchain *core.BlockChain, batcherConfig *BatcherConfig) *Batcher {
  l1Comms, err := l2utils.NewL1Comms(batcherConfig.L1NodeUrl, batcherConfig.L1ContractAddress, common.HexToAddress("0x0"))
  if err != nil {
    log.Fatalf("Error creating L1 comms: %s\n", err)
  }
  
  //TODO: Load info from L1
  return &Batcher{
    BatcherConfig: batcherConfig,
    L1Comms: l1Comms,  
    L2Blockchain: l2Blockchain,
    BlockIdx:     0,
    PostedBlockIdx: 0,
    LastPostTime: time.Now(), //TODO
    TxBatch:      make([]*types.Transaction, 0),
    BatchId: 0,
  }
}
  
func (batcher *Batcher) PostBatch() error {
  if len(batcher.TxBatch) == 0 {
    return fmt.Errorf("no txs to post")
  }
    
  log.Printf("Posting batch of %d txs\n", len(batcher.TxBatch))

  //TODO: Compress transaction data
  transactionByteData := make([]byte, 0)
  for _, tx := range batcher.TxBatch {
    txBin, err := tx.MarshalBinary()
    if err != nil {
      return err
    }
    transactionByteData = append(transactionByteData, txBin...)
  }

  byteDataHash := sha256.Sum256(transactionByteData)// TODO: Use blockchain root

  err := batcher.L1Comms.PostBatch(transactionByteData, batcher.BatchId, byteDataHash, batcher.BatcherConfig.PosterAddress)
  if err != nil {
    return err
  }

  return nil
}

func (batcher *Batcher) Start() error {
  runFunc := func() {
  for {
    block := batcher.L2Blockchain.GetBlockByNumber(batcher.BlockIdx)
    if block != nil {
      for _, tx := range block.Transactions() {
        batcher.TxBatch = append(batcher.TxBatch, tx)
      }

      if len(batcher.TxBatch) > batcher.BatcherConfig.BatchSize ||
         (len(batcher.TxBatch) > 0 && batcher.BatcherConfig.MaxBatchTimeMinutes > 0 &&
          time.Since(batcher.LastPostTime) > time.Duration(batcher.BatcherConfig.MaxBatchTimeMinutes) * time.Minute && len(batcher.TxBatch) > 0) {
        err := batcher.PostBatch()
        if err != nil {
          panic(err)
        }
        batcher.LastPostTime = time.Now()
        batcher.PostedBlockIdx = batcher.BlockIdx
        batcher.TxBatch = make([]*types.Transaction, 0)
        batcher.BatchId++
      }
      batcher.BlockIdx++
    } else {
      //TODO: Only sleep if caught up
      time.Sleep(100 * time.Millisecond)
    }
  }}

  go runFunc()

  return nil
}
