package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"

	naive_utils "naive-l2/src/utils"
)

type Batcher struct {
  BatcherConfig *BatcherConfig

  L1Comms *naive_utils.L1Comms
  L2Blockchain *core.BlockChain
  BlockIdx     uint64

  PostedBlockIdx uint64
  LastPostTime time.Time

  TxBatch      []*types.Transaction
  BatchId int64
}

type BatcherConfig struct {
  L1NodeUrl string
  L1ContractAddress common.Address
  PosterAddress common.Address
  BatchSize int
  MaxBatchTimeMinutes int
}


func NewBatcher(l2Blockchain *core.BlockChain, batcherConfig *BatcherConfig) *Batcher {
  l1Comms, err := naive_utils.NewL1Comms(batcherConfig.L1NodeUrl, batcherConfig.L1ContractAddress)
  if err != nil {
    log.Fatalf("Error creating L1 comms: %s\n", err)
  }

  //TODO: Load info from L1
  return &Batcher{
    BatcherConfig: batcherConfig,
    L1Comms: l1Comms,
    L2Blockchain: l2Blockchain,
    BlockIdx:     0,
    PostedBlockIdx: 0,
    LastPostTime: time.Now(), //TODO
    TxBatch:      make([]*types.Transaction, 0),
    BatchId: 0,
  }
}

func (batcher *Batcher) PostBatch() error {
  if len(batcher.TxBatch) == 0 {
    return fmt.Errorf("no txs to post")
  }

  log.Printf("Posting batch of %d txs\n", len(batcher.TxBatch))

  //TODO: Compress transaction data
  transactionByteData := make([]byte, 0)
  for _, tx := range batcher.TxBatch {
    txBin, err := tx.MarshalBinary()
    if err != nil {
      return err
    }
    transactionByteData = append(transactionByteData, txBin...)
  }

  byteDataHash := sha256.Sum256(transactionByteData)// TODO: Use blockchain root

  err := batcher.L1Comms.PostBatch(transactionByteData, batcher.BatchId, byteDataHash, batcher.BatcherConfig.PosterAddress)
  if err != nil {
    return err
  }

  return nil
}

func (batcher *Batcher) Start() error {
  runFunc := func() {
  for {
    block := batcher.L2Blockchain.GetBlockByNumber(batcher.BlockIdx)
    if block != nil {
      for _, tx := range block.Transactions() {
        batcher.TxBatch = append(batcher.TxBatch, tx)
      }

      if len(batcher.TxBatch) > batcher.BatcherConfig.BatchSize ||
         (len(batcher.TxBatch) > 0 && batcher.BatcherConfig.MaxBatchTimeMinutes > 0 &&
          time.Since(batcher.LastPostTime) > time.Duration(batcher.BatcherConfig.MaxBatchTimeMinutes) * time.Minute && len(batcher.TxBatch) > 0) {
        err := batcher.PostBatch()
        if err != nil {
          panic(err)
        }
        batcher.LastPostTime = time.Now()
        batcher.PostedBlockIdx = batcher.BlockIdx
        batcher.TxBatch = make([]*types.Transaction, 0)
        batcher.BatchId++
      }
      batcher.BlockIdx++
    } else {
      //TODO: Only sleep if caught up
      time.Sleep(100 * time.Millisecond)
    }
  }}

  go runFunc()

  return nil
}

type Node struct {
  ChainDb ethdb.Database
  Node    *node.Node
  L2Blockchain *core.BlockChain
  Engine  consensus.Engine
  Eth     *eth.Ethereum
  Batcher *Batcher
}

func NewNode(node *node.Node, chainDb ethdb.Database, l2Blockchain *core.BlockChain, engine consensus.Engine, config *ethconfig.Config, l1ContractAddress common.Address, posterAddress common.Address, l1Host string, l1Port int) (*Node, error) {
  txPool := txpool.NewTxPool(config.TxPool, l2Blockchain.Config(), l2Blockchain)
  naive_eth := eth.NewNaiveEthereum(l2Blockchain, chainDb, node, config, txPool, engine)
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

  l1Url := fmt.Sprintf("http://%s:%d", l1Host, l1Port)

  batcherConfig := &BatcherConfig{
    L1NodeUrl: l1Url,
    L1ContractAddress: l1ContractAddress,
    PosterAddress: posterAddress,
    BatchSize: 10,
    MaxBatchTimeMinutes: 1,
  }

  //TODO: APIs / RPC
  return &Node{
    ChainDb: chainDb,
    Node:    node,
    L2Blockchain: l2Blockchain,
    Engine:  engine,
    Eth:     naive_eth,
    Batcher:   NewBatcher(l2Blockchain, batcherConfig), //TODO: Hardcode
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

  if err := node.Eth.APIBackend.StartMining(4); err != nil { //TODO: Hardcode
    return fmt.Errorf("failed to start mining: %v", err)
  }

  if err := node.Batcher.Start(); err != nil {
    return fmt.Errorf("failed to start batcher: %v", err)
  }

  err = node.Eth.StartNaive()
  if err != nil {
    return fmt.Errorf("failed to start eth: %v", err)
  }

  return nil
}

func CreateNaiveNode(dataDir string, httpHost string, httpPort int, httpModules string, l1Host string, l1Port int,
                     l1ContractAddress common.Address, posterAddress common.Address) (*Node, error) {
  // Function used to create Naive Node mimicing eth/backend.go:New for Ethereum Node object

  // Setup Geth node/node
  nodeConfig := NodeConfig(dataDir, httpHost, httpPort, httpModules)
  node, err := node.New(nodeConfig)
  if err != nil {
    return nil, fmt.Errorf("failed to create node: %v", err)
  }

  // Add Keystore to backend for unlocking accounts later on
  am := node.AccountManager()
  log.Println("KeyStoreDir is", node.KeyStoreDir())
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
  file, err := os.Open(nodeConfig.DataDir + "/genesis.json") //TODO: Hardcode                                                                         
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
  ethConfig := EthConfig(address)
  cliqueConfig, err := core.LoadCliqueConfig(chainDb, genesis)
  if err != nil {
    return nil, fmt.Errorf("failed to load clique config: %v", err)
  }

  engine := ethconfig.CreateConsensusEngine(node, &ethConfig.Ethash, cliqueConfig, ethConfig.Miner.Notify, ethConfig.Miner.Noverify, chainDb)

  // Setup L2 Blockchain
  //TODO: l2blockchain more research on args
  var overrides core.ChainOverrides
  vmConfig := vm.Config{EnablePreimageRecording: false}
  txLookupLimi := uint64(31536000) // 1 year at 1 block per second
  l2BlockChain, err := core.NewBlockChain(chainDb, DefaultCacheConfigFor(node, false), genesis, &overrides, engine, vmConfig, ShouldPreserveFalse, &txLookupLimi)
  if err != nil {
    return nil, fmt.Errorf("failed to create L2 blockchain: %v", err)
  }

  //TODO: naiveDb, err := node.OpenDatabase("naivedata", 0, 0, "", false)

  naiveNode, err := NewNode(node, chainDb, l2BlockChain, engine, ethConfig, l1ContractAddress, posterAddress, l1Host, l1Port)
  if err != nil {
    return nil, fmt.Errorf("failed to create naive node: %v", err)
  }


  return naiveNode, nil
}

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Starting sequencer...")

  osHomeDir, err := os.UserHomeDir()
  dataDir := flag.String("datadir", osHomeDir + "/naive-sequencer-data", "data directory for the database and keystore")
  httpHost := flag.String("httphost", "localhost", "HTTP-RPC server listening interface")
  httpPort := flag.Int("httpport", 5055, "HTTP-RPC server listening port")
  httpModules := flag.String("httpmodules", "personal,naive", "Comma separated list of API modules to enable on the HTTP-RPC interface")
  l1ContractAddress := flag.String("l1contract", "", "Address of the L1 contract")
  sequencerAddress := flag.String("sequencer", "", "Address of the sequencer on L1")
  sequencerKeystore := flag.String("sequencerkeystore", "", "Keystore file for the sequencer on L1")
  l1Host := flag.String("l1host", "localhost", "L1 HTTP-RPC server listening interface")
  l1Port := flag.Int("l1port", 8545, "L1 HTTP-RPC server listening port")
  flag.Parse()

  log.Println("Connecting to L1 contract at", *l1Host, *l1Port, "with address", *l1ContractAddress)
  naiveNode, err := CreateNaiveNode(*dataDir, *httpHost, *httpPort, *httpModules, *l1Host, *l1Port, common.HexToAddress(*l1ContractAddress), common.HexToAddress(*sequencerAddress))
  if err != nil {
    utils.Fatalf("Failed to create naive sequencer node: %v", err)
  }

  naiveNode.Batcher.L1Comms.RegisterL2Address(common.HexToAddress(*sequencerAddress), *sequencerKeystore)

  ////TODO: close dbs & stop blockchain defers
  fatalErrChan := make(chan error, 10)

  genesis, err := core.ReadGenesis(naiveNode.ChainDb)
  if err != nil {
    fatalErrChan <- err
  }

  err = naiveNode.Batcher.L1Comms.L2GenesisOnL1(genesis, common.HexToAddress(*sequencerAddress))
  if err != nil {
    fatalErrChan <- err
  }
  naiveNode.Batcher.BatchId = 1

  err = naiveNode.Start()
  if err != nil {
    fatalErrChan <- err
  }
  log.Println("Naive Sequencer Node started", naiveNode)

  sigint := make(chan os.Signal, 1)
  signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
  
  exitCode := 0
  select {
  case err := <-fatalErrChan:
    log.Println("shutting down due to fatal error", "err", err)
    defer log.Println("shut down due to fatal error", "err", err)
    exitCode = 1
  case <-sigint:
    log.Println("shutting down because of sigint")
  }
  
  // cause future ctrl+c's to panic
  close(sigint)

  // node stop&wait

  return exitCode
}
