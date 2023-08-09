package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2core "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/core"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

func main() { os.Exit(mainImpl()) }


//TODO: Fix color difference in terminal output between this and geth
// l1 contract address, l1 bridge contract address, l2 bridge contract address, l1 token bridge contract address, l2 token bridge contract address
func mainImpl() int {
  log.Println("Starting sequencer...")
  osHomeDir, err := os.UserHomeDir()
  sequencerConfigFile := flag.String("config", osHomeDir + "/naive-sequencer-data/sequencer.config.json", "sequencer config file")
  sequencerAddress := flag.String("sequencer", "", "Address of the sequencer on L1")
  flag.Parse()


  naiveNode, err := l2core.NewSequencer(*sequencerConfigFile, common.HexToAddress(*sequencerAddress))
  if err != nil {
    utils.Fatalf("Failed to create naive sequencer node: %v", err)
  }
  defer naiveNode.Stop()

  //TODO: Geth metrics being recorded as 0 ( around this commit : a35e6aa7a4d36509c2bbee705aeeb3a2b79a7bb6 )
  l2config.SetupMetrics(naiveNode.L2Node.Config)

  ////TODO: close dbs & stop blockchain defers
  fatalErrChan := make(chan error, 10)

  genesis, err := core.ReadGenesis(naiveNode.L2Node.ChainDb)
  if err != nil {
    fatalErrChan <- err
  }

  //TODO: Dont do this if genesis already exists
  err = naiveNode.Batcher.L1Comms.L2GenesisOnL1(genesis, common.HexToAddress(*sequencerAddress))
  if err != nil {
    fatalErrChan <- err
  }
  naiveNode.Batcher.BatchId = 1

  err = naiveNode.Start()
  if err != nil {
    fatalErrChan <- err
  }

  if naiveNode.L2Node.Config.Metrics.Enabled {
    log.Println("Starting Metrics Server")
    l2utils.StartSystemMetrics()
  }

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
