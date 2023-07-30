package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"

	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

type Sequencer struct {
  L2Node  *Node
  Batcher *Batcher
  BridgeWatcher *BridgeWatcher

  MiningThreads int

  L2BridgeAddress common.Address
  L2TokenBridgeAddress common.Address
}

func NewSequencer(node *node.Node, chainDb ethdb.Database, l2Blockchain *core.BlockChain, engine consensus.Engine, config *ethconfig.Config, l1ContractAddress common.Address, l1BridgeAddress common.Address, l1TokenBridgeAddress common.Address, posterAddress common.Address, l1Url string, l1ChainId int, l2ChainId int, miningThreads int, contractsPath string) (*Sequencer, error) {

  l1BridgeConfig := &eth.L1BridgeConfig{
    L1BridgeAddress: l1BridgeAddress,
    L1TokenBridgeAddress: l1TokenBridgeAddress,
    L1BridgeUrl: l1Url,
    SequencerAddr: posterAddress,
  }

  l2Node, err := NewNode(node, chainDb, l2Blockchain, engine, config, l1BridgeConfig)
  if err != nil {
    return nil, fmt.Errorf("failed to create l2 node: %v", err)
  }

  l2utils.SetSequencer(l1BridgeConfig.SequencerAddr)

  batcherConfig := &BatcherConfig{
    L1NodeUrl: l1Url,
    L1ContractAddress: l1ContractAddress,
    L1ChainId: l1ChainId,
    PosterAddress: posterAddress,
    BatchSize: 10,
    MaxBatchTimeMinutes: 1,
    L1BridgeAddress: l1BridgeAddress,
    L1TokenBridgeAddress: l1TokenBridgeAddress,
    L2IPCPath: node.DataDir() + "/naive-sequencer.ipc",
    ContractsPath: contractsPath,
  }

  batcher := NewBatcher(l2Blockchain, batcherConfig)

  //TODO: APIs / RPC
  return &Sequencer{
    L2Node:   l2Node,
    Batcher:   batcher,
    BridgeWatcher: NewBridgeWatcher(l1BridgeAddress, common.HexToAddress("0x0"), l1TokenBridgeAddress, common.HexToAddress("0x0"), batcher.L1Comms, int64(l2ChainId), node.DataDir() + "/naive-sequencer.ipc", contractsPath),
    MiningThreads: miningThreads,
  }, nil
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

  sequencer.BridgeWatcher.Watch()

  return nil
}

func (sequencer *Sequencer) Stop() {
  sequencer.L2Node.Stop()
}
