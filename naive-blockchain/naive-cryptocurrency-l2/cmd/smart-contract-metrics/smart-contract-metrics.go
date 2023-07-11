package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: batchCount, lastConfirmedBatch, l1 block height, latest batches l1 block, latest batches latest proof block, latest batches( all for latest & latest proved )
// TODO: From root metrics such as : total rewards by account?, ...

var (
  BatchCount = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "batch_count",
    Help: "Number of batches",
  })
  LastConfirmedBatch = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "last_confirmed_batch",
    Help: "Last confirmed batch number",
  })
  L1BlockHeight = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_block_height",
    Help: "L1 block height at time of metric collection",
  })
  LatestBatchL1Block = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "latest_batch_l1_block",
    Help: "L1 block height of transaction storing latest batch",
  })
  LatestBatchProofL1Block = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "latest_batch_proof_l1_block",
    Help: "L1 block height of transaction storing latest batch proof",
  })
  LatestConfirmedBatchL1Block = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "latest_confirmed_batch_l1_block",
    Help: "L1 block height of transaction storing latest confirmed batch",
  })
  LatestConfirmedBatchProofL1Block = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "latest_confirmed_batch_proof_l1_block",
    Help: "L1 block height of transaction storing latest confirmed batch proof",
  })
  BridgeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "bridge_balance",
    Help: "Bridge balance",
  })
)

type SmartContractMetricExporter struct {
  L1Comms *l2utils.L1Comms
}

func NewSmartContractMetricExporter(l1Comms *l2utils.L1Comms) *SmartContractMetricExporter {
  smartContractMetricExporter := &SmartContractMetricExporter{
    L1Comms: l1Comms,
  }

  return smartContractMetricExporter
}

func SetupMetrics() {
  prometheus.MustRegister(BatchCount)
  prometheus.MustRegister(LastConfirmedBatch)
  prometheus.MustRegister(L1BlockHeight)
  prometheus.MustRegister(LatestBatchL1Block)
  prometheus.MustRegister(LatestBatchProofL1Block)
  prometheus.MustRegister(LatestConfirmedBatchL1Block)
  prometheus.MustRegister(LatestConfirmedBatchProofL1Block)
  prometheus.MustRegister(BridgeBalance)
}

func (p *SmartContractMetricExporter) Start() error {
  log.Println("Starting Smart Contract Metric Exporter...")

  go func() {
    for {
      log.Println("Updating smart contract metrics from ", p.L1Comms.TxStorageContractAddress.String())
      // Update metric values
      batchCount, err := p.L1Comms.TxStorageContract.GetBatchCount(nil)
      if err != nil {
        log.Fatalf("Failed to get batch count: %v", err)
      }
      BatchCount.Set(float64(batchCount.Int64()))

      lastConfirmedBatch, err := p.L1Comms.TxStorageContract.GetLastConfirmedBatch(nil)
      if err != nil {
        log.Fatalf("Failed to get last confirmed batch: %v", err)
      }
      LastConfirmedBatch.Set(float64(lastConfirmedBatch.Int64()))

      l1BlockHeight, err := p.L1Comms.L1Client.BlockNumber(context.Background())
      if err != nil {
        log.Fatalf("Failed to get L1 block height: %v", err)
      }
      L1BlockHeight.Set(float64(l1BlockHeight))

      latestBatchL1Block, err := p.L1Comms.TxStorageContract.GetBatchL1Block(nil, big.NewInt(int64(batchCount.Int64() - 1)))
      if err != nil {
        log.Fatalf("Failed to get latest batch L1 block: %v", err)
      }
      LatestBatchL1Block.Set(float64(latestBatchL1Block.Int64()))

      latestBatchProofL1Block, err := p.L1Comms.TxStorageContract.GetProofL1Block(nil, big.NewInt(int64(batchCount.Int64() - 1)))
      if err != nil {
        log.Fatalf("Failed to get latest batch proof L1 block: %v", err)
      }
      LatestBatchProofL1Block.Set(float64(latestBatchProofL1Block.Int64()))

      latestConfirmedBatchL1Block, err := p.L1Comms.TxStorageContract.GetBatchL1Block(nil, big.NewInt(int64(lastConfirmedBatch.Int64())))
      if err != nil {
        log.Fatalf("Failed to get latest confirmed batch L1 block: %v", err)
      }
      LatestConfirmedBatchL1Block.Set(float64(latestConfirmedBatchL1Block.Int64()))

      latestConfirmedBatchProofL1Block, err := p.L1Comms.TxStorageContract.GetProofL1Block(nil, big.NewInt(int64(lastConfirmedBatch.Int64())))
      if err != nil {
        log.Fatalf("Failed to get latest confirmed batch proof L1 block: %v", err)
      }
      LatestConfirmedBatchProofL1Block.Set(float64(latestConfirmedBatchProofL1Block.Int64()))

      bridgeBalance, err := p.L1Comms.BridgeContract.GetBridgeBalance(nil)
      if err != nil {
        log.Fatalf("Failed to get bridge balance: %v", err)
      }
      BridgeBalance.Set(float64(bridgeBalance.Int64()))

      // Sleep for 3 seconds
      time.Sleep(3 * time.Second)
    }
  }()

  return nil
}

func main() { os.Exit(mainImp()) }

func mainImp() int {
  l1ContractAddress := flag.String("l1-contract-address", "", "Main L1 contract address")
  l1BridgeAddress := flag.String("l1-bridge-address", "", "Main L1 contract address")
  l1Host := flag.String("l1-host", "http://localhost", "L1 host")
  l1Port := flag.String("l1-port", "8545", "L1 port")
  l1ChainId := flag.Int("l1-chainid", 505, "L1 chain ID")
  flag.Parse()

  l1Url := *l1Host + ":" + *l1Port
  l1Comms, err := l2utils.NewL1Comms(l1Url , common.HexToAddress(*l1ContractAddress), common.HexToAddress(*l1BridgeAddress), big.NewInt(int64(*l1ChainId)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
  }

  SetupMetrics()

  smartContractMetricExporter := NewSmartContractMetricExporter(l1Comms)

  fatalErrChan := make(chan error, 10)
  err = smartContractMetricExporter.Start()
  if err != nil {
    fatalErrChan <- err
  }

  log.Println("Starting Prometheus metrics server on port 6169...")
  http.Handle("/metrics", promhttp.Handler())
  http.ListenAndServe(":6169", nil)
  //sigint := make(chan os.Signal, 1)
  //signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

  exitCode := 0
  //select {
  //case err := <-fatalErrChan:
  //  log.Println("shutting down due to fatal error:", err)
  //  defer log.Println("shut down")
  //  exitCode = 1
  //case <-sigint:
  //  log.Println("shutting down due to interrupt")
  //}

  //close(sigint)

  return exitCode
}
