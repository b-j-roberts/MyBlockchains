package core

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

type Batcher struct {
  BatcherConfig *BatcherConfig
  
  L1Comms *l2utils.L1Comms
  L2Blockchain *core.BlockChain
  CurrL2BlockNumber     uint64
  LastPostTime time.Time

  TxBatch      []*types.Transaction
  BatchId int64
}  

type BatcherConfig struct {
  NodeConfig *l2config.NodeBaseConfig
  PosterAddress common.Address

  BatchSize int
  MaxBatchTimeMinutes int
}
    
   
func NewBatcher(l2Blockchain *core.BlockChain, batcherConfig *BatcherConfig) *Batcher {
  l1Comms, err := l2utils.NewL1Comms(batcherConfig.NodeConfig.L1URL, batcherConfig.NodeConfig.Contracts, big.NewInt(int64(batcherConfig.NodeConfig.L1ChainID)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Error creating L1 comms: %s\n", err)
  }
  
  return &Batcher{
    BatcherConfig: batcherConfig,
    L1Comms: l1Comms,  
    L2Blockchain: l2Blockchain,
    CurrL2BlockNumber: l2Blockchain.CurrentBlock().Number.Uint64(),
    LastPostTime: time.Now(), //TODO: load from file? or change this logic
    TxBatch:      make([]*types.Transaction, 0),
    BatchId: 0, //TODO: load from l1?
  }
}
  
func (batcher *Batcher) PostBatch() error {
  if len(batcher.TxBatch) == 0 {
    return fmt.Errorf("no txs to post")
  }
    
  log.Printf("Posting batch of %d txs\n", len(batcher.TxBatch))

  transactionByteData := make([]byte, 0)
  for _, tx := range batcher.TxBatch {
    txBin, err := tx.MarshalBinary()
    if err != nil {
      return err
    }
    transactionByteData = append(transactionByteData, txBin...)
  }
  // Compress transaction data
  compressedTransactionByteData, err := l2utils.CompressTransactionData(transactionByteData)
  if err != nil {
    return err
  }

  //var byteDataHash [32]byte
  //byteDataHash = sha256.Sum256(transactionByteData)
  //log.Printf("Batch %d: %x\n", batcher.BatchId, byteDataHash)

  blockchainRootHash := batcher.L2Blockchain.GetBlockByNumber(batcher.CurrL2BlockNumber).Root()

  err = batcher.L1Comms.PostBatch(compressedTransactionByteData, batcher.BatchId, blockchainRootHash, batcher.BatcherConfig.PosterAddress)
  if err != nil {
    return err
  }

  return nil
}

func (batcher *Batcher) Start() error {
  runFunc := func() {
  for {
    //TODO: Use similar only finalized block logic for this
    block := batcher.L2Blockchain.GetBlockByNumber(batcher.CurrL2BlockNumber)
    if block != nil {
      for _, tx := range block.Transactions() {
        batcher.TxBatch = append(batcher.TxBatch, tx)

        seqURL := "http://" + batcher.BatcherConfig.NodeConfig.Host + ":" + strconv.Itoa(batcher.BatcherConfig.NodeConfig.Port)
        rpc, err := rpc.Dial(seqURL)

        backend := ethclient.NewClient(rpc)                                                                                    
        receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
        if err != nil {
          log.Printf("Batcher got error: %v\n", err)
          panic(err)
        }

        receipt_logs := l2utils.ReceiptLogsWithEvent(receipt, crypto.Keccak256Hash([]byte("EthWithdraw(uint256,address,uint256)")).Bytes())
        for _, receipt_log := range receipt_logs {
          l2AddressConfig := l2utils.CreateL2ContractAddressConfig(batcher.BatcherConfig.NodeConfig.Contracts)
          if common.HexToAddress(receipt_log.Address.Hex()) == l2AddressConfig.BridgeContractAddress {
            nonce, addr, amount, err := l2utils.UnpackEthWithdraw(*receipt_log)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            currL1BridgeNonce, err := batcher.L1Comms.L1Contracts.BridgeContract.GetEthWithdrawNonce(&bind.CallOpts{})
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            if currL1BridgeNonce.Cmp(nonce) >= 0 {
              log.Printf("Skipping withdrawal: %v %v %v\n", nonce, addr, amount)
              continue
            }

            log.Printf("Withdrawal: %v %v %v\n", nonce, addr, amount)
            transactOpts, err := batcher.L1Comms.CreateL1TransactionOpts(batcher.BatcherConfig.PosterAddress, big.NewInt(0))
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            tx, err := batcher.L1Comms.L1Contracts.BridgeContract.WithdrawEth(transactOpts, addr, amount)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }
            log.Printf("Withdrawal tx to l1: %s\n", tx.Hash().Hex())
          }
        }

        receipt_logs = l2utils.ReceiptLogsWithEvent(receipt, crypto.Keccak256Hash([]byte("TokensWithdrawn(uint256,address,address,uint256)")).Bytes())
        for _, receipt_log := range receipt_logs {
          l2AddressConfig := l2utils.CreateL2ContractAddressConfig(batcher.BatcherConfig.NodeConfig.Contracts)
          if common.HexToAddress(receipt_log.Address.Hex()) == l2AddressConfig.TokenBridgeContractAddress {
            nonce, from, token, amount, err := l2utils.UnpackTokenWithdraw(*receipt_log)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            currL1TokenBridgeNonce, err := batcher.L1Comms.L1Contracts.TokenBridgeContract.GetTokenWithdrawNonce(&bind.CallOpts{})
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            if currL1TokenBridgeNonce.Cmp(nonce) >= 0 {
              log.Printf("Skipping token withdrawal: %v %v %v %v\n", nonce, from, token, amount)
              continue
            }

            transactOpts, err := batcher.L1Comms.CreateL1TransactionOpts(batcher.BatcherConfig.PosterAddress, big.NewInt(0))
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            tx, err := batcher.L1Comms.L1Contracts.TokenBridgeContract.WithdrawTokens(transactOpts, token, from, amount)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }
            log.Printf("Token withdrawal: %v %v %v %v %v\n", nonce, from, token, amount, tx.Hash().Hex())
          }
        }
      }

      if len(batcher.TxBatch) > batcher.BatcherConfig.BatchSize ||
         (len(batcher.TxBatch) > 0 && batcher.BatcherConfig.MaxBatchTimeMinutes > 0 &&
          time.Since(batcher.LastPostTime) > time.Duration(batcher.BatcherConfig.MaxBatchTimeMinutes) * time.Minute && len(batcher.TxBatch) > 0) {
        err := batcher.PostBatch()
        if err != nil {
          log.Printf("Batcher got error: %v\n", err)
          panic(err)
        }
        batcher.LastPostTime = time.Now()
        batcher.TxBatch = make([]*types.Transaction, 0)
        batcher.BatchId++
      }
      batcher.CurrL2BlockNumber++
    } else {
      l2CurrentBlock := batcher.L2Blockchain.CurrentBlock()
      if l2CurrentBlock != nil {
        if batcher.CurrL2BlockNumber >= l2CurrentBlock.Number.Uint64() {
          time.Sleep(100 * time.Millisecond)
        }
      }
    }
  }}

  go runFunc()

  return nil
}
