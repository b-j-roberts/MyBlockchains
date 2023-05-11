package utils

import (
	"log"
	"math/big"
	contracts "naive-l2/contracts/go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type L1Comms struct {
  RpcUrl string
  L1Client *ethclient.Client

  L1Contract *contracts.Contracts
  L1ContractAddress common.Address
}

func NewL1Comms(rpcUrl string, l1ContractAddress common.Address) (*L1Comms, error) {
  rawRpc, err := rpc.Dial(rpcUrl)
  if err != nil {
    log.Fatalf("Failed to connect to RPC: %v", err)
    return nil, err
  }
  l1Client := ethclient.NewClient(rawRpc)

  l1Contract, err := contracts.NewContracts(l1ContractAddress, l1Client)
  if err != nil {
    return nil, err
  }

  return &L1Comms{
    RpcUrl: rpcUrl,
    L1Client: l1Client,
    L1Contract: l1Contract,
    L1ContractAddress: l1ContractAddress,
  }, nil
}

func (l1Comms *L1Comms) RegisterL2Address(posterAddress common.Address, keystoreDir string) error {
  StoreKeyStoreDir(posterAddress, keystoreDir)

  return nil
}

func (l1Comms *L1Comms) L2GenesisOnL1(genesis *core.Genesis, posterAddress common.Address) error {
  var genesisHash [32]byte
  copy(genesisHash[:], genesis.ToBlock().Hash().Bytes())
  _, err := l1Comms.L1Contract.StoreGenesisState(&bind.TransactOpts{
    From: posterAddress,
    Value: big.NewInt(0),
    GasLimit:  3000000, //TODO: Hardcoded
    GasPrice: big.NewInt(200), //TODO: Hardcoded
    Signer: KeystoreSignTx,
  }, genesisHash)
  if err != nil {
    log.Fatalf("Failed to store genesis state on L1: %v", err)
    return err
  }

  log.Println("Stored genesis state on L1")
  return nil
}

func (l1Comms *L1Comms) PostBatch(transactionByteData []byte, id int64, hash [32]byte, posterAddress common.Address) error {
  log.Println("Posting batch to L1...")

  _, err := l1Comms.L1Contract.StoreBatch(&bind.TransactOpts{
    From: posterAddress,
    Value: big.NewInt(0),
    GasLimit:  3000000, //TODO: Hardcoded
    GasPrice: big.NewInt(200), //TODO: Hardcoded
    Signer: KeystoreSignTx,
  }, big.NewInt(id), hash, transactionByteData)
  if err != nil {
    log.Fatalf("Failed to post batch to L1: %v", err)
    return err
  }

  log.Println("Posted batch to L1")
  return nil
}

func (l1Comms *L1Comms) SubmitProof(proof []byte, batchNumber int, proverAddress common.Address) error {
  log.Println("Submitting proof to L1...")

  _, err := l1Comms.L1Contract.SubmitProof(&bind.TransactOpts{
    From: proverAddress,
    Value: big.NewInt(0),
    GasLimit:  3000000, //TODO: Hardcoded
    GasPrice: big.NewInt(200), //TODO: Hardcoded
    Signer: KeystoreSignTx,
  }, big.NewInt(int64(batchNumber)), proof)
  if err != nil {
    log.Fatalf("Failed to submit proof to L1: %v", err)
    return err
  }

  log.Println("Submitted proof to L1")
  return nil
}
