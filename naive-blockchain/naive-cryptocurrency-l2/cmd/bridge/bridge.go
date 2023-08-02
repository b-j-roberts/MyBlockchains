// Basic script to bridge eth from an account

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

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Starting bridge-eth")

  osHomeDir := os.Getenv("HOME")

  address := flag.String("address", "", "address to bridge with")
  value := flag.Uint64("value", 0, "value to send")
  toL1 := flag.Bool("to-l1", false, "send to l1")
  isERC := flag.Bool("is-erc", false, "is erc20")
  tokenAddress := flag.String("token", "", "token address")
  configFile := flag.String("config", osHomeDir + "/naive-sequencer-data/sequencer.config.json", "config file")
  flag.Parse()

  if *address == "" {
    flag.Usage()
    return 1
  }

  config, err := l2config.LoadNodeBaseConfig(*configFile)
  if err != nil {
    panic(err)
  }

  if *isERC {
    if !*toL1 {
      log.Println("Sending erc20 to l2")
      
      l1Comms, err := l2utils.NewL1Comms(config.L1URL, config.Contracts, big.NewInt(int64(config.L1ChainID)), l2utils.L1TransactionConfig{
        GasLimit: 3000000, 
        GasPrice: big.NewInt(200),
      })
      if err != nil {
        panic(err)
      }

      err = l1Comms.BridgeTokenToL2(common.HexToAddress(*tokenAddress), common.HexToAddress(*address), big.NewInt(int64(*value)))
      if err != nil {
        panic(err)
      }
    } else {
      log.Println("Sending erc20 to l1")

      l2Comms, err := l2utils.NewL2Comms(config.DataDir + "/naive-sequencer.ipc", config.Contracts, big.NewInt(int64(config.L2ChainID)), l2utils.GetDefaultL2TransactionConfig())
      if err != nil {
        panic(err)
      }

      log.Println("Bridging tokens with values :", common.HexToAddress(*tokenAddress), common.HexToAddress(*address), big.NewInt(int64(*value)))
      err = l2Comms.BridgeTokenToL1(common.HexToAddress(*tokenAddress), common.HexToAddress(*address), big.NewInt(int64(*value)))
      if err != nil {
        panic(err)
      }


      log.Println("Success")
    }
  } else {
    if !*toL1 {
      log.Println("Sending eth to l2")
      l1BridgeComms, err := l2utils.NewL1Comms(config.L1URL, config.Contracts, big.NewInt(int64(config.L1ChainID)), l2utils.L1TransactionConfig{
        GasLimit: 3000000,
        GasPrice: big.NewInt(200),
      })
      if err != nil {
        panic(err)
      }

      err = l1BridgeComms.BridgeEthToL2(common.HexToAddress(*address), *value)
      if err != nil {
        panic(err)
      }

      log.Println("Success")
    } else {
      l2BridgeComms, err := l2utils.NewL2Comms(config.DataDir + "/naive-sequencer.ipc", config.Contracts, big.NewInt(int64(config.L2ChainID)), l2utils.GetDefaultL2TransactionConfig())
      if err != nil {
        panic(err)
      }
      log.Println("Bridging eth with values :", common.HexToAddress(*address), big.NewInt(int64(*value)))
      err = l2BridgeComms.BridgeEthToL1(common.HexToAddress(*address), big.NewInt(int64(*value)))
      if err != nil {
        panic(err)
      }

      log.Println("Success")

    }
  }

  log.Println("Success")

  return 0
}
