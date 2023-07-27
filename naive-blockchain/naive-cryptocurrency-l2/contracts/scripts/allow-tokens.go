package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"os"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Allow ERC20 token transfer")

  sender := flag.String("sender", "", "sender address")
  erc20Address := flag.String("erc20", "", "ERC20 token address")
  l2erc20Address := flag.String("l2erc20", "", "L2 ERC20 token address")
  stableErc20Address := flag.String("stableErc20", "", "stable ERC20 token address")
  l2StableErc20Address := flag.String("l2stableErc20", "", "L2 stable ERC20 token address")
  erc721Address := flag.String("erc721", "", "ERC721 token address")
  l2erc721Address := flag.String("l2erc721", "", "L2 ERC721 token address")
  specialErc721Address := flag.String("specialErc721", "", "special ERC721 token address")
  l2SpecialErc721Address := flag.String("l2specialErc721", "", "L2 special ERC721 token address")
  tokenBridgeAddress := flag.String("tokenBridge", "", "token bridge address")
  l2TokenBridgeAddress := flag.String("l2tokenBridge", "", "L2 token bridge address")
  rpcUrl := flag.String("rpc", "http://localhost:8545", "RPC URL")
  keystore := flag.String("keystore", "", "keystore directory")
  chainId := flag.Int("chainId", 505, "chain ID")
  l2ChainId := flag.Int("l2chainId", 515, "L2 chain ID")
  flag.Parse()

  log.Println("All token addresses : ", "erc20", *erc20Address, "l2erc20", *l2erc20Address, "stableErc20", *stableErc20Address, "l2StableErc20", *l2StableErc20Address, "erc721", *erc721Address, "l2erc721", *l2erc721Address, "specialErc721", *specialErc721Address, "l2SpecialErc721", *l2SpecialErc721Address)

  l2utils.RegisterAccount(common.HexToAddress(*sender), *keystore)

  l1Comms, err := l2utils.NewL1Comms(*rpcUrl, common.HexToAddress("0x0"), common.HexToAddress("0x0"), common.HexToAddress(*tokenBridgeAddress), big.NewInt(int64(*chainId)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Println("Failed to create L1 comms", "error", err)
    return 1
  }

  ipcFile := *keystore + "/../naive-sequencer.ipc"
  rpcIPC, err := rpc.DialIPC(context.Background(), ipcFile)
  if err != nil {
    log.Println("Failed to connect to L2 RPC", "error", err)
    return 1
  }
  backend := ethclient.NewClient(rpcIPC)

  l2Comms, err := l2utils.NewL2Comms(common.HexToAddress("0x0"), common.HexToAddress(*l2TokenBridgeAddress), big.NewInt(int64(*l2ChainId)), backend, l2utils.GetDefaultL2TransactionConfig())
  if err != nil {
    log.Println("Failed to create L2 comms", "error", err)
    return 1
  }

  l1TransactionOpts, err := l1Comms.CreateL1TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err := l1Comms.TokenBridgeContract.AddAllowedToken(l1TransactionOpts, common.HexToAddress(*erc20Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }
  log.Println("Added allowed token", "txHash", tx.Hash().Hex())

  l2TransactionOpts, err := l2Comms.CreateL2TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l2Comms.L2TokenBridgeContract.AddAllowedToken(l2TransactionOpts, common.HexToAddress(*erc20Address), common.HexToAddress(*l2erc20Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }
  log.Println("Added allowed token", "txHash", tx.Hash().Hex())

  l1TransactionOpts, err = l1Comms.CreateL1TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l1Comms.TokenBridgeContract.AddAllowedToken(l1TransactionOpts, common.HexToAddress(*stableErc20Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }
  log.Println("Added allowed token", "txHash", tx.Hash().Hex())

  l2TransactionOpts, err = l2Comms.CreateL2TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l2Comms.L2TokenBridgeContract.AddAllowedToken(l2TransactionOpts, common.HexToAddress(*stableErc20Address), common.HexToAddress(*l2StableErc20Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }
  log.Println("Added allowed token", "txHash", tx.Hash().Hex())

  l1TransactionOpts, err = l1Comms.CreateL1TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l1Comms.TokenBridgeContract.AddAllowedToken(l1TransactionOpts, common.HexToAddress(*erc721Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }

  l2TransactionOpts, err = l2Comms.CreateL2TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l2Comms.L2TokenBridgeContract.AddAllowedToken(l2TransactionOpts, common.HexToAddress(*erc721Address), common.HexToAddress(*l2erc721Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }

  l1TransactionOpts, err = l1Comms.CreateL1TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l1Comms.TokenBridgeContract.AddAllowedToken(l1TransactionOpts, common.HexToAddress(*specialErc721Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }

  l2TransactionOpts, err = l2Comms.CreateL2TransactionOpts(common.HexToAddress(*sender), big.NewInt(0))
  tx, err = l2Comms.L2TokenBridgeContract.AddAllowedToken(l2TransactionOpts, common.HexToAddress(*specialErc721Address), common.HexToAddress(*l2SpecialErc721Address))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return 1
  }

  return 0
}
