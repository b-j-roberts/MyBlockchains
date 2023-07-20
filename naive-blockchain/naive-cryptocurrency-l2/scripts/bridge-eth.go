// Basic script to bridge eth from an account

package main

import (
	"context"
	"flag"
	"math/big"
	"os"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Info("Starting bridge-eth")

  address := flag.String("address", "", "address to bridge with")
  value := flag.Uint64("value", 0, "value to send")
  rpcURL := flag.String("rpc", "http://localhost:8545", "rpc address")
  bridgeAddress := flag.String("bridgeAddress", "", "bridge address")
  keystore := flag.String("keystore", "", "keystore directory")
  chainId := flag.Uint64("chainId", 505, "l1 chain id")
  toL1 := flag.Bool("to-l1", false, "send to l1")
  flag.Parse()

  if *value == 0 {
    flag.Usage()
    return 1
  }

  if *address == "" {
    flag.Usage()
    return 1
  }

  if !*toL1 {
    log.Info("Sending eth to l2")
    l1BridgeComms, err := l2utils.NewL1Comms(*rpcURL, common.HexToAddress("0x0"), common.HexToAddress(*bridgeAddress), big.NewInt(int64(*chainId)), l2utils.L1TransactionConfig{
      GasLimit: 3000000,
      GasPrice: big.NewInt(200),
    })
    if err != nil {
      panic(err)
    }

    l2utils.RegisterAccount(common.HexToAddress(*address), *keystore)

    err = l1BridgeComms.BridgeEthToL2(common.HexToAddress(*address), *value)
    if err != nil {
      panic(err)
    }
  } else {
    l2utils.RegisterAccount(common.HexToAddress(*address), *keystore)

    log.Info("Sending eth to l1")
    ipcFile := *keystore + "/../naive-sequencer.ipc"
    rpcIPC, err := rpc.DialIPC(context.Background(), ipcFile)
    if err != nil {
      panic(err)
    }
    backend := ethclient.NewClient(rpcIPC)

    l2BridgeComms, err := l2utils.NewL2Comms(common.HexToAddress(*bridgeAddress), backend)
    if err != nil {
      panic(err)
    }
    err = l2BridgeComms.BridgeEthToL1(common.HexToAddress(*address), big.NewInt(int64(*value)))
    if err != nil {
      panic(err)
    }

    log.Info("Success")

  }

  log.Info("Success")

  return 0
}
