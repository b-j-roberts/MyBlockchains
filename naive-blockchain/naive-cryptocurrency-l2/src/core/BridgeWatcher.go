package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type BridgeWatcher struct {
  L1BridgeAddress common.Address
  L2BridgeAddress common.Address
  L1TokenBridgeAddress common.Address
  L2TokenBridgeAddress common.Address

  L1BlockNumber uint64 // TODO: Store in file to prevent needing to rewatch all blocks
  L2BlockNumber uint64

  L1Comms *l2utils.L1Comms

  L2ChainID int64
}

func NewBridgeWatcher(l1BridgeAddress common.Address, l2BridgeAddress common.Address, L1TokenBridgeAddress common.Address, l2TokenBridgeAddress common.Address, l1Comms *l2utils.L1Comms, l2ChainID int64) *BridgeWatcher {
  return &BridgeWatcher{
    L1BridgeAddress: l1BridgeAddress,
    L2BridgeAddress: l2BridgeAddress,
    L1TokenBridgeAddress: L1TokenBridgeAddress,
    L2TokenBridgeAddress: l2TokenBridgeAddress,
    L1Comms: l1Comms,
    L2ChainID: l2ChainID,
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
        if common.HexToAddress(receipt_log.Address.Hex()) == bw.L1Comms.BridgeContractAddress {
          bridgeDep, err := bw.L1Comms.BridgeContract.ParseEthDeposited(*receipt_log)
          if err != nil {
            log.Fatalf("Failed to parse deposit: %v", err)
            return err
          }

          log.Printf("L1 Deposit: %v", bridgeDep)
          transactOpts, err := l2utils.CreateTransactOpts(accounts.Account{Address: l2utils.GetSequencer()}, big.NewInt(bw.L2ChainID))
          if err != nil {
            log.Fatalf("Failed to create transact opts: %v", err)
            return err
          }

          ipcFile := "/home/brandon/naive-sequencer-data/naive-sequencer.ipc"
          rpcIPC, err := rpc.DialIPC(context.Background(), ipcFile)
          if err != nil {
            log.Fatalf("Failed to dial ipc: %v", err)
            return err
          }
  
          // l2 bridge address file is json containing field address
          l2BridgeAddressFile := "/home/brandon/naive-sequencer-data/l2-bridge-address.txt"
          l2BridgeAddressBytes, err := ioutil.ReadFile(l2BridgeAddressFile)
          if err != nil {
            log.Fatalf("Failed to read l2 bridge address file: %v", err)
            return err
          }
          var l2BridgeAddressJSONMap map[string]interface{}
          err = json.Unmarshal(l2BridgeAddressBytes, &l2BridgeAddressJSONMap)
          if err != nil {
            log.Fatalf("Failed to unmarshal l2 bridge address json: %v", err)
            return err
          }
          l2BridgeAddress := common.HexToAddress(l2BridgeAddressJSONMap["address"].(string))

          l2TokenBridgeAddressFile := "/home/brandon/naive-sequencer-data/l2-token-bridge-address.txt"
          l2TokenBridgeAddressBytes, err := ioutil.ReadFile(l2TokenBridgeAddressFile)
          if err != nil {
            log.Fatalf("Failed to read l2 token bridge address file: %v", err)
            return err
          }
          var l2TokenBridgeAddressJSONMap map[string]interface{}
          err = json.Unmarshal(l2TokenBridgeAddressBytes, &l2TokenBridgeAddressJSONMap)
          if err != nil {
            log.Fatalf("Failed to unmarshal l2 token bridge address json: %v", err)
            return err
          }
          l2TokenBridgeAddress := common.HexToAddress(l2TokenBridgeAddressJSONMap["address"].(string))
  
          backend := ethclient.NewClient(rpcIPC)
          l2Comms, err := l2utils.NewL2Comms(l2BridgeAddress, l2TokenBridgeAddress, big.NewInt(bw.L2ChainID), backend, l2utils.GetDefaultL2TransactionConfig())
          if err != nil {
            log.Fatalf("Failed to create L2 Comms: %v", err)
            return err
          }

          currDepositNonce, err := l2Comms.L2BridgeContract.GetEthDepositNonce(&bind.CallOpts{})
          if err != nil {
            log.Fatalf("Failed to get deposit nonce: %v", err)
            return err
          }

          if currDepositNonce.Cmp(bridgeDep.Nonce) >= 0 {
            log.Printf("Skipping deposit nonce %v", bridgeDep.Nonce)
            continue
          }

          log.Println("Calling deposit eth on L2 Bridge Address : ", l2BridgeAddress.Hex())
          tx, err := l2Comms.L2BridgeContract.DepositEth(transactOpts, bridgeDep.Addr, bridgeDep.Amount)
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
        log.Printf("Checking if address matches token bridge address: %v    %v", receipt_log.Address.Hex(), bw.L1Comms.TokenBridgeContractAddress)
        if common.HexToAddress(receipt_log.Address.Hex()) == bw.L1Comms.TokenBridgeContractAddress {
          tokenDep, err := bw.L1Comms.TokenBridgeContract.ParseTokensDeposited(*receipt_log)
          if err != nil {
            log.Fatalf("Failed to parse deposit: %v", err)
            return err
          }

          log.Printf("L1 Token Deposit: %v", tokenDep)
          transactOpts, err := l2utils.CreateTransactOpts(accounts.Account{Address: l2utils.GetSequencer()}, big.NewInt(bw.L2ChainID))
          if err != nil {
            log.Fatalf("Failed to create transact opts: %v", err)
            return err
          }

          ipcFile := "/home/brandon/naive-sequencer-data/naive-sequencer.ipc"
          rpcIPC, err := rpc.DialIPC(context.Background(), ipcFile)
          if err != nil {
            log.Fatalf("Failed to dial ipc: %v", err)
            return err
          }

          l2BridgeAddressFile := "/home/brandon/naive-sequencer-data/l2-bridge-address.txt"
          l2BridgeAddressBytes, err := ioutil.ReadFile(l2BridgeAddressFile)
          if err != nil {
            log.Fatalf("Failed to read l2 bridge address file: %v", err)
            return err
          }
          var l2BridgeAddressJSONMap map[string]interface{}
          err = json.Unmarshal(l2BridgeAddressBytes, &l2BridgeAddressJSONMap)
          if err != nil {
            log.Fatalf("Failed to unmarshal l2 bridge address json: %v", err)
            return err
          }
          l2BridgeAddress := common.HexToAddress(l2BridgeAddressJSONMap["address"].(string))

          // l2 token bridge address file is json containing field address
          l2TokenBridgeAddressFile := "/home/brandon/naive-sequencer-data/l2-token-bridge-address.txt"
          l2TokenBridgeAddressBytes, err := ioutil.ReadFile(l2TokenBridgeAddressFile)
          if err != nil {
            log.Fatalf("Failed to read l2 token bridge address file: %v", err)
            return err
          }
          var l2TokenBridgeAddressJSONMap map[string]interface{}
          err = json.Unmarshal(l2TokenBridgeAddressBytes, &l2TokenBridgeAddressJSONMap)
          if err != nil {
            log.Fatalf("Failed to unmarshal l2 token bridge address json: %v", err)
            return err
          }
          l2TokenBridgeAddress := common.HexToAddress(l2TokenBridgeAddressJSONMap["address"].(string))

          backend := ethclient.NewClient(rpcIPC)
          l2Comms, err := l2utils.NewL2Comms(l2BridgeAddress, l2TokenBridgeAddress, big.NewInt(bw.L2ChainID), backend, l2utils.GetDefaultL2TransactionConfig())
          if err != nil {
            log.Fatalf("Failed to create L2 Comms: %v", err)
            return err
          }

          currDepositNonce, err := l2Comms.L2TokenBridgeContract.GetTokenDepositNonce(&bind.CallOpts{})
          if err != nil {
            log.Fatalf("Failed to get deposit nonce: %v", err)
            return err
          }

          if currDepositNonce.Cmp(tokenDep.Nonce) >= 0 {
            log.Printf("Skipping deposit nonce %v", tokenDep.Nonce)
            continue
          }

          log.Println("Calling mint tokens on L2 Token Bridge Address : ", l2TokenBridgeAddress.Hex())
          log.Println("Using args: ", tokenDep.TokenAddress.Hex(), tokenDep.From.Hex(), tokenDep.Amount)
          //TODO: Check if token is already deployed on L2by allowedTOkens
          tx, err := l2Comms.L2TokenBridgeContract.MintTokens(transactOpts, tokenDep.TokenAddress, tokenDep.From, tokenDep.Amount)
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
