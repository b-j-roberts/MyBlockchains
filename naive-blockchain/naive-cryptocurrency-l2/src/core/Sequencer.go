package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"

	l2consensus "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/consensus"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

type Sequencer struct {
  L2Node  *Node
  Batcher *Batcher
  BridgeWatcher *BridgeWatcher

  MiningThreads int

  L2ContractAddresses l2utils.L2ContractAddressConfig
}

func NewSequencer(sequencerConfigFile string, posterAddress common.Address) (*Sequencer, error) {

  l2Node, err := NewNode(sequencerConfigFile)
  if err != nil {
    return nil, fmt.Errorf("failed to create l2 node: %v", err)
  }

  l2utils.SetSequencer(posterAddress) // Set the sequencer address to L1 posting address

  batcherConfig := &BatcherConfig{
    PosterAddress: posterAddress,
    BatchSize: 10,
    MaxBatchTimeMinutes: 1,
    NodeConfig: l2Node.Config,
  }

  batcher := NewBatcher(l2Node.L2Blockchain, batcherConfig)

  return &Sequencer{
    L2Node:   l2Node,
    Batcher:   batcher,
    BridgeWatcher: NewBridgeWatcher(batcher.L1Comms, l2Node.Config),
    MiningThreads: l2Node.Config.MiningThreads,
  }, nil
}

func (sequencer *Sequencer) Start() error {
  if err := sequencer.L2Node.Start(); err != nil {
    return fmt.Errorf("failed to start l2 node: %v", err)
  }

  // Setup mining signer
  eb, err := sequencer.L2Node.Eth.Etherbase()
  if err != nil {
    return fmt.Errorf("failed to get etherbase: %v", err)
  }
  var consensus *l2consensus.L2Consensus
  if c, ok := sequencer.L2Node.Eth.Engine().(*l2consensus.L2Consensus); ok {
    consensus = c
  } else {
    return fmt.Errorf("failed to get consensus 1")
  }
  if consensus == nil {
    return fmt.Errorf("failed to get consensus 2")
  }
  wallet, err := sequencer.L2Node.Eth.AccountManager().Find(accounts.Account{Address: eb})
  if wallet == nil || err != nil {
    return fmt.Errorf("failed to find account: %v", err)
  }
  consensus.Authorize(eb, wallet.SignData)


  if err := sequencer.L2Node.Eth.APIBackend.StartMining(sequencer.MiningThreads); err != nil {
    return fmt.Errorf("failed to start mining: %v", err)
  }

  if err := sequencer.Batcher.Start(); err != nil {
    return fmt.Errorf("failed to start batcher: %v", err)
  }
  

  err = sequencer.L2Node.Eth.StartNaive()
  if err != nil {
    return fmt.Errorf("failed to start eth: %v", err)
  }

  sequencer.BridgeWatcher.Watch()

  return nil
}

func (sequencer *Sequencer) Stop() {
  sequencer.L2Node.Stop()
}
