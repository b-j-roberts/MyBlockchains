package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

type Sequencer struct {
  L2Node  *Node
  Batcher *Batcher

  MiningThreads int

  L1BridgeBlockIdx uint64 //TODO: BridgeWatcher type
  L2BridgeAddress common.Address
}

func NewSequencer(node *node.Node, chainDb ethdb.Database, l2Blockchain *core.BlockChain, engine consensus.Engine, config *ethconfig.Config, l1ContractAddress common.Address, l1BridgeAddress common.Address, posterAddress common.Address, l1Host string, l1Port int, l1ChainId int, miningThreads int) (*Sequencer, error) {

  l1BridgeConfig := &eth.L1BridgeConfig{
    L1BridgeAddress: l1BridgeAddress,
    L1BridgeUrl: fmt.Sprintf("http://%s:%d", l1Host, l1Port),
    SequencerAddr: posterAddress,
  }

  l2Node, err := NewNode(node, chainDb, l2Blockchain, engine, config, l1BridgeConfig)
  if err != nil {
    return nil, fmt.Errorf("failed to create l2 node: %v", err)
  }

  l2utils.SetSequencer(l1BridgeConfig.SequencerAddr)

  l1Url := fmt.Sprintf("http://%s:%d", l1Host, l1Port)
  batcherConfig := &BatcherConfig{
    L1NodeUrl: l1Url,
    L1ContractAddress: l1ContractAddress,
    L1ChainId: l1ChainId,
    PosterAddress: posterAddress,
    BatchSize: 10,
    MaxBatchTimeMinutes: 1,
    L1BridgeAddress: l1BridgeAddress,
  }

  //TODO: APIs / RPC
  return &Sequencer{
    L2Node:   l2Node,
    Batcher:   NewBatcher(l2Blockchain, batcherConfig),
    MiningThreads: miningThreads,

    L1BridgeBlockIdx: 0, //TODO: Get this from data to resume?
  }, nil
}

func (sequencer *Sequencer) L1BridgeWatcher() error {
  runFunc := func() {
  for {
    latestBlockIdx, err := sequencer.Batcher.L1Comms.L1Client.BlockNumber(context.Background()) //TODO: Check this is finalized
    if err != nil {
      log.Fatalf("Failed to get block number: %v", err)
      return
    }
    log.Printf("Watcher : Latest block idx: %v  -- Current block idx: %v", latestBlockIdx, sequencer.L1BridgeBlockIdx)

    if latestBlockIdx <= sequencer.L1BridgeBlockIdx {
      // Sleep for 100 milliseconds
      time.Sleep(1000 * time.Millisecond)
      continue
    }

    for i := sequencer.L1BridgeBlockIdx + 1; i <= latestBlockIdx; i++ {
      log.Printf("Watcher : Checking block %v", i)
      newblock, err := sequencer.Batcher.L1Comms.L1Client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
      if err != nil {
        log.Fatalf("Failed to get block: %v", err)
        return
      }

      if newblock == nil {
        break
      }

      for _, tx := range newblock.Transactions() {
        log.Printf("Watcher : Checking tx %v", tx.Hash().Hex())
        receipt, err := sequencer.Batcher.L1Comms.L1Client.TransactionReceipt(context.Background(), tx.Hash())
        if err != nil {
          log.Fatalf("Failed to get transaction receipt: %v", err)
          return 
        }

        for _, receipt_log := range receipt.Logs {
          eventSignature := crypto.Keccak256Hash([]byte("EthDeposited(uint256,address,uint256)"))
          if bytes.Equal(receipt_log.Topics[0].Bytes(), eventSignature.Bytes()) && common.HexToAddress(receipt_log.Address.Hex()) == sequencer.Batcher.L1Comms.BridgeContractAddress {
            bridgeDep, err := sequencer.Batcher.L1Comms.BridgeContract.ParseEthDeposited(*receipt_log)
            if err != nil {
              log.Fatalf("Failed to parse eth deposited event: %v", err)
              return }

            //TODO: Check nonce
            log.Printf("Got deposit: %v    ---- %v %v %v", bridgeDep, bridgeDep.Nonce, bridgeDep.Addr, bridgeDep.Amount)
            transactOpts, err := l2utils.CreateTransactOpts(accounts.Account{Address: l2utils.GetSequencer()}, big.NewInt(515)) //TODO: Hardcoded
            if err != nil {
              log.Fatalf("Failed to create transact opts: %v", err)
              return
            }

            ipcFile := sequencer.L2Node.Node.DataDir() + "/naive-sequencer.ipc"
            rpcIPC, err := rpc.DialIPC(context.Background(), ipcFile)
            if err != nil {
              log.Fatalf("Failed to dial ipc: %v", err)
              return
            }

            // l2 bridge address file is json containing field address
            l2BridgeAddressFile := sequencer.L2Node.Node.DataDir() + "/l2-bridge-address.txt"
            l2BridgeAddressBytes, err := ioutil.ReadFile(l2BridgeAddressFile)
            if err != nil {
              log.Fatalf("Failed to read l2 bridge address file: %v", err)
              return
            }
            var l2BridgeAddressJSONMap map[string]interface{}
            err = json.Unmarshal(l2BridgeAddressBytes, &l2BridgeAddressJSONMap)
            if err != nil {
              log.Fatalf("Failed to unmarshal l2 bridge address json: %v", err)
              return
            }
            l2BridgeAddress := common.HexToAddress(l2BridgeAddressJSONMap["address"].(string))

            backend := ethclient.NewClient(rpcIPC)
            log.Println("Calling deposit eth on L2 Bridge Address : ", l2BridgeAddress.Hex())
            l2Comms, err := l2utils.NewL2Comms(l2BridgeAddress, backend)

            currDepositNonce, err := l2Comms.L2BridgeContract.GetEthDepositNonce(&bind.CallOpts{})
            if err != nil {
              log.Fatalf("Failed to get deposit nonce: %v", err)
              return
            }

            if currDepositNonce.Cmp(bridgeDep.Nonce) >= 0 {
              log.Printf("Skipping deposit nonce %v since it is already processed", bridgeDep.Nonce)
              continue
            }

            tx, err := l2Comms.L2BridgeContract.DepositEth(transactOpts, bridgeDep.Addr, bridgeDep.Amount)
            if err != nil {
              log.Fatalf("Failed to deposit eth: %v", err)
              return
            }

            log.Printf("Sent deposit tx to l2: %v", tx.Hash().Hex())
          }
        }
      }
    }

    sequencer.L1BridgeBlockIdx = latestBlockIdx
  }}

  go runFunc()
  return nil
}

func (sequencer *Sequencer) Start() error {
  if err := sequencer.L2Node.Start(); err != nil {
    return fmt.Errorf("failed to start l2 node: %v", err)
  }

  if err := sequencer.L2Node.Eth.APIBackend.StartMining(sequencer.MiningThreads); err != nil {
    return fmt.Errorf("failed to start mining: %v", err)
  }

  if err := sequencer.Batcher.Start(); err != nil {
    return fmt.Errorf("failed to start batcher: %v", err)
  }
  

  err := sequencer.L2Node.Eth.StartNaive()
  if err != nil {
    return fmt.Errorf("failed to start eth: %v", err)
  }

  //TODO: Wait for eth to be ready & setup l2 comms & things?

  if err := sequencer.L1BridgeWatcher(); err != nil {
    return fmt.Errorf("failed to start l1 bridge watcher: %v", err)
  }

  return nil
}

func (sequencer *Sequencer) Stop() {
  sequencer.L2Node.Stop()
}
