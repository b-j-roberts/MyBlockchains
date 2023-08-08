package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicerc20"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicerc721"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicl2erc20"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicl2erc721"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/specialerc721"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/speciall2erc721"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/stableerc20"
	"github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/stablel2erc20"
	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type L1TokenContracts struct {
  ERC20Contract *basicerc20.Basicerc20
  ERC721Contract *basicerc721.Basicerc721
  StableERC20Contract *stableerc20.Stableerc20
  SpecialERC721Contract *specialerc721.Specialerc721
}

func LoadL1TokenContracts(l1Comms *l2utils.L1Comms, tokenAddresses l2utils.TokenAddresses) L1TokenContracts {
  var err error
  l1TokenContracts := L1TokenContracts{}

  l1TokenContracts.ERC20Contract, err = basicerc20.NewBasicerc20(tokenAddresses.Erc20Address, l1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  l1TokenContracts.ERC721Contract, err = basicerc721.NewBasicerc721(tokenAddresses.Erc721Address, l1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  l1TokenContracts.StableERC20Contract, err = stableerc20.NewStableerc20(tokenAddresses.StableErc20Address, l1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  l1TokenContracts.SpecialERC721Contract, err = specialerc721.NewSpecialerc721(tokenAddresses.SpecialErc721Address, l1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  return l1TokenContracts
}

type L2TokenContracts struct {
  ERC20Contract *basicl2erc20.Basicl2erc20
  ERC721Contract *basicl2erc721.Basicl2erc721
  StableERC20Contract *stablel2erc20.Stablel2erc20
  SpecialERC721Contract *speciall2erc721.Speciall2erc721
}

func LoadL2TokenContracts(l2Comms *l2utils.L2Comms, tokenAddresses l2utils.TokenAddresses) L2TokenContracts {
  var err error
  l2TokenContracts := L2TokenContracts{}

  l2TokenContracts.ERC20Contract, err = basicl2erc20.NewBasicl2erc20(tokenAddresses.L2Erc20Address, l2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  l2TokenContracts.ERC721Contract, err = basicl2erc721.NewBasicl2erc721(tokenAddresses.L2Erc721Address, l2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  l2TokenContracts.StableERC20Contract, err = stablel2erc20.NewStablel2erc20(tokenAddresses.L2StableErc20Address, l2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  l2TokenContracts.SpecialERC721Contract, err = speciall2erc721.NewSpeciall2erc721(tokenAddresses.L2SpecialErc721Address, l2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  return l2TokenContracts
}

type SmartContractMetricExporter struct {
  L1Comms *l2utils.L1Comms
  L2Comms *l2utils.L2Comms

  TokenAddresses l2utils.TokenAddresses

  L1TokenContracts L1TokenContracts
  L2TokenContracts L2TokenContracts
}

func NewSmartContractMetricExporter(l1Comms *l2utils.L1Comms, l2Comms *l2utils.L2Comms, tokenAddresses l2utils.TokenAddresses) *SmartContractMetricExporter {
  return &SmartContractMetricExporter{
    L1Comms: l1Comms,
    L2Comms: l2Comms,
    TokenAddresses: tokenAddresses,
    L1TokenContracts: LoadL1TokenContracts(l1Comms, tokenAddresses),
    L2TokenContracts: LoadL2TokenContracts(l2Comms, tokenAddresses),
  }
}

func (p *SmartContractMetricExporter) CollectBatchMetrics() {
  log.Println("Updating batch metrics from ", p.L1Comms.L1ContractAddressConfig.TxStorageContractAddress.String())
  // Update metric values
  batchCount, err := p.L1Comms.L1Contracts.TxStorageContract.GetBatchCount(nil)
  if err != nil {
    log.Fatalf("Failed to get batch count: %v", err)
  }
  BatchCount.Set(float64(batchCount.Int64()))

  lastConfirmedBatch, err := p.L1Comms.L1Contracts.TxStorageContract.GetLastConfirmedBatch(nil)
  if err != nil {
    log.Fatalf("Failed to get last confirmed batch: %v", err)
  }
  LastConfirmedBatch.Set(float64(lastConfirmedBatch.Int64()))

  latestBatchL1Block, err := p.L1Comms.L1Contracts.TxStorageContract.GetBatchL1Block(nil, big.NewInt(int64(batchCount.Int64() - 1)))
  if err != nil {
    log.Fatalf("Failed to get latest batch L1 block: %v", err)
  }
  LatestBatchL1Block.Set(float64(latestBatchL1Block.Int64()))

  latestBatchProofL1Block, err := p.L1Comms.L1Contracts.TxStorageContract.GetProofL1Block(nil, big.NewInt(int64(batchCount.Int64() - 1)))
  if err != nil {
    log.Fatalf("Failed to get latest batch proof L1 block: %v", err)
  }
  LatestBatchProofL1Block.Set(float64(latestBatchProofL1Block.Int64()))

  latestConfirmedBatchL1Block, err := p.L1Comms.L1Contracts.TxStorageContract.GetBatchL1Block(nil, big.NewInt(int64(lastConfirmedBatch.Int64())))
  if err != nil {
    log.Fatalf("Failed to get latest confirmed batch L1 block: %v", err)
  }
  LatestConfirmedBatchL1Block.Set(float64(latestConfirmedBatchL1Block.Int64()))

  latestConfirmedBatchProofL1Block, err := p.L1Comms.L1Contracts.TxStorageContract.GetProofL1Block(nil, big.NewInt(int64(lastConfirmedBatch.Int64())))
  if err != nil {
    log.Fatalf("Failed to get latest confirmed batch proof L1 block: %v", err)
  }
  LatestConfirmedBatchProofL1Block.Set(float64(latestConfirmedBatchProofL1Block.Int64()))
}

func (p *SmartContractMetricExporter) CollectEthBridgeMetrics() {
  bridgeBalance, err := p.L1Comms.L1Contracts.BridgeContract.GetBridgeBalance(nil)
  if err != nil {
    log.Fatalf("Failed to get bridge balance: %v", err)
  }
  L1BridgeBalance.Set(float64(bridgeBalance.Int64()))

  depositNonce, err := p.L1Comms.L1Contracts.BridgeContract.GetEthDepositNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get deposit nonce: %v", err)
  }
  L1DepositNonce.Set(float64(depositNonce.Int64()))

  withdrawalNonce, err := p.L1Comms.L1Contracts.BridgeContract.GetEthWithdrawNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get withdrawal nonce: %v", err)
  }
  L1WithdrawalNonce.Set(float64(withdrawalNonce.Int64()))

  depositNonce, err = p.L2Comms.L2Contracts.L2BridgeContract.GetEthDepositNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get deposit nonce: %v", err)
  }
  L2DepositNonce.Set(float64(depositNonce.Int64()))

  withdrawalNonce, err = p.L2Comms.L2Contracts.L2BridgeContract.GetEthWithdrawNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get withdrawal nonce: %v", err)
  }
  L2WithdrawalNonce.Set(float64(withdrawalNonce.Int64()))

  // Burn balance is l2backend balance on l2 for 0x0 address
  burnBalance, err := p.L2Comms.L2Contracts.L2BridgeContract.GetBurntBalance(nil)
  if err != nil {
    log.Fatalf("Failed to get burn balance: %v", err)
  }
  L2BurnBalance.Set(float64(burnBalance.Int64()))
}

func (p *SmartContractMetricExporter) CollectTokenBridgeMetrics() {
  basicTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.Erc20Address)
  if err != nil {
    log.Fatalf("Failed to get basic token allowed: %v", err)
  }
  L1BasicTokenAllowed.Set(float64(basicTokenAllowed))

  basicTokenSupply, err := p.L1TokenContracts.ERC20Contract.TotalSupply(nil)
  if err != nil {
    log.Fatalf("Failed to get basic token supply: %v", err)
  }
  L1BasicTokenSupply.Set(float64(basicTokenSupply.Int64()))

  log.Println("Getting basic token sequencer balance from ", l2utils.GetSequencer())
  basicTokenSequencerBalance, err := p.L1TokenContracts.ERC20Contract.BalanceOf(nil, l2utils.GetSequencer())
  if err != nil {
    log.Fatalf("Failed to get basic token sequencer balance: %v", err)
  }
  L1BasicTokenSequencerBalance.Set(float64(basicTokenSequencerBalance.Int64()))

  log.Println("Getting token balance for token bridge account ", p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
  basicTokenBalance, err := p.L1TokenContracts.ERC20Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
  if err != nil {
    log.Fatalf("Failed to get basic token balance: %v", err)
  }
  L1BasicTokenBridgeBalance.Set(float64(basicTokenBalance.Int64()))

  stableTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.StableErc20Address)
  if err != nil {
    log.Fatalf("Failed to get stable token allowed: %v", err)
  }
  L1StableTokenAllowed.Set(float64(stableTokenAllowed))

  stableTokenSupply, err := p.L1TokenContracts.StableERC20Contract.TotalSupply(nil)
  if err != nil {
    log.Fatalf("Failed to get stable token supply: %v", err)
  }
  L1StableTokenSupply.Set(float64(stableTokenSupply.Int64()))

  log.Println("Getting stable token sequencer balance from ", l2utils.GetSequencer())
  stableTokenSequencerBalance, err := p.L1TokenContracts.StableERC20Contract.BalanceOf(nil, l2utils.GetSequencer())
  if err != nil {
    log.Fatalf("Failed to get stable token sequencer balance: %v", err)
  }
  L1StableTokenSequencerBalance.Set(float64(stableTokenSequencerBalance.Int64()))

  log.Println("Getting token balance for token bridge account ", p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
  stableTokenBalance, err := p.L1TokenContracts.StableERC20Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
  if err != nil {
    log.Fatalf("Failed to get stable token balance: %v", err)
  }
  L1StableTokenBridgeBalance.Set(float64(stableTokenBalance.Int64()))

  basicTokenDepositNonce, err := p.L1Comms.L1Contracts.TokenBridgeContract.GetTokenDepositNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get basic token deposit nonce: %v", err)
  }
  L1TokenDepositNonce.Set(float64(basicTokenDepositNonce.Int64()))

  basicTokenWithdrawalNonce, err := p.L1Comms.L1Contracts.TokenBridgeContract.GetTokenWithdrawNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get basic token withdrawal nonce: %v", err)
  }
  L1TokenWithdrawalNonce.Set(float64(basicTokenWithdrawalNonce.Int64()))

  l2BasicTokenAllowed, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetAllowedToken(nil, p.TokenAddresses.Erc20Address)
  if err != nil {
    log.Fatalf("Failed to get basic token allowed: %v", err)
  }
  if l2BasicTokenAllowed {
    L2BasicTokenAllowed.Set(1)
  } else {
    L2BasicTokenAllowed.Set(0)
  }

  l2BasicTokenSupply, err := p.L2TokenContracts.ERC20Contract.TotalSupply(nil)
  if err != nil {
    log.Fatalf("Failed to get basic token supply: %v", err)
  }
  L2BasicTokenSupply.Set(float64(l2BasicTokenSupply.Int64()))

  l2StableTokenAllowed, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetAllowedToken(nil, p.TokenAddresses.StableErc20Address)
  if err != nil {
    log.Fatalf("Failed to get stable token allowed: %v", err)
  }
  if l2StableTokenAllowed {
    L2StableTokenAllowed.Set(1)
  } else {
    L2StableTokenAllowed.Set(0)
  }
  
  l2StableTokenSupply, err := p.L2TokenContracts.ERC20Contract.TotalSupply(nil)
  if err != nil {
    log.Fatalf("Failed to get stable token supply: %v", err)
  }
  L2StableTokenSupply.Set(float64(l2StableTokenSupply.Int64()))

  l2BasicTokenDepositNonce, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetTokenDepositNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get basic token deposit nonce: %v", err)
  }
  L2TokenDepositNonce.Set(float64(l2BasicTokenDepositNonce.Int64()))

  l2BasicTokenWithdrawalNonce, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetTokenWithdrawNonce(nil)
  if err != nil {
    log.Fatalf("Failed to get basic token withdrawal nonce: %v", err)
  }
  L2TokenWithdrawalNonce.Set(float64(l2BasicTokenWithdrawalNonce.Int64()))

  l2BasicTokenSequencerBalance, err := p.L2TokenContracts.ERC20Contract.BalanceOf(nil, common.HexToAddress(l2utils.GetSequencer().Hex()))
  if err != nil {
    log.Fatalf("Failed to get basic token sequencer balance: %v", err)
  }
  L2BasicTokenSequencerBalance.Set(float64(l2BasicTokenSequencerBalance.Int64()))

  l2StableTokenSequencerBalance, err := p.L2TokenContracts.ERC20Contract.BalanceOf(nil, common.HexToAddress(l2utils.GetSequencer().Hex()))
  if err != nil {
    log.Fatalf("Failed to get stable token sequencer balance: %v", err)
  }
  L2StableTokenSequencerBalance.Set(float64(l2StableTokenSequencerBalance.Int64()))

  basicNFTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.Erc721Address)
  if err != nil {
    log.Fatalf("Failed to get basic NFT allowed: %v", err)
  }
  L1BasicNFTAllowed.Set(float64(basicNFTokenAllowed))

  basicNFTokenSupply, err := p.L1TokenContracts.ERC721Contract.TotalSupply(nil)
  if err != nil {
    log.Fatalf("Failed to get basic NFT token supply: %v", err)
  }
  L1BasicNFTSupply.Set(float64(basicNFTokenSupply.Int64()))

  basicNFTokenSequencerBalance, err := p.L1TokenContracts.ERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
  if err != nil {
    log.Fatalf("Failed to get basic NFT token sequencer balance: %v", err)
  }
  L1BasicNFTSequencerBalance.Set(float64(basicNFTokenSequencerBalance.Int64()))

  basicNFTBridgeBalance, err := p.L1TokenContracts.ERC721Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
  if err != nil {
    log.Fatalf("Failed to get basic NFT token bridge balance: %v", err)
  }
  L1BasicNFTBridgeBalance.Set(float64(basicNFTBridgeBalance.Int64()))

  specialNFTTokenName, err := p.L1TokenContracts.SpecialERC721Contract.Name(nil)
  if err != nil {
    log.Fatalf("Failed to get special NFT token name: %v", err)
  }
  if specialNFTTokenName != "" {
    L1SpecialNFTSupply.Set(float64(1))
  } else {
    L1SpecialNFTSupply.Set(float64(0))
  }

  specialNFTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.SpecialErc721Address)
  if err != nil {
    log.Fatalf("Failed to get special NFT allowed: %v", err)
  }
  L1SpecialNFTAllowed.Set(float64(specialNFTokenAllowed))

  specialNFTSequncerBalance, err := p.L1TokenContracts.SpecialERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
  if err != nil {
    log.Fatalf("Failed to get special NFT token sequencer balance: %v", err)
  }
  L1SpecialNFTSequencerBalance.Set(float64(specialNFTSequncerBalance.Int64()))

  specialNFTBridgeBalance, err := p.L1TokenContracts.SpecialERC721Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
  if err != nil {
    log.Fatalf("Failed to get special NFT token bridge balance: %v", err)
  }
  L1SpecialNFTBridgeBalance.Set(float64(specialNFTBridgeBalance.Int64()))

  l2SpecialNFTSequncerBalance, err := p.L2TokenContracts.SpecialERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
  if err != nil {
    log.Fatalf("Failed to get special NFT token sequencer balance: %v", err)
  }
  L2SpecialNFTSequencerBalance.Set(float64(l2SpecialNFTSequncerBalance.Int64()))

  l2BasicNFTokenAllowed, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetAllowedToken(nil, p.TokenAddresses.Erc721Address)
  if err != nil {
    log.Fatalf("Failed to get basic NFT allowed: %v", err)
  }
  if l2BasicNFTokenAllowed {
    L2BasicNFTAllowed.Set(float64(1))
  } else {
    L2BasicNFTAllowed.Set(float64(0))
  }

  l2BasicNFTSupply, err := p.L2TokenContracts.ERC721Contract.TotalSupply(nil)
  if err != nil {
    log.Fatalf("Failed to get basic NFT token supply: %v", err)
  }
  L2BasicNFTSupply.Set(float64(l2BasicNFTSupply.Int64()))

  l2BasicNFTSequencerBalance, err := p.L2TokenContracts.ERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
  if err != nil {
    log.Fatalf("Failed to get basic NFT token sequencer balance: %v", err)
  }
  L2BasicNFTSequencerBalance.Set(float64(l2BasicNFTSequencerBalance.Int64()))

  l2SpecialNFTokenAllowed, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetAllowedToken(nil, p.TokenAddresses.SpecialErc721Address)
  if err != nil {
    log.Fatalf("Failed to get special NFT allowed: %v", err)
  }
  if l2SpecialNFTokenAllowed {
    L2SpecialNFTAllowed.Set(float64(1))
  } else {
    L2SpecialNFTAllowed.Set(float64(0))
  }
}

func (p *SmartContractMetricExporter) Start() error {
  log.Println("Starting Smart Contract Metric Exporter...")

  go func() {
    for {
      l1BlockHeight, err := p.L1Comms.L1Client.BlockNumber(context.Background())
      if err != nil {
        log.Fatalf("Failed to get L1 block height: %v", err)
      }
      L1BlockHeight.Set(float64(l1BlockHeight))

      p.CollectBatchMetrics()
      p.CollectEthBridgeMetrics()
      p.CollectTokenBridgeMetrics()

      // Sleep for 3 seconds
      time.Sleep(3 * time.Second)
    }
  }()

  return nil
}

func main() { os.Exit(mainImp()) }

func mainImp() int {
  osHomeDir, err := os.UserHomeDir()
  sequencer := flag.String("sequencer", "", "Sequencer address")
  configFile := flag.String("config", osHomeDir + "/naive-sequencer-data/sequencer.config.json", "Config file")
  flag.Parse()

  config, err := l2config.LoadNodeBaseConfig(*configFile)
  if err != nil {
    log.Fatalf("Failed to load config: %v", err)
    return 1
  }

  tokenAddresses, err := l2utils.LoadTokenAddresses(config.Contracts)
  if err != nil {
    log.Fatalf("Failed to load token addresses: %v", err)
    return 1
  }

  if *sequencer == "" {
    log.Fatalf("Must provide sequencer address")
    return 1
  }
  l2utils.SetSequencer(common.HexToAddress(*sequencer))

  l1Comms, err := l2utils.NewL1Comms(config.L1URL, config.Contracts, big.NewInt(int64(config.L1ChainID)), l2utils.L1TransactionConfig{
    GasLimit: 3000000,
    GasPrice: big.NewInt(200),
  })
  if err != nil {
    log.Fatalf("Failed to create L1 comms: %v", err)
    return 1
  }

  l2Comms, err := l2utils.NewL2Comms(config, l2utils.GetDefaultL2TransactionConfig())
  if err != nil {
    log.Fatalf("Failed to create L2 comms: %v", err)
  }

  SetupMetrics()

  smartContractMetricExporter := NewSmartContractMetricExporter(l1Comms, l2Comms, tokenAddresses)

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
