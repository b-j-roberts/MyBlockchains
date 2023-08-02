package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/cmd/utils"

	l2core "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/core"
)

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Starting rpc...")

  osHomeDir, err := os.UserHomeDir()
  rpcConfigFile := flag.String("config", osHomeDir + "/naive-rpc-data/rpc.config.json", "Path to rpc config file")
  flag.Parse()

  naiveNode, err := l2core.NewNode(*rpcConfigFile)
  if err != nil {
    utils.Fatalf("Failed to create naive rpc node: %v", err)
  }

  ////TODO: close dbs & stop blockchain defers
  fatalErrChan := make(chan error, 10)

  err = naiveNode.Start()
  if err != nil {
    fatalErrChan <- err
  }
  log.Println("Naive Rpc Node started", naiveNode)

  sigint := make(chan os.Signal, 1)
  signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
  
  exitCode := 0
  select {
  case err := <-fatalErrChan:
    log.Println("shutting down due to fatal error", "err", err)
    defer log.Println("shut down due to fatal error", "err", err)
    exitCode = 1
  case <-sigint:
    log.Println("shutting down because of sigint")
  }
  
  // cause future ctrl+c's to panic
  close(sigint)

  // node stop&wait

  return exitCode
}
