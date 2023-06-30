package utils

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	l1bridge "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l1bridge"
	txstorage "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/txstorage"
)

type L1Comms struct {
  RpcUrl string
  L1Client *ethclient.Client

  // L1 Tx Storage
  TxStorageContract *txstorage.Txstorage
  TxStorageContractAddress common.Address

  // L1 Bridge
  BridgeContract *l1bridge.L1bridge
  BridgeContractAddress common.Address
}

func NewL1Comms(rpcUrl string, txStorageContractAddress common.Address, bridgeContractAddress common.Address) (*L1Comms, error) {
  l1Comms := &L1Comms{
    RpcUrl: rpcUrl,
    TxStorageContractAddress: txStorageContractAddress,
    BridgeContractAddress: bridgeContractAddress,
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

  txOpts := MakeTransactOpts(posterAddress, 0)
  _, err := l1Comms.TxStorageContract.StoreGenesisState(&txOpts, genesisHash)
  if err != nil {
    log.Fatalf("Failed to store genesis state on L1: %v", err)
    return err
  }

  log.Println("Stored genesis state on L1")
  return nil
}

func (l1Comms *L1Comms) PostBatch(transactionByteData []byte, id int64, hash [32]byte, posterAddress common.Address) error {
  log.Println("Posting batch to L1...")

  txOpts := MakeTransactOpts(posterAddress, 0)
  _, err := l1Comms.TxStorageContract.StoreBatch(&txOpts, big.NewInt(id), hash, transactionByteData)
  if err != nil {
    log.Fatalf("Failed to post batch to L1: %v", err)
    return err
  }

  log.Println("Posted batch to L1")
  return nil
}

func (l1Comms *L1Comms) SubmitProof(proof []byte, batchNumber int, proverAddress common.Address) error {
  log.Println("Submitting proof to L1...")

  txOpts := MakeTransactOpts(proverAddress, 0)
  _, err := l1Comms.TxStorageContract.SubmitProof(&txOpts, big.NewInt(int64(batchNumber)), proof)
  if err != nil {
    log.Fatalf("Failed to submit proof to L1: %v", err)
    return err
  }

  log.Println("Submitted proof to L1")
  return nil
}

func (l1BridgeComms *L1Comms) BridgeEthToL2(address common.Address, amount uint64) error {
  log.Println("Bridging ", amount, " ETH to L2 for address ", address.Hex())

  txOpts := MakeTransactOpts(address, amount)
  bridgeTx, err := l1BridgeComms.BridgeContract.DepositEth(&txOpts)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}

func (l1BridgeComms *L1Comms) BridgeEthToL1(address common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " ETH to L1 for address ", address.Hex())

  txOpts := MakeTransactOpts(GetSequencer(), 0)
  bridgeTx, err := l1BridgeComms.BridgeContract.WithdrawEth(&txOpts, address, amount)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}
