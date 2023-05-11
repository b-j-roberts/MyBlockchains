package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"naive-l2/src/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)


type Prover struct {
  L1Comms *utils.L1Comms
  ProverL1Address common.Address

  LastProvedBatch uint64
  TotalProofs uint64
  TotalProofsVerifiedOnchain uint64
  TotalRewards uint64
}

func NewProver(l1Comms *utils.L1Comms, proverL1Address common.Address) *Prover {
  prover := &Prover{
    L1Comms: l1Comms,
    ProverL1Address: proverL1Address,

    TotalProofs: 0,
    TotalProofsVerifiedOnchain: 0,
    TotalRewards: 0,
  }

  return prover
}

func unpackCalldata(tx *types.Transaction) (id *big.Int, root [32]byte, batchData []byte, err error) {
    data := tx.Data()
    if len(data) < 4 {
        err = errors.New("Invalid calldata length")
        return
    }

    id = new(big.Int).SetBytes(data[4:36])
    copy(root[:], data[36:68])
    batchData = data[68:]

    return
}

func (p *Prover) GetBatchProofParams(batchNumber uint64) ([]byte, [32]byte, [32]byte) {
  log.Println("Getting Batch Proof Params...")

  batchL1Block, err := p.L1Comms.L1Contract.GetBatchL1Block(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch L1 block: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }

  block, err := p.L1Comms.L1Client.BlockByNumber(context.Background(), batchL1Block)
  if err != nil {
    log.Fatalf("Failed to get batch L1 block: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }


  var batchData []byte
  for _, tx := range block.Transactions() {
    receipt, err := p.L1Comms.L1Client.TransactionReceipt(context.Background(), tx.Hash())
    if err != nil {
      log.Fatalf("Failed to get transaction receipt: %v", err)
      return nil, [32]byte{}, [32]byte{}
    }

    for _, receipt_log := range receipt.Logs {
      eventSignature := []byte("BatchStored(uint256, uint256, byte32)")
      if bytes.Equal(receipt_log.Topics[0].Bytes(), eventSignature) && common.HexToAddress(receipt_log.Address.Hex()) == p.L1Comms.L1ContractAddress {
        batchStored, err := p.L1Comms.L1Contract.ParseBatchStored(*receipt_log)
        if err != nil {
          log.Fatalf("Failed to unpack log: %v", err)
          return nil, [32]byte{}, [32]byte{}
        }

        if batchStored.Id == big.NewInt(int64(batchNumber)) {
          //This is the correct event / transaction / batch
          id, root, batchData, err := unpackCalldata(tx)
          if err != nil {
            log.Fatalf("Failed to unpack calldata: %v", err)
            return nil, [32]byte{}, [32]byte{}
          }

          log.Println("Batch Proof Params are", id, root, batchData)
        }
      }
    }
  }

  batchPreStateRoot, err := p.L1Comms.L1Contract.GetBatchPreStateRoot(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch pre state root: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }

  batchPostStateRoot, err := p.L1Comms.L1Contract.GetBatchPostStateRoot(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch post state root: %v", err)
    return nil, [32]byte{}, [32]byte{}
  }

  log.Println("Batch Proof Params are", batchData, batchPreStateRoot, batchPostStateRoot)
  log.Println("Batch Proof Params Complete!")

  return batchData, batchPreStateRoot, batchPostStateRoot
}

func (p *Prover) CreateAggProof(batchBytes []byte, batchNumber uint64, batchPreStateRoot []byte, batchPostStateRoot []byte) ([]byte, error) {
  log.Println("Proving...")

  proof := make([]byte, 0)

  //TODO: This is a mock proof for development reasons
  if len(batchBytes) > 0 {
    proof = append(proof, batchBytes[0])
  }
  proof = append(proof, byte(batchNumber))

  if len(batchPreStateRoot) > 0 {
    proof = append(proof, batchPreStateRoot[0])
  }
  if len(batchPostStateRoot) > 0 {
    proof = append(proof, batchPostStateRoot[0])
  }

  if len(proof) == 0 {
    log.Fatalf("Failed to create proof: %v", "Proof is empty")
    return nil, nil
  }

  log.Println("New Proof is", proof)
  log.Println("Proof Complete!")

  return proof, nil
}

func (p *Prover) SubmitProof(proof []byte, batchNumber uint64) error {
  log.Println("Submitting Proof...")

  p.L1Comms.SubmitProof(proof, int(batchNumber), p.ProverL1Address)
  p.TotalProofs += 1
  p.LastProvedBatch = batchNumber

  log.Println("Proof Submitted for batch %d!", batchNumber)
  return nil
}

func (p *Prover) Start() error {
  log.Println("Starting Prover...")

  for {
    lastBatchConfirmed, err := p.L1Comms.L1Contract.GetLastConfirmedBatch(nil)
    if err != nil {
      log.Fatalf("Failed to get last confirmed batch: %v", err)
    }

    if lastBatchConfirmed.Int64() + 1 == int64(p.LastProvedBatch) {
      log.Println("No new batches to prove...")
      sleepTime := 10 * time.Second
      time.Sleep(sleepTime)
      continue
    }

    batchCount, err := p.L1Comms.L1Contract.GetBatchCount(nil)
    if err != nil {
      log.Fatalf("Failed to get batch count: %v", err)
    }

    batchNumber := uint64(lastBatchConfirmed.Int64()) + 1

    if batchNumber == uint64(batchCount.Int64()) {
      log.Println("No new batches to prove...")
      sleepTime := 10 * time.Second
      time.Sleep(sleepTime)
      continue
    }

    if batchNumber > 0 && batchNumber < uint64(batchCount.Int64()) {
      log.Println("New batch to prove:", batchNumber)
      batchData, batchPreStateRoot, batchPostStateRoot := p.GetBatchProofParams(batchNumber)
      proof, err := p.CreateAggProof(batchData, batchNumber, batchPreStateRoot[:], batchPostStateRoot[:])
      if err != nil {
        log.Fatalf("Failed to create proof: %v", err)
      }
      p.SubmitProof(proof, batchNumber)
      log.Println("Proof submitted for batch", batchNumber)
      time.Sleep(5 * time.Second)
    } else {
      log.Println("Batch number is invalid:", batchNumber, batchCount)
      sleepTime := 10 * time.Second
      time.Sleep(sleepTime)
      continue
    }
  }

  return nil
}

//TODO: This is a mock function for development reasons, keep track of proofs verified onchain & store vales
func (p *Prover) VerifyBatchValid(batchNumber int) error {
  log.Println("Verifying Proof submitted on L1...")

  isConfirmed, err := p.L1Comms.L1Contract.GetBatchConfirmed(nil, big.NewInt(int64(batchNumber)))
  if err != nil {
    log.Fatalf("Failed to get batch confirmed: %v", err)
    return err
  }

  if isConfirmed {
    log.Println("Proof verified on L1!")
    p.TotalProofsVerifiedOnchain += 1

    reward, err := p.L1Comms.L1Contract.GetReward(nil, big.NewInt(int64(batchNumber)))
    if err != nil {
      log.Fatalf("Failed to get reward: %v", err)
      return err
    }

    p.TotalRewards += reward
  } else {
    log.Println("Proof not verified on L1!")
  }

  return nil
}

func main() { os.Exit(mainImp()) }

func mainImp() int {
  l1ContractAddressFile := flag.String("l1-contract-address-file", "", "Path to file containing L1 contract address")
  l1Host := flag.String("l1-host", "http://localhost", "L1 host")
  l1Port := flag.String("l1-port", "8545", "L1 port")
  proverL1AddressFile := flag.String("prover-address-file", "", "Path to file containing prover address")
  proverL1Keystore := flag.String("prover-keystore", "", "Path to prover keystore")
  flag.Parse()

  l1ContractAddress, err := utils.AddressFromFile(*l1ContractAddressFile)
  if err != nil {
    log.Fatalf("Failed to read address from file: %v", err)
  }

  l1Url := *l1Host + ":" + *l1Port
  l1Comms, err := utils.NewL1Comms(l1Url , l1ContractAddress)
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
  }

  proverL1Address, err := utils.AddressFromFile(*proverL1AddressFile)
  if err != nil {
    log.Fatalf("Failed to read address from file: %v", err)
  }

  prover := NewProver(l1Comms, proverL1Address)

  prover.L1Comms.RegisterL2Address(proverL1Address, *proverL1Keystore)


  fatalErrChan := make(chan error, 10)
  err = prover.Start()
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
