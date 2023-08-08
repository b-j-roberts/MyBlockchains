package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"

	"github.com/ethereum/go-ethereum/common"
)


type Verifier struct {
  L1Comms *l2utils.L1Comms
}

func NewVerifier(l1Comms *l2utils.L1Comms) *Verifier {
  verifier := &Verifier{
    L1Comms: l1Comms,
  }

  return verifier
}

func (v *Verifier) GetBatchProofParams(batchNumber uint64) ([]byte, [32]byte, [32]byte) {
  log.Println("Getting Batch Proof Params...")

  batchL1Block, err := v.L1Comms.L1Contracts.TxStorageContract.GetBatchL1Block(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch L1 block: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }

  block, err := v.L1Comms.L1Client.BlockByNumber(context.Background(), batchL1Block)
  if err != nil {
    log.Fatalf("Failed to get batch L1 block: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }


  var batchData []byte
  for _, tx := range block.Transactions() {
    receipt, err := v.L1Comms.L1Client.TransactionReceipt(context.Background(), tx.Hash())
    if err != nil {
      log.Fatalf("Failed to get transaction receipt: %v", err)
      return nil, [32]byte{}, [32]byte{}
    }

    receipt_logs := l2utils.ReceiptLogsWithEvent(receipt, []byte("BatchStored(uint256, uint256, byte32)"))
    for _, receipt_log := range receipt_logs {
      eventSignature := []byte("BatchStored(uint256, uint256, byte32)")
      if bytes.Equal(receipt_log.Topics[0].Bytes(), eventSignature) && common.HexToAddress(receipt_log.Address.Hex()) == v.L1Comms.L1ContractAddressConfig.TxStorageContractAddress {
        batchStored, err := v.L1Comms.L1Contracts.TxStorageContract.ParseBatchStored(*receipt_log)
        if err != nil {
          log.Fatalf("Failed to unpack log: %v", err)
          return nil, [32]byte{}, [32]byte{}
        }

        if batchStored.Id == big.NewInt(int64(batchNumber)) {
          //This is the correct event / transaction / batch
          id, root, batchData, err := l2utils.UnpackBatchStoredCalldata(tx)
          if err != nil {
            log.Fatalf("Failed to unpack calldata: %v", err)
            return nil, [32]byte{}, [32]byte{}
          }

          log.Println("Batch Proof Params are", id, root, batchData)
        }
      }
    }
  }
  

  batchPreStateRoot, err := v.L1Comms.L1Contracts.TxStorageContract.GetBatchPreStateRoot(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch pre state root: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }

  batchPostStateRoot, err := v.L1Comms.L1Contracts.TxStorageContract.GetBatchPostStateRoot(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch post state root: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }

  log.Println("Batch Proof Params are", batchData, batchPreStateRoot, batchPostStateRoot)
  log.Println("Batch Proof Params Complete!")

  return batchData, batchPreStateRoot, batchPostStateRoot
}

func (v *Verifier) GetProof(batchNumber uint64) []byte {
  log.Println("Getting Proof...")

  proofL1Block, err := v.L1Comms.L1Contracts.TxStorageContract.GetProofL1Block(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get proof: %v", err)
    return nil
  }

  block, err := v.L1Comms.L1Client.BlockByNumber(context.Background(), proofL1Block)
  if err != nil {
    log.Fatalf("Failed to get proof: %v", err)
    return nil
  }

  for _, tx := range block.Transactions() {
    receipt, err := v.L1Comms.L1Client.TransactionReceipt(context.Background(), tx.Hash())
    if err != nil {
      log.Fatalf("Failed to get transaction receipt: %v", err)
      return nil
    }

    receipt_logs := l2utils.ReceiptLogsWithEvent(receipt, []byte("BatchConfirmed(uint256, uint256, byte32)"))
    for _, receipt_log := range receipt_logs {
      if common.HexToAddress(receipt_log.Address.Hex()) == v.L1Comms.L1ContractAddressConfig.TxStorageContractAddress {
        batchConfirmed, err := v.L1Comms.L1Contracts.TxStorageContract.ParseBatchConfirmed(*receipt_log)
        if err != nil {
          log.Fatalf("Failed to unpack log: %v", err)
          return nil
        }

        if batchConfirmed.Id == big.NewInt(int64(batchNumber)) {
          //This is the correct event / transaction / batch
          _, proof, err := l2utils.UnpackProofCalldata(tx)
          if err != nil {
            log.Fatalf("Failed to unpack calldata: %v", err)
            return nil
          }

          return proof
        }
      }
    }
  }

  return nil
}

func (v *Verifier) CheckProof(batchNumber uint64) error {
  log.Println("Checking Proof...")

  batchData, batchPreStateRoot, batchPostStateRoot := v.GetBatchProofParams(batchNumber)

  if batchData == nil || batchPreStateRoot == [32]byte{} || batchPostStateRoot == [32]byte{} {
    log.Println("Batch Proof Params not found, skipping...")
    return nil
  }

  proof := v.GetProof(batchNumber)
  if proof == nil {
    log.Println("Proof not found, skipping...")
    return nil
  }

  //TODO: Verify proof w/ input data

  log.Println("Checking Proof Complete!")
  return nil
}

func main() { os.Exit(mainImp()) }

func mainImp() int {
  osHomeDir, err := os.UserHomeDir()
  configFile := flag.String("config", osHomeDir + "/naive-sequencer-data/sequencer.config.json", "Path to config file")
  batchId := flag.Uint64("batch-id", 0, "Batch ID to verify")
  flag.Parse()

  config, err := l2config.LoadNodeBaseConfig(*configFile)
  if err != nil {
    log.Fatalf("Failed to load config: %v", err)
  }

  l1Comms, err := l2utils.NewL1Comms(config.L1URL, config.Contracts, big.NewInt(int64(config.L1ChainID)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
  }

  verifier := NewVerifier(l1Comms)

  fatalErrChan := make(chan error, 10)
  err = verifier.CheckProof(*batchId)
  if err != nil {
    fatalErrChan <- err
  }

  sigint := make(chan os.Signal, 1)
  signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

  exitCode := 0
  select {
  case err := <-fatalErrChan:
    log.Println("shutting down due to fatal error:", err)
    defer log.Println("shut down")
    exitCode = 1
  case <-sigint:
    log.Println("shutting down due to interrupt")
  }

  close(sigint)

  return exitCode
}
