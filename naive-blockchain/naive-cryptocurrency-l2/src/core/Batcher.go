package core

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

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
  L1NodeUrl string
  L1ContractAddress common.Address
  L1ChainId int
  PosterAddress common.Address
  BatchSize int
  MaxBatchTimeMinutes int
  L1BridgeAddress common.Address
  L1TokenBridgeAddress common.Address
  L2IPCPath string
  ContractsPath string
}
    
   
func NewBatcher(l2Blockchain *core.BlockChain, batcherConfig *BatcherConfig) *Batcher {
  l1Comms, err := l2utils.NewL1Comms(batcherConfig.L1NodeUrl, batcherConfig.L1ContractAddress, batcherConfig.L1BridgeAddress, batcherConfig.L1TokenBridgeAddress, big.NewInt(int64(batcherConfig.L1ChainId)), l2utils.L1TransactionConfig{
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
    CurrL2BlockNumber:     l2Blockchain.CurrentBlock().Number.Uint64(),
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

  var byteDataHash [32]byte
  byteDataHash = sha256.Sum256(transactionByteData)
  log.Printf("Batch %d: %x\n", batcher.BatchId, byteDataHash)

  blockchainRootHash := batcher.L2Blockchain.GetBlockByNumber(batcher.CurrL2BlockNumber).Root()
  log.Printf("Blockchain root: %x\n", blockchainRootHash)

  err = batcher.L1Comms.PostBatch(compressedTransactionByteData, batcher.BatchId, blockchainRootHash, batcher.BatcherConfig.PosterAddress)
  if err != nil {
    return err
  }

  return nil
}

func UnpackEthWithdraw(receiptLog types.Log) (nonce *big.Int, addr common.Address, amount *big.Int, err error) {
    data := receiptLog.Data
    if len(data) < 96 {
        err = fmt.Errorf("invalid data")
        return 
    }

    offset := 12
    nonce = new(big.Int).SetBytes(data[:32])
    addr = common.BytesToAddress(data[32:52+offset])
    amount = new(big.Int).SetBytes(data[52+offset:84+offset])

    return
}

func UnpackTokenWithdraw(receiptLog types.Log) (nonce *big.Int, addr common.Address, tokenAddr common.Address, amount *big.Int, err error) {
  data := receiptLog.Data
  if len(data) < 128 {
    err = fmt.Errorf("invalid data")
    return 
  }

  nonce = new(big.Int).SetBytes(data[:32])
  log.Println("Got nonce:", nonce)
  addr = common.BytesToAddress(data[32:64])
  log.Println("Got addr:", addr)
  tokenAddr = common.BytesToAddress(data[64:96])
  log.Println("Got tokenAddr:", tokenAddr)
  amount = new(big.Int).SetBytes(data[96:128])
  log.Println("Got amount:", amount)

  return
}

func (batcher *Batcher) Start() error {
  runFunc := func() {
  for {
    //TODO: Use similar only finalized block logic for this
    block := batcher.L2Blockchain.GetBlockByNumber(batcher.CurrL2BlockNumber)
    if block != nil {
      log.Printf("Batcher block %d has %d txs\n", batcher.CurrL2BlockNumber, len(block.Transactions()))
      for _, tx := range block.Transactions() {
        log.Printf("Batcher got tx: %v\n", tx)
        batcher.TxBatch = append(batcher.TxBatch, tx)

        //TODO: To function
        rpcIPC, err := rpc.DialIPC(context.Background(), batcher.BatcherConfig.L2IPCPath)
        if err != nil {
          log.Fatalf("Failed to dial ipc: %v", err)
          panic(err)
        }

        backend := ethclient.NewClient(rpcIPC)                                                                                    
        receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
        if err != nil {
          log.Printf("Batcher got error: %v\n", err)
          panic(err)
        }

        log.Printf("Batcher got Receipt: %v\n", receipt)
        receipt_logs := l2utils.ReceiptLogsWithEvent(receipt, crypto.Keccak256Hash([]byte("EthWithdraw(uint256,address,uint256)")).Bytes())
        for _, receipt_log := range receipt_logs {
          //TODO: l2 bridge address hardcoded
          //TODO: To function
         l2BridgeAddressFile := batcher.BatcherConfig.ContractsPath + "/l2-bridge-address.txt"
         l2BridgeAddressBytes, err := ioutil.ReadFile(l2BridgeAddressFile)
         if err != nil {
           log.Fatalf("Failed to read l2 bridge address file: %v", err)
           panic(err)
         }
         var l2BridgeAddressJSONMap map[string]interface{}
         err = json.Unmarshal(l2BridgeAddressBytes, &l2BridgeAddressJSONMap)
         if err != nil {
           log.Fatalf("Failed to unmarshal l2 bridge address json: %v", err)
           panic(err)
         }                                                          
         l2BridgeAddress := common.HexToAddress(l2BridgeAddressJSONMap["address"].(string))                                        
                                                                                                                                   
         if err != nil {          
           log.Fatalf("Failed to create l2 comms: %v", err)
           panic(err)
         }

          if common.HexToAddress(receipt_log.Address.Hex()) == l2BridgeAddress {

            nonce, addr, amount, err := UnpackEthWithdraw(*receipt_log)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            currL1BridgeNonce, err := batcher.L1Comms.BridgeContract.GetEthWithdrawNonce(&bind.CallOpts{})
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

            //TODO: Check nonce
            //TODO: WHat happens if briding more than has / exists? & try and start state with 0 extra tokens / eth
            tx, err := batcher.L1Comms.BridgeContract.WithdrawEth(transactOpts, addr, amount)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }
            log.Printf("Withdrawal tx to l1: %s\n", tx.Hash().Hex())
          }
        }

        receipt_logs = l2utils.ReceiptLogsWithEvent(receipt, crypto.Keccak256Hash([]byte("TokensWithdrawn(uint256,address,address,uint256)")).Bytes())
        log.Println("Got receipt logs:", len(receipt_logs))
        for _, receipt_log := range receipt_logs {
         l2TokenBridgeAddressFile := batcher.BatcherConfig.ContractsPath + "/l2-token-bridge-address.txt"
         l2TokenBridgeAddressBytes, err := ioutil.ReadFile(l2TokenBridgeAddressFile)
         if err != nil {
           log.Fatalf("Failed to read l2 token bridge address file: %v", err)
           panic(err)
         } 
         var l2TokenBridgeAddressJSONMap map[string]interface{}
         err = json.Unmarshal(l2TokenBridgeAddressBytes, &l2TokenBridgeAddressJSONMap)
         if err != nil {
           log.Fatalf("Failed to unmarshal l2 token bridge address json: %v", err)
           panic(err)
         }
         l2TokenBridgeAddress := common.HexToAddress(l2TokenBridgeAddressJSONMap["address"].(string))

          if common.HexToAddress(receipt_log.Address.Hex()) == l2TokenBridgeAddress {
            log.Printf("Got token withdrawal: %v\n", receipt_log)
            nonce, from, token, amount, err := UnpackTokenWithdraw(*receipt_log)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            currL1TokenBridgeNonce, err := batcher.L1Comms.TokenBridgeContract.GetTokenWithdrawNonce(&bind.CallOpts{})
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            if currL1TokenBridgeNonce.Cmp(nonce) >= 0 {
              log.Printf("Skipping token withdrawal: %v %v %v %v\n", nonce, from, token, amount)
              continue
            }

            log.Printf("Token withdrawal: %v %v %v %v\n", nonce, from, token, amount)
            transactOpts, err := batcher.L1Comms.CreateL1TransactionOpts(batcher.BatcherConfig.PosterAddress, big.NewInt(0))
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }

            tx, err := batcher.L1Comms.TokenBridgeContract.WithdrawTokens(transactOpts, token, from, amount)
            if err != nil {
              log.Printf("Batcher got error: %v\n", err)
              panic(err)
            }
            log.Printf("Token withdrawal tx to l1: %s\n", tx.Hash().Hex())
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
