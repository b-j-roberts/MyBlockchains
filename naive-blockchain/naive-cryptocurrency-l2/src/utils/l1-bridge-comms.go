package utils

import (
	"log"
	"math/big"

	l1bridge "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l1bridge"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type L1BridgeComms struct {
  RpcUrl string
  L1Client *ethclient.Client
  L1BridgeContract *l1bridge.L1bridge
  L1BridgeContractAddress common.Address
}

func NewL1BridgeComms(rpcUrl string, l1BridgeContractAddress common.Address) (*L1BridgeComms, error) {
  log.Println("Initializing L1BridgeComms w/ RPC URL: ", rpcUrl, " and L1BridgeContractAddress: ", l1BridgeContractAddress)

  rawRpc, err := rpc.Dial(rpcUrl)
  if err != nil {
    log.Println("Failed to connect to RPC", err)
    return nil, err
  }
  l1Client := ethclient.NewClient(rawRpc)

  l1BridgeContract, err := l1bridge.NewL1bridge(l1BridgeContractAddress, l1Client)
  if err != nil {
    log.Println("Failed to instantiate L1BridgeContract", err)
    return nil, err
  }

  return &L1BridgeComms{
    RpcUrl: rpcUrl,
    L1Client: l1Client,
    L1BridgeContract: l1BridgeContract,
    L1BridgeContractAddress: l1BridgeContractAddress,
  }, nil
}

func (l1BridgeComms *L1BridgeComms) RegisterL1Address(address common.Address, keystoreDir string) error {
  StoreKeyStoreDir(address, keystoreDir)

  return nil
}

func (l1BridgeComms *L1BridgeComms) BridgeEthToL2(address common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " ETH to L2 for address ", address.Hex())

  bridgeTx, err := l1BridgeComms.L1BridgeContract.DepositEth(&bind.TransactOpts{
    From: address,
    Value: amount,
    GasLimit:  3000000, //TODO: Hardcoded
    GasPrice: big.NewInt(200), //TODO: Hardcoded
    Signer: KeystoreSignTx,
  })
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())

  return nil
}

func (l1BridgeComms *L1BridgeComms) BridgeEthToL1(address common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " ETH to L1 for address ", address.Hex())

  bridgeTx, err := l1BridgeComms.L1BridgeContract.WithdrawEth(&bind.TransactOpts{
    From: GetSequencer(), // Sequencer L1 addr
    GasLimit:  3000000, //TODO: Hardcoded
    GasPrice: big.NewInt(200), //TODO: Hardcoded
    Signer: KeystoreSignTx,
  }, address, amount)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())

  return nil
}

