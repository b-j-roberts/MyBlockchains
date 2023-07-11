package utils

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	l1bridge "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l1bridge"
	txstorage "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/txstorage"
)

type L1TransactionConfig struct {
  GasLimit uint64
  GasPrice *big.Int
}

type L1Comms struct {
  RpcUrl string
  L1Client *ethclient.Client
  L1ChainID *big.Int
  L1TransactionConfig L1TransactionConfig

  // L1 Tx Storage
  TxStorageContract *txstorage.Txstorage
  TxStorageContractAddress common.Address

  // L1 Bridge
  BridgeContract *l1bridge.L1bridge
  BridgeContractAddress common.Address
}

func (l1Comms *L1Comms) CreateL1TransactionOpts(fromAddress common.Address, value *big.Int) (*bind.TransactOpts, error) {
  transactOpts, err := CreateTransactOpts(accounts.Account{Address: fromAddress}, l1Comms.L1ChainID)
  if err != nil {
    return nil, err
  }
  transactOpts.GasLimit = l1Comms.L1TransactionConfig.GasLimit
  transactOpts.GasPrice = l1Comms.L1TransactionConfig.GasPrice
  transactOpts.Value = value

  return transactOpts, nil
}

func NewL1Comms(rpcUrl string, txStorageContractAddress common.Address, bridgeContractAddress common.Address, chainID *big.Int, l1TransactionConfig L1TransactionConfig) (*L1Comms, error) {
  l1Comms := &L1Comms{
    RpcUrl: rpcUrl,
    TxStorageContractAddress: txStorageContractAddress,
    BridgeContractAddress: bridgeContractAddress,
    L1ChainID: chainID,
    L1TransactionConfig: l1TransactionConfig,
  }

  // Connect to L1
  rawRpc, err := rpc.Dial(l1Comms.RpcUrl)
  if err != nil {
    return nil, err
  }
  l1Comms.L1Client = ethclient.NewClient(rawRpc)

  l1Comms.TxStorageContract, err = txstorage.NewTxstorage(l1Comms.TxStorageContractAddress, l1Comms.L1Client)
  if err != nil {
    return nil, err
  }

  l1Comms.BridgeContract, err = l1bridge.NewL1bridge(l1Comms.BridgeContractAddress, l1Comms.L1Client)
  if err != nil {
    return nil, err
  }

  return l1Comms, nil
}

func (l1Comms *L1Comms) L2GenesisOnL1(genesis *core.Genesis, posterAddress common.Address) error {
  var genesisHash [32]byte
  copy(genesisHash[:], genesis.ToBlock().Hash().Bytes())

  transactOpts, err := l1Comms.CreateL1TransactionOpts(posterAddress, big.NewInt(0))
  if err != nil {
    log.Fatalf("Failed to create L1 transaction opts: %v", err)
    return err
  }

  _, err = l1Comms.TxStorageContract.StoreGenesisState(transactOpts, genesisHash)
  if err != nil {
    log.Fatalf("Failed to store genesis state on L1: %v", err)
    return err
  }

  log.Println("Stored genesis state on L1")
  return nil
}

func (l1Comms *L1Comms) PostBatch(transactionByteData []byte, id int64, hash [32]byte, posterAddress common.Address) error {
  log.Println("Posting batch to L1...")

  transactOpts, err := l1Comms.CreateL1TransactionOpts(posterAddress, big.NewInt(0))
  if err != nil {
    log.Fatalf("Failed to create L1 transaction opts: %v", err)
    return err
  }

  _, err = l1Comms.TxStorageContract.StoreBatch(transactOpts, big.NewInt(id), hash, transactionByteData)
  if err != nil {
    log.Fatalf("Failed to post batch to L1: %v", err)
    return err
  }

  log.Println("Posted batch to L1")
  return nil
}

func (l1Comms *L1Comms) SubmitProof(proof []byte, batchNumber int, proverAddress common.Address) error {
  log.Println("Submitting proof to L1...")

  transactOpts, err := l1Comms.CreateL1TransactionOpts(proverAddress, big.NewInt(0))
  if err != nil {
    log.Fatalf("Failed to create L1 transaction opts: %v", err)
    return err
  }

  _, err = l1Comms.TxStorageContract.SubmitProof(transactOpts, big.NewInt(int64(batchNumber)), proof)
  if err != nil {
    log.Fatalf("Failed to submit proof to L1: %v", err)
    return err
  }

  log.Println("Submitted proof to L1")
  return nil
}

func (l1BridgeComms *L1Comms) BridgeEthToL2(address common.Address, amount uint64) error {
  log.Println("Bridging ", amount, " ETH to L2 for address ", address.Hex())

  transactOpts, err := l1BridgeComms.CreateL1TransactionOpts(address, big.NewInt(int64(amount)))
  if err != nil {
    log.Println("Failed to create L1 transaction opts", err)
    return err
  }

  bridgeTx, err := l1BridgeComms.BridgeContract.DepositEth(transactOpts)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}

func (l1BridgeComms *L1Comms) BridgeEthToL1(address common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " ETH to L1 for address ", address.Hex())

  transactOpts, err := l1BridgeComms.CreateL1TransactionOpts(GetSequencer(), big.NewInt(0))
  if err != nil {
    log.Println("Failed to create L1 transaction opts", err)
    return err
  }

  bridgeTx, err := l1BridgeComms.BridgeContract.WithdrawEth(transactOpts, address, amount)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}
