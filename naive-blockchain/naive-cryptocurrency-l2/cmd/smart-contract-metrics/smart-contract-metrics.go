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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

  L1BasicTokenAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_token_allowed",
    Help: "L1 basic token allowed",
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
  L1StableTokenAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_stable_token_allowed",
    Help: "L1 stable token allowed",
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
  L2BasicTokenAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_token_allowed",
    Help: "L2 basic token allowed",
  })
  L2BasicTokenSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_token_supply",
    Help: "L2 basic token supply",
  })
  L2StableTokenAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_stable_token_allowed",
    Help: "L2 basic token allowed",
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

  L1BasicNFTAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_nft_allowed",
    Help: "L1 basic NFT allowed",
  })
  L1BasicNFTSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_nft_supply",
    Help: "L1 basic NFT supply",
  })
  L1BasicNFTSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_nft_sequencer_balance",
    Help: "L1 basic NFT sequencer balance",
  })
  L1BasicNFTBridgeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_basic_nft_bridge_balance",
    Help: "L1 basic NFT bridge balance",
  })
  L1SpecialNFTAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_special_nft_allowed",
    Help: "L1 special NFT allowed",
  })
  L1SpecialNFTSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_special_nft_supply",
    Help: "L1 special NFT supply",
  })
  L1SpecialNFTSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_special_nft_sequencer_balance",
    Help: "L1 special NFT sequencer balance",
  })
  L1SpecialNFTBridgeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l1_special_nft_bridge_balance",
    Help: "L1 special NFT bridge balance",
  })
  L2BasicNFTAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_nft_allowed",
    Help: "L2 basic NFT allowed",
  })
  L2BasicNFTSupply = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_nft_supply",
    Help: "L2 basic NFT supply",
  })
  L2SpecialNFTAllowed = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_special_nft_allowed",
    Help: "L2 special NFT allowed",
  })
  L2BasicNFTSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_basic_nft_sequencer_balance",
    Help: "L2 basic NFT sequencer balance",
  })
  L2SpecialNFTSequencerBalance = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "l2_special_nft_sequencer_balance",
    Help: "L2 special NFT sequencer balance",
  })
)

type SmartContractMetricExporter struct {
  L1Comms *l2utils.L1Comms
  L2Comms *l2utils.L2Comms

  TokenAddresses l2utils.TokenAddresses

  ERC20Contract *basicerc20.Basicerc20
  L2ERC20Contract *basicl2erc20.Basicl2erc20
  StableERC20Contract *stableerc20.Stableerc20
  L2StableERC20Contract *stablel2erc20.Stablel2erc20
  ERC721Contract *basicerc721.Basicerc721
  L2ERC721Contract *basicl2erc721.Basicl2erc721
  SpecialERC721Contract *specialerc721.Specialerc721
  L2SpecialERC721Contract *speciall2erc721.Speciall2erc721
}

func NewSmartContractMetricExporter(l1Comms *l2utils.L1Comms, l2Comms *l2utils.L2Comms, tokenAddresses l2utils.TokenAddresses) *SmartContractMetricExporter {
  smartContractMetricExporter := &SmartContractMetricExporter{
    L1Comms: l1Comms,
    L2Comms: l2Comms,
    TokenAddresses: tokenAddresses,
  }

  var err error
  smartContractMetricExporter.ERC20Contract, err = basicerc20.NewBasicerc20(smartContractMetricExporter.TokenAddresses.Erc20Address, smartContractMetricExporter.L1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  //TODO: Get l2erc20address from l2comms contract using allowedTokens
  smartContractMetricExporter.L2ERC20Contract, err = basicl2erc20.NewBasicl2erc20(smartContractMetricExporter.TokenAddresses.L2Erc20Address, smartContractMetricExporter.L2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.StableERC20Contract, err = stableerc20.NewStableerc20(smartContractMetricExporter.TokenAddresses.StableErc20Address, smartContractMetricExporter.L1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.L2StableERC20Contract, err = stablel2erc20.NewStablel2erc20(smartContractMetricExporter.TokenAddresses.L2StableErc20Address, smartContractMetricExporter.L2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.ERC721Contract, err = basicerc721.NewBasicerc721(smartContractMetricExporter.TokenAddresses.Erc721Address, smartContractMetricExporter.L1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.L2ERC721Contract, err = basicl2erc721.NewBasicl2erc721(smartContractMetricExporter.TokenAddresses.L2Erc721Address, smartContractMetricExporter.L2Comms.L2Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.SpecialERC721Contract, err = specialerc721.NewSpecialerc721(smartContractMetricExporter.TokenAddresses.SpecialErc721Address, smartContractMetricExporter.L1Comms.L1Client)
  if err != nil {
    log.Fatal(err)
  }

  smartContractMetricExporter.L2SpecialERC721Contract, err = speciall2erc721.NewSpeciall2erc721(smartContractMetricExporter.TokenAddresses.L2SpecialErc721Address, smartContractMetricExporter.L2Comms.L2Client)
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

  prometheus.MustRegister(L1BasicTokenAllowed)
  prometheus.MustRegister(L1BasicTokenSupply)
  prometheus.MustRegister(L1BasicTokenSequencerBalance)
  prometheus.MustRegister(L1BasicTokenBridgeBalance)
  prometheus.MustRegister(L1StableTokenAllowed)
  prometheus.MustRegister(L1StableTokenSupply)
  prometheus.MustRegister(L1StableTokenSequencerBalance)
  prometheus.MustRegister(L1StableTokenBridgeBalance)
  prometheus.MustRegister(L1TokenDepositNonce)
  prometheus.MustRegister(L1TokenWithdrawalNonce)
  prometheus.MustRegister(L2BasicTokenAllowed)
  prometheus.MustRegister(L2BasicTokenSupply)
  prometheus.MustRegister(L2StableTokenAllowed)
  prometheus.MustRegister(L2StableTokenSupply)
  prometheus.MustRegister(L2TokenDepositNonce)
  prometheus.MustRegister(L2TokenWithdrawalNonce)
  prometheus.MustRegister(L2BasicTokenSequencerBalance)
  prometheus.MustRegister(L2StableTokenSequencerBalance)

  prometheus.MustRegister(L1BasicNFTAllowed)
  prometheus.MustRegister(L1BasicNFTSupply)
  prometheus.MustRegister(L1BasicNFTSequencerBalance)
  prometheus.MustRegister(L1BasicNFTBridgeBalance)
  prometheus.MustRegister(L1SpecialNFTAllowed)
  prometheus.MustRegister(L1SpecialNFTSupply)
  prometheus.MustRegister(L1SpecialNFTSequencerBalance)
  prometheus.MustRegister(L1SpecialNFTBridgeBalance)
  prometheus.MustRegister(L2BasicNFTAllowed)
  prometheus.MustRegister(L2BasicNFTSupply)
  prometheus.MustRegister(L2SpecialNFTAllowed)
  prometheus.MustRegister(L2BasicNFTSequencerBalance)
  prometheus.MustRegister(L2SpecialNFTSequencerBalance)
}

func (p *SmartContractMetricExporter) Start() error {
  log.Println("Starting Smart Contract Metric Exporter...")

  go func() {
    for {
      log.Println("Updating smart contract metrics from ", p.L1Comms.L1ContractAddressConfig.TxStorageContractAddress.String())
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

      l1BlockHeight, err := p.L1Comms.L1Client.BlockNumber(context.Background())
      if err != nil {
        log.Fatalf("Failed to get L1 block height: %v", err)
      }
      L1BlockHeight.Set(float64(l1BlockHeight))

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


      basicTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.Erc20Address)
      if err != nil {
        log.Fatalf("Failed to get basic token allowed: %v", err)
      }
      L1BasicTokenAllowed.Set(float64(basicTokenAllowed))

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

      log.Println("Getting token balance for token bridge account ", p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
      basicTokenBalance, err := p.ERC20Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
      if err != nil {
        log.Fatalf("Failed to get basic token balance: %v", err)
      }
      L1BasicTokenBridgeBalance.Set(float64(basicTokenBalance.Int64()))

      stableTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.StableErc20Address)
      if err != nil {
        log.Fatalf("Failed to get stable token allowed: %v", err)
      }
      L1StableTokenAllowed.Set(float64(stableTokenAllowed))

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

      log.Println("Getting token balance for token bridge account ", p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
      stableTokenBalance, err := p.StableERC20Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
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

      l2BasicTokenSupply, err := p.L2ERC20Contract.TotalSupply(nil)
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
      
      l2StableTokenSupply, err := p.L2StableERC20Contract.TotalSupply(nil)
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

      basicNFTokenAllowed, err := p.L1Comms.L1Contracts.TokenBridgeContract.AllowedTokens(nil, p.TokenAddresses.Erc721Address)
      if err != nil {
        log.Fatalf("Failed to get basic NFT allowed: %v", err)
      }
      L1BasicNFTAllowed.Set(float64(basicNFTokenAllowed))

      basicNFTokenSupply, err := p.ERC721Contract.TotalSupply(nil)
      if err != nil {
        log.Fatalf("Failed to get basic NFT token supply: %v", err)
      }
      L1BasicNFTSupply.Set(float64(basicNFTokenSupply.Int64()))

      basicNFTokenSequencerBalance, err := p.ERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
      if err != nil {
        log.Fatalf("Failed to get basic NFT token sequencer balance: %v", err)
      }
      L1BasicNFTSequencerBalance.Set(float64(basicNFTokenSequencerBalance.Int64()))

      basicNFTBridgeBalance, err := p.ERC721Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
      if err != nil {
        log.Fatalf("Failed to get basic NFT token bridge balance: %v", err)
      }
      L1BasicNFTBridgeBalance.Set(float64(basicNFTBridgeBalance.Int64()))

      specialNFTTokenName, err := p.SpecialERC721Contract.Name(nil)
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

      specialNFTSequncerBalance, err := p.SpecialERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
      if err != nil {
        log.Fatalf("Failed to get special NFT token sequencer balance: %v", err)
      }
      L2SpecialNFTSequencerBalance.Set(float64(specialNFTSequncerBalance.Int64()))

      specialNFTBridgeBalance, err := p.SpecialERC721Contract.BalanceOf(nil, p.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress)
      if err != nil {
        log.Fatalf("Failed to get special NFT token bridge balance: %v", err)
      }
      L1SpecialNFTBridgeBalance.Set(float64(specialNFTBridgeBalance.Int64()))

      specialNFTSequncerBalance, err = p.SpecialERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
      if err != nil {
        log.Fatalf("Failed to get special NFT token sequencer balance: %v", err)
      }
      L1SpecialNFTSequencerBalance.Set(float64(specialNFTSequncerBalance.Int64()))

      l2BasicNFTokenAllowed, err := p.L2Comms.L2Contracts.L2TokenBridgeContract.GetAllowedToken(nil, p.TokenAddresses.Erc721Address)
      if err != nil {
        log.Fatalf("Failed to get basic NFT allowed: %v", err)
      }
      if l2BasicNFTokenAllowed {
        L2BasicNFTAllowed.Set(float64(1))
      } else {
        L2BasicNFTAllowed.Set(float64(0))
      }

      l2BasicNFTSupply, err := p.L2ERC721Contract.TotalSupply(nil)
      if err != nil {
        log.Fatalf("Failed to get basic NFT token supply: %v", err)
      }
      L2BasicNFTSupply.Set(float64(l2BasicNFTSupply.Int64()))

      l2BasicNFTSequencerBalance, err := p.L2ERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
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

      l2SpecialNFTSequncerBalance, err := p.L2SpecialERC721Contract.BalanceOf(nil, l2utils.GetSequencer())
      if err != nil {
        log.Fatalf("Failed to get special NFT token sequencer balance: %v", err)
      }
      L2SpecialNFTSequencerBalance.Set(float64(l2SpecialNFTSequncerBalance.Int64()))

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
