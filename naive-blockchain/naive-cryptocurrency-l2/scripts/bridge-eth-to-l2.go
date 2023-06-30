// Basic script to bridge eth from an account

package main

import (
	"flag"
	"os"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
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

  l1BridgeComms, err := l2utils.NewL1Comms(*rpc, common.HexToAddress("0x0"), common.HexToAddress(*bridgeAddress))
  if err != nil {
    panic(err)
  }

  l2utils.RegisterAccount(common.HexToAddress(*address), *keystore)

  err = l1BridgeComms.BridgeEthToL2(common.HexToAddress(*address), *value)
  if err != nil {
    panic(err)
  }

  log.Info("Success")

  return 0
}
