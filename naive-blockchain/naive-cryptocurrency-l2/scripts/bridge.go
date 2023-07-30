// Basic script to bridge eth from an account

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
  log.Println("Starting bridge-eth")

  address := flag.String("address", "", "address to bridge with")
  value := flag.Uint64("value", 0, "value to send")
  rpcURL := flag.String("rpc", "http://localhost:8545", "rpc address")
  bridgeAddress := flag.String("bridgeAddress", "", "bridge address")
  keystore := flag.String("keystore", "", "keystore directory")
  chainId := flag.Uint64("chainId", 505, "l1 chain id")
  l2chainId := flag.Uint64("l2chainId", 515, "l1 chain id")
  toL1 := flag.Bool("to-l1", false, "send to l1")
  isERC := flag.Bool("is-erc", false, "is erc20")
  tokenAddress := flag.String("token", "", "token address")
  ipcFile := flag.String("ipc", "", "File for sequencer ipc")
  flag.Parse()

  if *address == "" {
    flag.Usage()
    return 1
  }

  if *isERC {
    if !*toL1 {
      log.Println("Sending erc20 to l2")
      
      log.Println("Bridging with these values:", "bridge", *bridgeAddress, "address", *address, "tokenAddress", *tokenAddress, "value", *value, "chainId", *chainId, "l2chainId", *l2chainId)
      l1Comms, err := l2utils.NewL1Comms(*rpcURL, common.HexToAddress("0x0"), common.HexToAddress("0x0"), common.HexToAddress(*bridgeAddress), big.NewInt(int64(*chainId)), l2utils.L1TransactionConfig{
        GasLimit: 3000000, 
        GasPrice: big.NewInt(200),
      })
      if err != nil {
        panic(err)
      }

      l2utils.RegisterAccount(common.HexToAddress(*address), *keystore)

      err = l1Comms.BridgeTokenToL2(common.HexToAddress(*tokenAddress), common.HexToAddress(*address), big.NewInt(int64(*value)))
      if err != nil {
        panic(err)
      }
    } else {
      log.Println("Sending erc20 to l1")

      l2utils.RegisterAccount(common.HexToAddress(*address), *keystore)

      rpcIPC, err := rpc.DialIPC(context.Background(), *ipcFile)
      if err != nil {
        panic(err)
      }
      backend := ethclient.NewClient(rpcIPC)

      l2Comms, err := l2utils.NewL2Comms(common.HexToAddress("0x0"), common.HexToAddress(*bridgeAddress), big.NewInt(int64(*l2chainId)), backend, l2utils.GetDefaultL2TransactionConfig())
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
      l1BridgeComms, err := l2utils.NewL1Comms(*rpcURL, common.HexToAddress("0x0"), common.HexToAddress(*bridgeAddress), common.HexToAddress("0x0"), big.NewInt(int64(*chainId)), l2utils.L1TransactionConfig{
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

      log.Println("Success")
    } else {
      l2utils.RegisterAccount(common.HexToAddress(*address), *keystore)

      log.Println("Sending eth to l1")
      rpcIPC, err := rpc.DialIPC(context.Background(), *ipcFile)
      if err != nil {
        panic(err)
      }
      backend := ethclient.NewClient(rpcIPC)

      log.Println("Creating l2 comms with values :", common.HexToAddress(*bridgeAddress), common.HexToAddress("0x0"), big.NewInt(int64(*l2chainId)))
      l2BridgeComms, err := l2utils.NewL2Comms(common.HexToAddress(*bridgeAddress), common.HexToAddress("0x0"), big.NewInt(int64(*l2chainId)), backend, l2utils.GetDefaultL2TransactionConfig())
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
