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

  MiningThreads int
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
  }

  //TODO: APIs / RPC
  return &Sequencer{
    L2Node:   l2Node,
    Batcher:   NewBatcher(l2Blockchain, batcherConfig),
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

  return nil
}
