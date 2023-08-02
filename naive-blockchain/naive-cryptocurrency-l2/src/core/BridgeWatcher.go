package core

import (
	"context"
	"log"
	"math/big"
	"time"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type BridgeWatcher struct {
  L1BlockNumber uint64 // TODO: Store in file to prevent needing to rewatch all blocks
  L2BlockNumber uint64

  L1Comms *l2utils.L1Comms
  Config *l2config.NodeBaseConfig
}

func NewBridgeWatcher(l1Comms *l2utils.L1Comms, config *l2config.NodeBaseConfig) *BridgeWatcher {
  return &BridgeWatcher{
    L1Comms: l1Comms,
    Config: config,
  }
}

func (bw *BridgeWatcher) WatchL1() error {
  latestBlockNumber, err := bw.L1Comms.L1Client.BlockNumber(context.Background())
  if err != nil {
    log.Fatalf("Failed to get latest block number: %v", err)
    return err
  }

  log.Printf("Latest L1 block number: %v   %v", latestBlockNumber, bw.L1BlockNumber)
  if latestBlockNumber <= bw.L1BlockNumber {
    time.Sleep(1 * time.Second)
    return err
  }

  for i := bw.L1BlockNumber + 1; i <= latestBlockNumber; i++ {
    log.Printf("Watching L1 block %v", i)
    newBlock, err := bw.L1Comms.L1Client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
    if err != nil {
      log.Fatalf("Failed to get block: %v", err)
      return err
    }

    if newBlock == nil {
      continue
    }

    for _, tx := range newBlock.Transactions() {
      receipt, err := bw.L1Comms.L1Client.TransactionReceipt(context.Background(), tx.Hash())
      if err != nil {
        log.Fatalf("Failed to get receipt: %v", err)
        return err
      }

      receipt_logs := l2utils.ReceiptLogsWithEvent(receipt, crypto.Keccak256Hash([]byte("EthDeposited(uint256,address,uint256)")).Bytes())
      for _, receipt_log := range receipt_logs {
        if common.HexToAddress(receipt_log.Address.Hex()) == bw.L1Comms.L1ContractAddressConfig.BridgeContractAddress {
          bridgeDep, err := bw.L1Comms.L1Contracts.BridgeContract.ParseEthDeposited(*receipt_log)
          if err != nil {
            log.Fatalf("Failed to parse deposit: %v", err)
            return err
          }

          log.Printf("L1 Deposit: %v", bridgeDep)
          transactOpts, err := l2utils.CreateTransactOpts(accounts.Account{Address: l2utils.GetSequencer()}, big.NewInt(int64(bw.Config.L2ChainID)))
          if err != nil {
            log.Fatalf("Failed to create transact opts: %v", err)
            return err
          }

          l2Comms, err := l2utils.NewL2Comms(bw.Config.DataDir + "/naive-sequencer.ipc", bw.Config.Contracts, big.NewInt(int64(bw.Config.L2ChainID)), l2utils.GetDefaultL2TransactionConfig())
          if err != nil {
            log.Fatalf("Failed to create L2 Comms: %v", err)
            return err
          }

          currDepositNonce, err := l2Comms.L2Contracts.L2BridgeContract.GetEthDepositNonce(&bind.CallOpts{})
          if err != nil {
            log.Fatalf("Failed to get deposit nonce: %v", err)
            return err
          }

          if currDepositNonce.Cmp(bridgeDep.Nonce) >= 0 {
            log.Printf("Skipping deposit nonce %v", bridgeDep.Nonce)
            continue
          }

          tx, err := l2Comms.L2Contracts.L2BridgeContract.DepositEth(transactOpts, bridgeDep.Addr, bridgeDep.Amount)
          if err != nil {
            log.Fatalf("Failed to deposit eth: %v", err)
            return err
          }

          log.Printf("L2 Deposit: %v", tx.Hash().Hex())
        }
      }

      receipt_logs = l2utils.ReceiptLogsWithEvent(receipt, crypto.Keccak256Hash([]byte("TokensDeposited(uint256,address,address,uint256)")).Bytes())
      log.Println("found x receipt logs: ", len(receipt_logs))
      for _, receipt_log := range receipt_logs { 
        //TODO: nonce check
        log.Printf("Watcher found Receipt log: %v", receipt_log)
        if common.HexToAddress(receipt_log.Address.Hex()) == bw.L1Comms.L1ContractAddressConfig.TokenBridgeContractAddress {
          tokenDep, err := bw.L1Comms.L1Contracts.TokenBridgeContract.ParseTokensDeposited(*receipt_log)
          if err != nil {
            log.Fatalf("Failed to parse deposit: %v", err)
            return err
          }

          log.Printf("L1 Token Deposit: %v", tokenDep)
          transactOpts, err := l2utils.CreateTransactOpts(accounts.Account{Address: l2utils.GetSequencer()}, big.NewInt(int64(bw.Config.L2ChainID)))
          if err != nil {
            log.Fatalf("Failed to create transact opts: %v", err)
            return err
          }

          l2Comms, err := l2utils.NewL2Comms(bw.Config.DataDir + "/naive-sequencer.ipc", bw.Config.Contracts, big.NewInt(int64(bw.Config.L2ChainID)), l2utils.GetDefaultL2TransactionConfig())
          if err != nil {
            log.Fatalf("Failed to create L2 Comms: %v", err)
            return err
          }

          currDepositNonce, err := l2Comms.L2Contracts.L2TokenBridgeContract.GetTokenDepositNonce(&bind.CallOpts{})
          if err != nil {
            log.Fatalf("Failed to get deposit nonce: %v", err)
            return err
          }

          if currDepositNonce.Cmp(tokenDep.Nonce) >= 0 {
            log.Printf("Skipping deposit nonce %v", tokenDep.Nonce)
            continue
          }

          log.Println("Using args: ", tokenDep.TokenAddress.Hex(), tokenDep.From.Hex(), tokenDep.Value)
          //TODO: Check if token is already deployed on L2 by allowedTokens
          tx, err := l2Comms.L2Contracts.L2TokenBridgeContract.MintTokens(transactOpts, tokenDep.TokenAddress, tokenDep.From, tokenDep.Value)
          if err != nil {
            log.Fatalf("Failed to deposit tokens: %v", err)
            return err
          }

          log.Printf("L2 Token Deposit: %v", tx.Hash().Hex())
        }
      }
    }

    log.Println("Done processing block: ", i)
    bw.L1BlockNumber = i
  }

  bw.L1BlockNumber = latestBlockNumber
  return nil
}

func (bw *BridgeWatcher) Watch() {
  runFunc := func() {
    for {
      err := bw.WatchL1()
      if err != nil {
        log.Fatalf("Failed to watch L1: %v", err)
        return
      }
    }
  }

  go runFunc()
}
