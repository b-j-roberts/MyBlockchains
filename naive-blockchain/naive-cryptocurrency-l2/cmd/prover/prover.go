package main

import (
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	l2core "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/core"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

func main() { os.Exit(mainImpl()) }

//TODO: One over for error handling & logging
func mainImpl() int {
  l1Url := flag.String("l1-url", "http://localhost:8545", "L1 Url") // TODO: Better descriptions
  l1ChainId := flag.Int("l1-chainid", 505, "L1 Chain ID")
  l1TxStorageAddress := flag.String("l1-txstorage-address", "", "Main L1 contract address")
  proverL1Address := flag.String("prover-address", "", "prover address")
  proverL1Keystore := flag.String("prover-keystore", "", "Path to prover keystore")
  flag.Parse()

  l1Comms, err := l2utils.NewL1Comms(*l1Url , common.HexToAddress(*l1TxStorageAddress), common.HexToAddress("0x0"), big.NewInt(int64(*l1ChainId)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
    return 1
  }

  l2core.SetupMetrics()

  prover := l2core.NewProver(l1Comms, common.HexToAddress(*proverL1Address))
  l2utils.RegisterAccount(common.HexToAddress(*proverL1Address), *proverL1Keystore)

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
