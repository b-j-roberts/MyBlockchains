package utils

import (
	"fmt"
	"log"
	"math/big"

	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l2bridge"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l2tokenbridge"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type L2TransactionConfig struct {
  GasLimit uint64
  GasPrice *big.Int
}

type L2Comms struct {
  // L2 Bridge
  L2BridgeContract *l2bridge.L2bridge
  BridgeContractAddress common.Address

  // L2 Token Bridge
  L2TokenBridgeContract *l2tokenbridge.L2tokenbridge
  TokenBridgeContractAddress common.Address

  L2Backend *ethclient.Client
  L2ChainId *big.Int
  L2TransactionConfig L2TransactionConfig
}

func GetDefaultL2TransactionConfig() L2TransactionConfig {
  return L2TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  }
}

func NewL2Comms(bridgeContractAddress common.Address, tokenBridgeContractAddress common.Address, l2ChainId *big.Int, l2Backend *ethclient.Client, l2TransactionConfig L2TransactionConfig) (*L2Comms, error) {
  l2Comms := &L2Comms{
    BridgeContractAddress: bridgeContractAddress,
    TokenBridgeContractAddress: tokenBridgeContractAddress,
    L2Backend: l2Backend,
    L2ChainId: l2ChainId,
    L2TransactionConfig: l2TransactionConfig,
  }

  var err error
  l2Comms.L2BridgeContract, err = l2bridge.NewL2bridge(l2Comms.BridgeContractAddress, l2Comms.L2Backend)
  if err != nil {
    return nil, err
  }

  l2Comms.L2TokenBridgeContract, err = l2tokenbridge.NewL2tokenbridge(l2Comms.TokenBridgeContractAddress, l2Comms.L2Backend)
  if err != nil {
    return nil, err
  }

  return l2Comms, nil
}

func (l2Comms *L2Comms) CreateL2TransactionOpts(fromAddress common.Address, value *big.Int) (*bind.TransactOpts, error) {
  transactOpts, err := CreateTransactOpts(accounts.Account{Address: fromAddress}, l2Comms.L2ChainId)
  if err != nil {
    return nil, err
  }
  transactOpts.GasLimit = l2Comms.L2TransactionConfig.GasLimit
  transactOpts.GasPrice = l2Comms.L2TransactionConfig.GasPrice
  transactOpts.Value = value
  
  return transactOpts, nil
}

func (l2Comms *L2Comms) BridgeEthToL1(address common.Address, amount *big.Int) error {
  log.Println("BridgeEthToL1 called with address:", address.Hex(), "amount:", amount.String())

  transactOpts, err := l2Comms.CreateL2TransactionOpts(address, amount)
  if err != nil {
    log.Println("BridgeEthToL1 CreateTransactOpts error:", err)
    return err
  }

  log.Println("BridgeEthToL1 transactOpts created w/ value:", transactOpts.From.Hex())
  tx, err := l2Comms.L2BridgeContract.WithdrawEth(transactOpts)
  if err != nil {
    log.Println("BridgeEthToL1 WithdrawEth error:", err)
    return err
  }

  fmt.Println("BridgeEthToL1 tx sent:", tx.Hash().Hex())
  return nil
}

func (l2Comms *L2Comms) BridgeTokenToL1(tokenAddress common.Address, address common.Address, amount *big.Int) error {
  log.Println("BridgeTokenToL1 called with address:", address.Hex(), "amount:", amount.String())

  transactOpts, err := l2Comms.CreateL2TransactionOpts(address, big.NewInt(0))
  if err != nil {
    log.Println("BridgeTokenToL1 CreateTransactOpts error:", err)
    return err
  }

  log.Println("BridgeTokenToL1 transactOpts created w/ value:", transactOpts.From.Hex())
  tx, err := l2Comms.L2TokenBridgeContract.WithdrawTokens(transactOpts, tokenAddress, amount)
  if err != nil {
    log.Println("BridgeTokenToL1 WithdrawToken error:", err)
    return err
  }

  fmt.Println("BridgeTokenToL1 tx sent:", tx.Hash().Hex())
  return nil
}
