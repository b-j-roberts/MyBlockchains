package main

import (
	"github.com/prometheus/client_golang/prometheus"
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
