package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"

	l2config "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/config"
	l2consensus "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/consensus"
)

type Node struct {
  ChainDb ethdb.Database
  Node    *node.Node
  L2Blockchain *core.BlockChain
  Engine  consensus.Engine
  Eth     *eth.Ethereum
  Config  *l2config.NodeBaseConfig
}

func CreateL2ConsensusEngine(config *l2consensus.L2ConsensusConfig, db ethdb.Database) (consensus.Engine, error) {
  engine := l2consensus.New(config, db)
  return engine, nil
}

func NewNode(rpcConfigFile string) (*Node, error) {
  nodeBaseConfig, err := l2config.LoadNodeBaseConfig(rpcConfigFile)
  if err != nil {
    return nil, fmt.Errorf("failed to load node base config: %v", err)
  }

  // Setup Geth node/node
  node, err := node.New(l2config.NodeConfig(nodeBaseConfig))
  if err != nil {
    return nil, fmt.Errorf("failed to create node: %v", err)
  }

  // Add Keystore to backend for unlocking accounts later on
  am := node.AccountManager()
  am.AddBackend(keystore.NewKeyStore(node.KeyStoreDir(), keystore.StandardScryptN, keystore.StandardScryptP))
  backends := am.Backends(keystore.KeyStoreType)
  if len(backends) == 0 {
    return nil, fmt.Errorf("no key store backends found")
  }
  ks := backends[0].(*keystore.KeyStore)

  if len(ks.Accounts()) == 0 {
    return nil, fmt.Errorf("no accounts found in key store")
  }
  address := ks.Accounts()[0].Address//TODO: Is this just posterAddress?

  // Setup Database
  //TODO: chainDb more research on args
  // Handles, Persistent Chain Dir, & Ancient from nitro/cmd/conf/database.go
  // Caching from arbnode/execution/blockchain.go DatabaseCache
  // Namespace is prefix for metrics
  // Open rawdb from geth/core with ancients freezer & configs from arbitrum chainDb ( Disk based db )
  chainDb, err := node.OpenDatabaseWithFreezer("l2-chain", 2048, 512, "", "naive-l2/chaindb", false)
  if err != nil {
    return nil, fmt.Errorf("failed to open chain database: %v", err)
  }

  // Setup Genesis
  file, err := os.Open(nodeBaseConfig.Genesis)
  if err != nil {
    return nil, fmt.Errorf("failed to open genesis file: %v", err)
  }
  defer file.Close()

  genesis := new(core.Genesis)
  if err := json.NewDecoder(file).Decode(genesis); err != nil {
    return nil, fmt.Errorf("invalid genesis file: %v", err)
  }
  trieDb := trie.NewDatabaseWithConfig(chainDb, &trie.Config{Preimages: true})
  _, _, err = core.SetupGenesisBlock(chainDb, trieDb, genesis)
  if err != nil {
    return nil, fmt.Errorf("failed to setup genesis block: %v", err)
  }

  // Setup Consensus Engine
  ethConfig := l2config.EthConfig(address)
  cliqueConfig, err := core.LoadCliqueConfig(chainDb, genesis)
  l2ConsenusConfig := l2consensus.L2ConsensusConfig{
    CliqueConfig: cliqueConfig,
    ContractsPath: nodeBaseConfig.Contracts,
  }
  if err != nil {
    return nil, fmt.Errorf("failed to load clique config: %v", err)
  }

  engine, err := CreateL2ConsensusEngine(&l2ConsenusConfig, chainDb)
  if err != nil {
    return nil, fmt.Errorf("failed to create consensus engine: %v", err)
  }

  // Setup L2 Blockchain
  //TODO: l2blockchain more research on args
  var overrides core.ChainOverrides
  vmConfig := vm.Config{EnablePreimageRecording: false}
  txLookupLimit := uint64(31536000) // 1 year at 1 block per second
  l2Blockchain, err := core.NewBlockChain(chainDb, l2config.DefaultCacheConfigFor(node, false), genesis, &overrides, engine, vmConfig, l2config.ShouldPreserveFalse, &txLookupLimit)
  if err != nil {
    return nil, fmt.Errorf("failed to create L2 blockchain: %v", err)
  }

  txPool := txpool.NewTxPool(ethConfig.TxPool, l2Blockchain.Config(), l2Blockchain)
  naive_eth := eth.NewNaiveEthereum(l2Blockchain, chainDb, node, ethConfig, txPool, engine)
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
    Config:  nodeBaseConfig,
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
  if err != nil {
    return fmt.Errorf("failed to make address: %v", err)
  }

  // Load password from environment variable NODE_PASS
  node_pass, exists := os.LookupEnv("NODE_PASS") //TODO: Think about security of this
  if !exists {
    node_pass = "password"
  }
  err = ks.Unlock(account, node_pass)
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

func (node *Node) Stop() error {
  node.L2Blockchain.Stop()

  return nil
}
