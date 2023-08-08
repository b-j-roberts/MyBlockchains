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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
)

type L2TransactionConfig struct {
  GasLimit uint64
  GasPrice *big.Int
}

type L2ContractAddressConfig struct {
  BridgeContractAddress common.Address
  TokenBridgeContractAddress common.Address
}

type L2Contracts struct {
  L2BridgeContract *l2bridge.L2bridge
  L2TokenBridgeContract *l2tokenbridge.L2tokenbridge
}

func CreateL2ContractAddressConfig(contractsAddressDir string) L2ContractAddressConfig {
  //TODO: Load into memeory and cache
  bridgeContractAddress, err := ReadContractAddressFromFile(contractsAddressDir + "/l2-bridge-address.txt")
  if err != nil {
    log.Fatal("CreateL2ContractAddressConfig ReadContractAddressFromFile error:", err)
  }

  tokenBridgeContractAddress, err := ReadContractAddressFromFile(contractsAddressDir + "/l2-token-bridge-address.txt")
  if err != nil {
    log.Fatal("CreateL2ContractAddressConfig ReadContractAddressFromFile error:", err)
  }

  return L2ContractAddressConfig{
    BridgeContractAddress: bridgeContractAddress,
    TokenBridgeContractAddress: tokenBridgeContractAddress,
  }
}

func CreateL2Contracts(client *ethclient.Client, l2ContractAddressConfig L2ContractAddressConfig) L2Contracts {
  l2BridgeContract, err := l2bridge.NewL2bridge(l2ContractAddressConfig.BridgeContractAddress, client)
  if err != nil {
    log.Fatal("CreateL2Contracts NewL2bridge error:", err)
  }

  l2TokenBridgeContract, err := l2tokenbridge.NewL2tokenbridge(l2ContractAddressConfig.TokenBridgeContractAddress, client)
  if err != nil {
    log.Fatal("CreateL2Contracts NewL2tokenbridge error:", err)
  }

  return L2Contracts{
    L2BridgeContract: l2BridgeContract,
    L2TokenBridgeContract: l2TokenBridgeContract,
  }
}

type L2Comms struct {
  L2Client *ethclient.Client
  L2ChainId *big.Int
  L2TransactionConfig L2TransactionConfig
  L2ContractAddressConfig  L2ContractAddressConfig
  L2Contracts L2Contracts
}

func GetDefaultL2TransactionConfig() L2TransactionConfig {
  return L2TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  }
}

func NewL2Comms(nodeConfig *l2config.NodeBaseConfig, l2TransactionConfig L2TransactionConfig) (*L2Comms, error) {
  //rawIPC, err := rpc.DialIPC(context.Background(), ipcFile)
  //if err != nil {
  //  return nil, err
  //}
  l2Url := fmt.Sprintf("http://%s:%d", nodeConfig.Host, nodeConfig.Port)
  rpc, err := rpc.Dial(l2Url)
  if err != nil {
    return nil, err
  }
  client := ethclient.NewClient(rpc)

  l2Comms := &L2Comms{
    L2Client: client,
    L2ChainId: big.NewInt(int64(nodeConfig.L2ChainID)),
    L2TransactionConfig: l2TransactionConfig,
    L2ContractAddressConfig: CreateL2ContractAddressConfig(nodeConfig.Contracts),
  }
  l2Comms.L2Contracts = CreateL2Contracts(l2Comms.L2Client, l2Comms.L2ContractAddressConfig)

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
  tx, err := l2Comms.L2Contracts.L2BridgeContract.WithdrawEth(transactOpts)
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
  tx, err := l2Comms.L2Contracts.L2TokenBridgeContract.WithdrawTokens(transactOpts, tokenAddress, amount)
  if err != nil {
    log.Println("BridgeTokenToL1 WithdrawToken error:", err)
    return err
  }

  fmt.Println("BridgeTokenToL1 tx sent:", tx.Hash().Hex())
  return nil
}

func UnpackEthDeposited(receiptLog types.Log) (nonce *big.Int, addr common.Address, amount *big.Int, err error) {
    data := receiptLog.Data
    if len(data) < 10 {
        err = fmt.Errorf("invalid data")
        return
    }

    offset := 12
    nonce = new(big.Int).SetBytes(data[:32])
    addr = common.BytesToAddress(data[32:52+offset])
    amount = new(big.Int).SetBytes(data[52+offset:84+offset])

    return
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
  addr = common.BytesToAddress(data[32:64])
  tokenAddr = common.BytesToAddress(data[64:96])
  amount = new(big.Int).SetBytes(data[96:128])
             
  return
}
