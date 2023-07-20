package utils

import (
	"fmt"
	"log"
	"math/big"

	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l2bridge"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type L2Comms struct {
  // L2 Bridge
  L2BridgeContract *l2bridge.L2bridge
  BridgeContractAddress common.Address
  L2Backend *ethclient.Client
}

func NewL2Comms(bridgeContractAddress common.Address, l2Backend *ethclient.Client) (*L2Comms, error) {
  l2Comms := &L2Comms{
    BridgeContractAddress: bridgeContractAddress,
    L2Backend: l2Backend,
  }

  var err error
  l2Comms.L2BridgeContract, err = l2bridge.NewL2bridge(l2Comms.BridgeContractAddress, l2Comms.L2Backend)
  if err != nil {
    return nil, err
  }

  return l2Comms, nil
}

func (l2Comms *L2Comms) CreateL2TransactionOpts(fromAddress common.Address, value *big.Int) (*bind.TransactOpts, error) {
  transactOpts, err := CreateTransactOpts(accounts.Account{Address: fromAddress}, big.NewInt(515)) //TODO: Hardcoded
  if err != nil {
    return nil, err
  }
  transactOpts.GasLimit = 3000000
  transactOpts.GasPrice = big.NewInt(200) //TODO: Hardcoded
  //transactOpts.GasLimit = l2Comms.L2TransactionConfig.GasLimit
  //transactOpts.GasPrice = l2Comms.L2TransactionConfig.GasPrice
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
