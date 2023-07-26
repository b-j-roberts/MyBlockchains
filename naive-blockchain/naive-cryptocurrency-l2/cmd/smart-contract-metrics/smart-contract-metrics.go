package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	basicerc20 "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicerc20"
	basicl2erc20 "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicl2erc20"
	stableerc20 "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/stableerc20"
	stablel2erc20 "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/stablel2erc20"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: From root metrics such as : total rewards by account?, ...
//TODO: Token bridge metrics
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

  L1BridgeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "bridge_balance",
    Help: "Bridge balance",
  })
  L1DepositNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "deposit_nonce",
    Help: "Deposit nonce",
  })
  L1WithdrawalNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "withdrawal_nonce",
    Help: "Withdrawal nonce",
  })
  L2DepositNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_deposit_nonce",
    Help: "L2 deposit nonce",
  })
  L2WithdrawalNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_withdrawal_nonce",
    Help: "L2 withdrawal nonce",
  })
  L2BurnBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_burn_balance",
    Help: "L2 burn balance",
  })

  L1BasicTokenSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_token_supply",
    Help: "L1 basic token supply",
  })
  L1BasicTokenSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_token_sequencer_balance",
    Help: "L1 basic token sequencer balance",
  })
  L1BasicTokenBridgeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_token_bridge_balance",
    Help: "L1 basic token bridge balance",
  })
  L1StableTokenSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_stable_token_supply",
    Help: "L1 stable token supply",
  })
  L1StableTokenSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_stable_token_sequencer_balance",
    Help: "L1 stable token sequencer balance",
  })
  L1StableTokenBridgeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_stable_token_bridge_balance",
    Help: "L1 stable token bridge balance",
  })
  L1TokenDepositNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_token_deposit_nonce",
    Help: "L1  token deposit nonce",
  })
  L1TokenWithdrawalNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_token_withdrawal_nonce",
    Help: "L1  token withdrawal nonce",
  })
  L2BasicTokenSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_token_supply",
    Help: "L2 basic token supply",
  })
  L2StableTokenSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_stable_token_supply",
    Help: "L2 stable token supply",
  })
  L2TokenDepositNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_token_deposit_nonce",
    Help: "L2  token deposit nonce",
  })
  L2TokenWithdrawalNonce = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_token_withdrawal_nonce",
    Help: "L2  token withdrawal nonce",
  })
  L2BasicTokenSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_token_sequencer_balance",
    Help: "L2 basic token sequencer balance",
  })
  L2StableTokenSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_stable_token_sequencer_balance",
    Help: "L2 stable token sequencer balance",
  })
)

type SmartContractMetricExporter struct {
  L1Comms *l2utils.L1Comms
  L2Comms *l2utils.L2Comms

  ERC20Address common.Address
  ERC20Contract *basicerc20.Basicerc20

  L2ERC20Address common.Address
  L2ERC20Contract *basicl2erc20.Basicl2erc20

  StableERC20Address common.Address
  StableERC20Contract *stableerc20.Stableerc20

  L2StableERC20Address common.Address
  L2StableERC20Contract *stablel2erc20.Stablel2erc20
}

func NewSmartContractMetricExporter(l1Comms *l2utils.L1Comms, l2Comms *l2utils.L2Comms, erc20Address common.Address, l2erc20Address common.Address, stableerc20Address common.Address, l2stableerc20Address common.Address) *SmartContractMetricExporter {
  smartContractMetricExporter := &SmartContractMetricExporter{
    L1Comms: l1Comms,
    L2Comms: l2Comms,
    ERC20Address: erc20Address,
    L2ERC20Address: l2erc20Address,
    StableERC20Address: stableerc20Address,
    L2StableERC20Address: l2stableerc20Address,
  }

  var err error
  smartContractMetricExporter.ERC20Contract, err = basicerc20.NewBasicerc20(smartContractMetricExporter.ERC20Address, smartContractMetricExporter.L1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  //TODO: Get l2erc20address from l2comms contract
  smartContractMetricExporter.L2ERC20Contract, err = basicl2erc20.NewBasicl2erc20(smartContractMetricExporter.L2ERC20Address, smartContractMetricExporter.L2Comms.L2Backend)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.StableERC20Contract, err = stableerc20.NewStableerc20(smartContractMetricExporter.StableERC20Address, smartContractMetricExporter.L1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.L2StableERC20Contract, err = stablel2erc20.NewStablel2erc20(smartContractMetricExporter.L2StableERC20Address, smartContractMetricExporter.L2Comms.L2Backend)
  if err != nil {
    log.Fatal(err)
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

  prometheus.MustRegister(L1BridgeBalance)
  prometheus.MustRegister(L1DepositNonce)
  prometheus.MustRegister(L1WithdrawalNonce)
  prometheus.MustRegister(L2DepositNonce)
  prometheus.MustRegister(L2WithdrawalNonce)
  prometheus.MustRegister(L2BurnBalance)

  prometheus.MustRegister(L1BasicTokenSupply)
  prometheus.MustRegister(L1BasicTokenSequencerBalance)
  prometheus.MustRegister(L1BasicTokenBridgeBalance)
  prometheus.MustRegister(L1StableTokenSupply)
  prometheus.MustRegister(L1StableTokenSequencerBalance)
  prometheus.MustRegister(L1StableTokenBridgeBalance)
  prometheus.MustRegister(L1TokenDepositNonce)
  prometheus.MustRegister(L1TokenWithdrawalNonce)
  prometheus.MustRegister(L2BasicTokenSupply)
  prometheus.MustRegister(L2StableTokenSupply)
  prometheus.MustRegister(L2TokenDepositNonce)
  prometheus.MustRegister(L2TokenWithdrawalNonce)
  prometheus.MustRegister(L2BasicTokenSequencerBalance)
  prometheus.MustRegister(L2StableTokenSequencerBalance)
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
      L1BridgeBalance.Set(float64(bridgeBalance.Int64()))

      depositNonce, err := p.L1Comms.BridgeContract.GetEthDepositNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get deposit nonce: %v", err)
      }
      L1DepositNonce.Set(float64(depositNonce.Int64()))

      withdrawalNonce, err := p.L1Comms.BridgeContract.GetEthWithdrawNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get withdrawal nonce: %v", err)
      }
      L1WithdrawalNonce.Set(float64(withdrawalNonce.Int64()))

      depositNonce, err = p.L2Comms.L2BridgeContract.GetEthDepositNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get deposit nonce: %v", err)
      }
      L2DepositNonce.Set(float64(depositNonce.Int64()))

      withdrawalNonce, err = p.L2Comms.L2BridgeContract.GetEthWithdrawNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get withdrawal nonce: %v", err)
      }
      L2WithdrawalNonce.Set(float64(withdrawalNonce.Int64()))

      // Burn balance is l2backend balance on l2 for 0x0 address
      burnBalance, err := p.L2Comms.L2BridgeContract.GetBurntBalance(nil)
      if err != nil {
        log.Fatalf("Failed to get burn balance: %v", err)
      }
      L2BurnBalance.Set(float64(burnBalance.Int64()))



      basicTokenSupply, err := p.ERC20Contract.TotalSupply(nil)
      if err != nil {
        log.Fatalf("Failed to get basic token supply: %v", err)
      }
      L1BasicTokenSupply.Set(float64(basicTokenSupply.Int64()))

      log.Println("Getting basic token sequencer balance from ", l2utils.GetSequencer())
      basicTokenSequencerBalance, err := p.ERC20Contract.BalanceOf(nil, l2utils.GetSequencer())
      if err != nil {
        log.Fatalf("Failed to get basic token sequencer balance: %v", err)
      }
      L1BasicTokenSequencerBalance.Set(float64(basicTokenSequencerBalance.Int64()))

      log.Println("Getting token balance for token bridge account ", p.L1Comms.TokenBridgeContractAddress)
      basicTokenBalance, err := p.ERC20Contract.BalanceOf(nil, p.L1Comms.TokenBridgeContractAddress)
      if err != nil {
        log.Fatalf("Failed to get basic token balance: %v", err)
      }
      L1BasicTokenBridgeBalance.Set(float64(basicTokenBalance.Int64()))

      stableTokenSupply, err := p.StableERC20Contract.TotalSupply(nil)
      if err != nil {
        log.Fatalf("Failed to get stable token supply: %v", err)
      }
      L1StableTokenSupply.Set(float64(stableTokenSupply.Int64()))

      log.Println("Getting stable token sequencer balance from ", l2utils.GetSequencer())
      stableTokenSequencerBalance, err := p.StableERC20Contract.BalanceOf(nil, l2utils.GetSequencer())
      if err != nil {
        log.Fatalf("Failed to get stable token sequencer balance: %v", err)
      }
      L1StableTokenSequencerBalance.Set(float64(stableTokenSequencerBalance.Int64()))

      log.Println("Getting token balance for token bridge account ", p.L1Comms.TokenBridgeContractAddress)
      stableTokenBalance, err := p.StableERC20Contract.BalanceOf(nil, p.L1Comms.TokenBridgeContractAddress)
      if err != nil {
        log.Fatalf("Failed to get stable token balance: %v", err)
      }
      L1StableTokenBridgeBalance.Set(float64(stableTokenBalance.Int64()))

      basicTokenDepositNonce, err := p.L1Comms.TokenBridgeContract.GetTokenDepositNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get basic token deposit nonce: %v", err)
      }
      L1TokenDepositNonce.Set(float64(basicTokenDepositNonce.Int64()))

      basicTokenWithdrawalNonce, err := p.L1Comms.TokenBridgeContract.GetTokenWithdrawNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get basic token withdrawal nonce: %v", err)
      }
      L1TokenWithdrawalNonce.Set(float64(basicTokenWithdrawalNonce.Int64()))

      l2BasicTokenSupply, err := p.L2ERC20Contract.TotalSupply(nil)
      if err != nil {
        log.Fatalf("Failed to get basic token supply: %v", err)
      }
      L2BasicTokenSupply.Set(float64(l2BasicTokenSupply.Int64()))
      
      l2StableTokenSupply, err := p.L2StableERC20Contract.TotalSupply(nil)
      if err != nil {
        log.Fatalf("Failed to get stable token supply: %v", err)
      }
      L2StableTokenSupply.Set(float64(l2StableTokenSupply.Int64()))

      l2BasicTokenDepositNonce, err := p.L2Comms.L2TokenBridgeContract.GetTokenDepositNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get basic token deposit nonce: %v", err)
      }
      L2TokenDepositNonce.Set(float64(l2BasicTokenDepositNonce.Int64()))

      l2BasicTokenWithdrawalNonce, err := p.L2Comms.L2TokenBridgeContract.GetTokenWithdrawNonce(nil)
      if err != nil {
        log.Fatalf("Failed to get basic token withdrawal nonce: %v", err)
      }
      L2TokenWithdrawalNonce.Set(float64(l2BasicTokenWithdrawalNonce.Int64()))

      l2BasicTokenSequencerBalance, err := p.L2ERC20Contract.BalanceOf(nil, common.HexToAddress(l2utils.GetSequencer().Hex()))
      if err != nil {
        log.Fatalf("Failed to get basic token sequencer balance: %v", err)
      }
      L2BasicTokenSequencerBalance.Set(float64(l2BasicTokenSequencerBalance.Int64()))

      l2StableTokenSequencerBalance, err := p.L2StableERC20Contract.BalanceOf(nil, common.HexToAddress(l2utils.GetSequencer().Hex()))
      if err != nil {
        log.Fatalf("Failed to get stable token sequencer balance: %v", err)
      }
      L2StableTokenSequencerBalance.Set(float64(l2StableTokenSequencerBalance.Int64()))

      // Sleep for 3 seconds
      time.Sleep(3 * time.Second)
    }
  }()

  return nil
}

func main() { os.Exit(mainImp()) }

func mainImp() int {
  sequencer := flag.String("sequencer", "", "Sequencer address")
  l1ContractAddress := flag.String("l1-tx-storage-address", "", "Main L1 contract address")
  l1BridgeAddress := flag.String("l1-bridge-address", "", "Main L1 contract address")
  //TODO: l2bridgeaddress?
  l1TokenBridgeAddress := flag.String("l1-token-bridge-address", "", "Main L1 contract address")
  l2TokenBridgeAddress := flag.String("l2-token-bridge-address", "", "Main L1 contract address")
  l1Url := flag.String("l1-url", "http://localhost:8545", "L1 URL")
  l1ChainId := flag.Int("l1-chainid", 505, "L1 chain ID")
  l2ChainId := flag.Int("l2-chainid", 515, "L1 chain ID")
  erc20Address := flag.String("erc20-address", "", "ERC20 address")
  l2erc20Address := flag.String("l2-erc20-address", "", "ERC20 address on L2")
  stableERC20Address := flag.String("stable-erc20-address", "", "Stable ERC20 address")
  l2StableERC20Address := flag.String("l2-stable-erc20-address", "", "Stable ERC20 address on L2")
  l2BridgeAddress := flag.String("l2-bridge-address", "", "Main L2 contract address")
  l2IPCPath := flag.String("l2-ipc-path", "/home/brandon/naive-sequencer-data/naive-sequencer.ipc", "L2 IPC path")
  flag.Parse()

  if *sequencer == "" {
    log.Fatalf("Must provide sequencer address")
  }
  l2utils.SetSequencer(common.HexToAddress(*sequencer))
  l1Comms, err := l2utils.NewL1Comms(*l1Url , common.HexToAddress(*l1ContractAddress), common.HexToAddress(*l1BridgeAddress), common.HexToAddress(*l1TokenBridgeAddress), big.NewInt(int64(*l1ChainId)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
  }

  rpcIPC, err := rpc.DialIPC(context.Background(), *l2IPCPath)
  if err != nil {
    log.Fatalf("Failed to dial ipc: %v", err)
  }

  backend := ethclient.NewClient(rpcIPC)

  l2Comms, err := l2utils.NewL2Comms(common.HexToAddress(*l2BridgeAddress), common.HexToAddress(*l2TokenBridgeAddress), big.NewInt(int64(*l2ChainId)), backend, l2utils.GetDefaultL2TransactionConfig())
  if err != nil {
    log.Fatalf("Failed to create L2 comms: %v", err)
  }

  SetupMetrics()

  smartContractMetricExporter := NewSmartContractMetricExporter(l1Comms, l2Comms, common.HexToAddress(*erc20Address), common.HexToAddress(*l2erc20Address), common.HexToAddress(*stableERC20Address), common.HexToAddress(*l2StableERC20Address))

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
