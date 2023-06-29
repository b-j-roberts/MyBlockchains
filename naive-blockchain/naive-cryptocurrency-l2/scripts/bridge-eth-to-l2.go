// Basic script to bridge eth from an account

package main

import (
	"flag"
	"math/big"
	"os"

	naive_utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Info("Starting bridge-eth")

  address := flag.String("address", "", "address to send to")
  value := flag.Uint64("value", 0, "value to send")
  rpc := flag.String("rpc", "http://localhost:8545", "rpc address")
  bridgeAddress := flag.String("bridgeAddress", "", "bridge address")
  keystore := flag.String("keystore", "", "keystore directory")
  flag.Parse()

  if *value == 0 {
    flag.Usage()
    return 1
  }

  if *address == "" {
    flag.Usage()
    return 1
  }

  l1BridgeComms, err := naive_utils.NewL1BridgeComms(*rpc, common.HexToAddress(*bridgeAddress))
  if err != nil {
    panic(err)
  }

  err = l1BridgeComms.RegisterL1Address(common.HexToAddress(*address), *keystore)
  if err != nil {
    panic(err)
  }
  
  err = l1BridgeComms.BridgeEthToL2(common.HexToAddress(*address), big.NewInt(int64(*value)))
  if err != nil {
    panic(err)
  }

  log.Info("Success")

  return 0
}
