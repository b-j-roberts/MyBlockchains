package main

import (
	"flag"
	"log"
	"math/big"
	"os"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
	"github.com/ethereum/go-ethereum/common"
)

var Sender common.Address

type TokenType int
const (
  ERC20Type = 1
  ERC721Type = 2
)

func AddL1Token(l1Comms *l2utils.L1Comms, tokenAddr common.Address, tokenType TokenType) error {
  log.Println("Adding allowed token", "tokenAddr", tokenAddr.Hex(), "sender", Sender.Hex())
  l1TransactionOpts, err := l1Comms.CreateL1TransactionOpts(Sender, big.NewInt(0))
  tx, err := l1Comms.L1Contracts.TokenBridgeContract.AddAllowedToken(l1TransactionOpts, tokenAddr, uint8(tokenType))
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return err
  }
  log.Println("Added allowed token", "txHash", tx.Hash().Hex())
  return nil
}

func AddL2Token(l2Comms *l2utils.L2Comms, l1TokenAddr common.Address, l2TokenAddr common.Address) error {
  l2TransactionOpts, err := l2Comms.CreateL2TransactionOpts(Sender, big.NewInt(0))
  tx, err := l2Comms.L2Contracts.L2TokenBridgeContract.AddAllowedToken(l2TransactionOpts, l1TokenAddr, l2TokenAddr)
  if err != nil {
    log.Println("Failed to add allowed token", "error", err)
    return err
  }
  log.Println("Added allowed token", "txHash", tx.Hash().Hex())
  return nil
}

func AddL1Tokens(tokenAddresses l2utils.TokenAddresses, l1Comms *l2utils.L1Comms) error {
  err := AddL1Token(l1Comms, tokenAddresses.Erc20Address, ERC20Type)
  if err != nil {
    return err
  }
  err = AddL1Token(l1Comms, tokenAddresses.StableErc20Address, ERC20Type)
  if err != nil {
    return err
  }
  err = AddL1Token(l1Comms, tokenAddresses.Erc721Address, ERC721Type)
  if err != nil {
    return err
  }
  err = AddL1Token(l1Comms, tokenAddresses.SpecialErc721Address, ERC721Type)
  if err != nil {
    return err
  }
  return nil
}

func AddL2Tokens(tokenAddresses l2utils.TokenAddresses, l2Comms *l2utils.L2Comms) error {
  err := AddL2Token(l2Comms, tokenAddresses.Erc20Address, tokenAddresses.L2Erc20Address)
  if err != nil {
    return err
  }
  err = AddL2Token(l2Comms, tokenAddresses.StableErc20Address, tokenAddresses.L2StableErc20Address)
  if err != nil {
    return err
  }
  err = AddL2Token(l2Comms, tokenAddresses.Erc721Address, tokenAddresses.L2Erc721Address)
  if err != nil {
    return err
  }
  err = AddL2Token(l2Comms, tokenAddresses.SpecialErc721Address, tokenAddresses.L2SpecialErc721Address)
  if err != nil {
    return err
  }
  return nil
}

//TODO: allow-tokens to be single token use and call this multiple times
func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Allow ERC20 token transfer")
  osHomeDir, err := os.UserHomeDir()

  sender := flag.String("sender", "", "sender address")
  //TODO: Config file instead for alll this stuff
  configFile := flag.String("config", osHomeDir + "/naive-sequencer-data/sequencer.config.json", "sequencer config file")
  flag.Parse()

  Sender = common.HexToAddress(*sender)

  config, err := l2config.LoadNodeBaseConfig(*configFile)
  if err != nil {
    log.Println("Failed to load config", "error", err)
    return 1
  }

  tokenAddresses, err := l2utils.LoadTokenAddresses(config.Contracts)
  if err != nil {
    log.Println("Failed to load token addresses", "error", err)
    return 1
  }

  l1Comms, err := l2utils.NewL1Comms(config.L1URL, config.Contracts, big.NewInt(int64(config.L1ChainID)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Println("Failed to create L1 comms", "error", err)
    return 1
  }

  l2Comms, err := l2utils.NewL2Comms(config, l2utils.GetDefaultL2TransactionConfig())
  if err != nil {
    log.Println("Failed to create L2 comms", "error", err)
    return 1
  }

  err = AddL1Tokens(tokenAddresses, l1Comms)
  if err != nil {
    log.Println("Failed to add L1 tokens", "error", err)
    return 1
  }

  err = AddL2Tokens(tokenAddresses, l2Comms)
  if err != nil {
    log.Println("Failed to add L2 tokens", "error", err)
    return 1
  }

  return 0
}
