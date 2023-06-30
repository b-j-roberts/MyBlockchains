package core

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

type Node struct {
  ChainDb ethdb.Database
  Node    *node.Node
  L2Blockchain *core.BlockChain
  Engine  consensus.Engine
  Eth     *eth.Ethereum
}

func NewNode(node *node.Node, chainDb ethdb.Database, l2Blockchain *core.BlockChain, engine consensus.Engine, config *ethconfig.Config, l1BridgeConfig *eth.L1BridgeConfig) (*Node, error) {
  txPool := txpool.NewTxPool(config.TxPool, l2Blockchain.Config(), l2Blockchain)
  naive_eth := eth.NewNaiveEthereum(l2Blockchain, chainDb, node, config, txPool, engine, l1BridgeConfig)
  //TODO: Learn more about APIs & which to enable/disable based on public / ...?
  apis := eth.GetNaiveEthAPIs(naive_eth)
  apis = append(apis, engine.APIs(l2Blockchain)...)
  apis = append(apis, []rpc.API{
    {
      Namespace: "eth",
      Service:   eth.NewEthereumAPI(naive_eth),
    }, {
      Namespace: "admin",
      Service:   eth.NewAdminAPI(naive_eth),
    }, {
      Namespace: "debug",
      Service:   eth.NewDebugAPI(naive_eth),
    },
  }...)
  node.RegisterAPIs(apis)
  node.RegisterProtocols(naive_eth.Protocols())
  node.RegisterLifecycle(naive_eth)

  //TODO: APIs / RPC
  return &Node{
    ChainDb: chainDb,
    Node:    node,
    L2Blockchain: l2Blockchain,
    Engine:  engine,
    Eth:     naive_eth,
  }, nil
}

func (node *Node) Start() error {
  backends := node.Node.AccountManager().Backends(keystore.KeyStoreType)
  if len(backends) == 0 {
    return fmt.Errorf("no key store backends found")
  }
  ks := backends[0].(*keystore.KeyStore)
  address := ks.Accounts()[0].Address.Hex()
  log.Println("Address is", address)

  account, err := utils.MakeAddress(ks, address)
  err = ks.Unlock(account, "password") //TODO: Hardcode
  if err != nil {
    return fmt.Errorf("failed to unlock account: %v", err)
  } else {
    log.Println("Unlocked account", ks.Accounts()[0].Address.Hex())
  }

  err = node.Node.Start()
  if err != nil {
    return fmt.Errorf("failed to start node stack: %v", err)
  }

  return nil
}
