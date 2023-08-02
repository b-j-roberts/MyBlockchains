package main

import (
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2core "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/core"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

func main() { os.Exit(mainImpl()) }

//TODO: One over for error handling & logging
func mainImpl() int {
  osHomeDir, err := os.UserHomeDir()

  configFile := flag.String("config", osHomeDir + "/naive-sequencer-data/sequencer.config.json", "Config file")
  proverL1Address := flag.String("prover-address", "", "prover address")
  flag.Parse()

  config, err := l2config.LoadNodeBaseConfig(*configFile)
  if err != nil {
    log.Fatalf("Failed to load config: %v", err)
    return 1
  }

  l1Comms, err := l2utils.NewL1Comms(config.L1URL, config.Contracts, big.NewInt(int64(config.L1ChainID)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
    return 1
  }

  l2core.SetupMetrics()

  prover := l2core.NewProver(l1Comms, common.HexToAddress(*proverL1Address))

  err = prover.Start()
  if err != nil {
    log.Fatalf("Failed to start prover: %v", err)
    return 1
  }

  log.Println("Starting Prometheus metrics server on port 6162...")
  http.Handle("/metrics", promhttp.Handler())
  http.ListenAndServe(":6162", nil)
  return 0
}
