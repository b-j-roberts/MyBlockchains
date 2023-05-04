package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"

	contracts "naive-l2/contracts/go"
)

// Don't preserve reorg'd out blocks
func shouldPreserveFalse(_ *types.Header) bool {
  return false                                      
}

type CachingConfig struct {
  Archive               bool          `koanf:"archive"`
  BlockCount            uint64        `koanf:"block-count"`
  BlockAge              time.Duration `koanf:"block-age"`
  TrieTimeLimit         time.Duration `koanf:"trie-time-limit"`
  TrieDirtyCache        int           `koanf:"trie-dirty-cache"`
  TrieCleanCache        int           `koanf:"trie-clean-cache"`
  SnapshotCache         int           `koanf:"snapshot-cache"`
  DatabaseCache         int           `koanf:"database-cache"`
  SnapshotRestoreMaxGas uint64        `koanf:"snapshot-restore-gas-limit"`
}

// See arbnode/execution/blockchain.go
var DefaultCachingConfig = CachingConfig{
  Archive:               false,
  BlockCount:            128,
  BlockAge:              30 * time.Minute,
  TrieTimeLimit:         time.Hour,
  TrieDirtyCache:        1024,
  TrieCleanCache:        600,
  SnapshotCache:         400,
  DatabaseCache:         2048,
  SnapshotRestoreMaxGas: 300_000_000_000,
}

func DefaultCacheConfigFor(stack *node.Node, archive bool) *core.CacheConfig {
  baseConf := ethconfig.Defaults
  //if archive {
  //  baseConf = ethconfig.ArchiveDefaults
  //}

  return &core.CacheConfig{
    TrieCleanLimit:        DefaultCachingConfig.TrieCleanCache,
    TrieCleanJournal:      stack.ResolvePath(baseConf.TrieCleanCacheJournal),
    TrieCleanRejournal:    baseConf.TrieCleanCacheRejournal,
    TrieCleanNoPrefetch:   baseConf.NoPrefetch,
    TrieDirtyLimit:        DefaultCachingConfig.TrieDirtyCache,
    TrieDirtyDisabled:     DefaultCachingConfig.Archive,
    TrieTimeLimit:         DefaultCachingConfig.TrieTimeLimit,
    SnapshotLimit:         DefaultCachingConfig.SnapshotCache,
    Preimages:             baseConf.Preimages,
  }
}

type L1InOut struct {
  RpcUrl string
  L1Client *ethclient.Client
  Contract *contracts.Contracts
  ContractAddress common.Address
  node *node.Node
}

func NewL1InOut(rpcUrl string, node *node.Node) *L1InOut {
  l1Comms := &L1InOut{
    RpcUrl: rpcUrl,
    node: node,
  }

  rawRpc, err := rpc.Dial(l1Comms.RpcUrl)
  if err != nil {
    log.Fatalf("Failed to connect to l1 node: %v", err)
    return nil
  }

  l1Comms.L1Client = ethclient.NewClient(rawRpc)

  // Read address from file
  addressFile, err := os.Open("/home/b-j-roberts/workspace/blockchain/my-chains/naive-blockchain/naive-cryptocurrency-l2/contracts/builds/contract-address.txt") //TODO: Hardcode
  addressStringBytes := make([]byte, 42) //TODO: Hardcode
  log.Println("Address length is", 42)
  addressFile.Read(addressStringBytes)
  addressString := string(addressStringBytes)
  log.Println("Address string is", addressString)
  l1Comms.ContractAddress = common.HexToAddress(addressString)

  if err != nil {
    log.Fatalf("Failed to open address file: %v", err)
    return nil
  }
  defer addressFile.Close()

  //backends := node.AccountManager().Backends(keystore.KeyStoreType)
  //if len(backends) == 0 {
  //  log.Fatalf("No key store backends found")
  //}
  //ks := backends[0].(*keystore.KeyStore)
  //address := ks.Accounts()[0].Address
        
  log.Println("Contract Address is", l1Comms.ContractAddress.Hex())
  
  l1Comms.Contract, err = contracts.NewContracts(l1Comms.ContractAddress, l1Comms.L1Client)
  if err != nil {
    log.Fatalf("Failed to instantiate contract: %v", err)
    return nil
  }

  batchCount, err := l1Comms.Contract.GetBatchCount(nil)
  if err != nil {
    log.Fatalf("Failed to get batch count: %v", err)
    return nil
  }

  log.Println("Batch number is", batchCount)

  return l1Comms
}

func (l1InOut *L1InOut) PostBatch(transactionByteData []byte, id int64, hash [32]byte) error {
  log.Println("Posting batch here", id, "with hash", hash)
  // Read address from json file under .address field
  var jsonObject map[string]interface{}
  jsonFile, err := os.Open("/home/b-j-roberts/workspace/blockchain/my-chains/eth-private-network/data/keystore/UTC--2023-05-04T04-33-35.152600358Z--d966e954a01644a89eeb9c70157d5a1c4410f31b") //TODO: Hardcode
  if err != nil {
    return fmt.Errorf("failed to open address.json: %v", err)
  }
  defer jsonFile.Close()

  jsonParser := json.NewDecoder(jsonFile)
  if err = jsonParser.Decode(&jsonObject); err != nil {
    return fmt.Errorf("failed to decode address.json: %v", err)
  }
  addressString := jsonObject["address"].(string)
  log.Println("Address string is", addressString)

  batchCount, err := l1InOut.Contract.GetBatchCount(nil)  
  if err != nil {  
    log.Fatalf("Failed to get batch count: %v", err)  
    return nil  
  }  
  log.Println("Batch number is", batchCount)

  SignerFunc := func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
    log.Println("Signing transaction with address", address.Hex())
    keystore := keystore.NewKeyStore("/home/b-j-roberts/workspace/blockchain/my-chains/eth-private-network/data/keystore", keystore.StandardScryptN, keystore.StandardScryptP) //TODO: Hardcode
    //account, err := l1AccManager.Find(accounts.Account{Address: address})
    //if err != nil {
    //  return nil, fmt.Errorf("failed to find account: %v", err)
    //}
    return keystore.SignTxWithPassphrase(accounts.Account{Address: address}, "password", tx, big.NewInt(505)) //TODO: Hardcode
  }

  log.Println("Posting batch with id", id, "and hash", hash, "from address", addressString)
  tx, err := l1InOut.Contract.StoreBatch(&bind.TransactOpts{
    From: common.HexToAddress(addressString),
    Value: big.NewInt(0),
    GasLimit: 3000000, //TODO: Hardcode
    GasPrice: big.NewInt(200), //TODO: Hardcode
    Signer: SignerFunc,
  //  Signer: accounts.NewManage
  }, big.NewInt(id), hash, transactionByteData)
  if err != nil {
    return fmt.Errorf("failed to store batch: %v", err)
  }
  log.Println("Batch stored with transaction hash", tx.Hash().Hex())

  batchCount, err = l1InOut.Contract.GetBatchCount(nil)  
  if err != nil {  
    log.Fatalf("Failed to get batch count: %v", err)  
    return nil  
  }  
  log.Println("Batch number is", batchCount)

  return nil
}

type Batcher struct {
  L2Blockchain *core.BlockChain
  BlockIdx     uint64
  TxBatch      []*types.Transaction
  PostedBlockIdx uint64
  LastPostTime time.Time

  L1Comms *L1InOut
  BatchId int64
}

func NewBatcher(l2Blockchain *core.BlockChain, l1NodeUrl string, node *node.Node) *Batcher {
  return &Batcher{
    L2Blockchain: l2Blockchain,
    BlockIdx:     0,
    TxBatch:      make([]*types.Transaction, 0),
    PostedBlockIdx: 0,
    LastPostTime: time.Now(), //TODO
    L1Comms: NewL1InOut(l1NodeUrl, node),
    BatchId: 0,
  }
}

func (batcher *Batcher) PostBatch() error {
  if len(batcher.TxBatch) == 0 {
    return fmt.Errorf("no txs to post")
  }

  log.Printf("Posting batch of %d txs\n", len(batcher.TxBatch))
  transactionByteData := make([]byte, 0)
  //TODO: Compress transaction data
  for _, tx := range batcher.TxBatch {
    log.Printf("Posting tx: %s\n", tx.Hash().Hex())
    txBin, err := tx.MarshalBinary()
    if err != nil {
      return err
    }
    log.Printf("Tx: %s\n", txBin)
    transactionByteData = append(transactionByteData, txBin...)
  }

  //TODO: Post to l1 contract
  //TODO: Temporarily posting to file
  f, err := os.OpenFile("/home/b-j-roberts/naive-sequencer-data/batched_txs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //TODO: Hardcode
  if err != nil {
    return err
  }
  defer f.Close()

  f.WriteString(fmt.Sprintf("Idx Before Post: %d\n", batcher.PostedBlockIdx))
  f.WriteString(fmt.Sprintf("Posting Count: %d\n", len(batcher.TxBatch)))
  f.WriteString(fmt.Sprintf("Idx After Post: %d\n", batcher.BlockIdx))
  f.Write(transactionByteData)
  f.WriteString("\n")

  byteDataHash := sha256.Sum256(transactionByteData)// TODO: Use blockchain root

  log.Printf("Posting batch with values: %d, %d, %x\n", batcher.BatchId, byteDataHash, transactionByteData)
  err = batcher.L1Comms.PostBatch(transactionByteData, batcher.BatchId, byteDataHash)
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

      if len(batcher.TxBatch) > 10 || (time.Since(batcher.LastPostTime) > 1 * time.Minute && len(batcher.TxBatch) > 0) { //TODO: Hardcode
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
      time.Sleep(100 * time.Millisecond)
    }
  }}

  go runFunc()

  return nil
}

type Node struct {
  NaiveDb ethdb.Database
  Node    *node.Node
  L2Blockchain *core.BlockChain
  Engine  consensus.Engine
  Eth     *eth.Ethereum
  Batcher *Batcher
  // TODO: L1Reader, BatchPoster, ...
}

type NaiveAPI struct {
  version string
}

func NewNaiveAPI() *NaiveAPI {
  return &NaiveAPI{
    version: "1.0",
  }
}

func NewNode(naiveDb ethdb.Database, node *node.Node, chainDb ethdb.Database, l2Blockchain *core.BlockChain, engine consensus.Engine, config *ethconfig.Config) (*Node, error) {
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

  //TODO: APIs / RPC
  return &Node{
    NaiveDb: naiveDb,
    Node:    node,
    L2Blockchain: l2Blockchain,
    Engine:  engine,
    Eth:     naive_eth,
    Batcher:   NewBatcher(l2Blockchain, "http://localhost:8545", node), //TODO: Hardcode
  }, nil
}

func (node *Node) Start() error {
  //unauth, apis := node.Node.GetAPIs2()
  //log.Println("These are the APIs", unauth, apis)
  //log.Println("Some Configs", node.Node.IPCEndpoint(), node.Node.HTTPEndpoint())

  // Read address from json file under .address field
  //address := ""
  //jsonFile, err := os.Open(node.config.DataDir + j)
  //if err != nil {
  //  return fmt.Errorf("failed to open address.json: %v", err)
  //}
  //defer jsonFile.Close()
  //jsonParser := json.NewDecoder(jsonFile)
  //if err = jsonParser.Decode(&address); err != nil {
  //  return fmt.Errorf("failed to decode address.json: %v", err)
  //}
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

  return nil
}

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Starting sequencer...")

  //TODO: Learn more about default config
  nodeConfig := node.DefaultConfig
  nodeConfig.DataDir = "/home/b-j-roberts/naive-sequencer-data" //TODO: Hardcode
  nodeConfig.P2P.ListenAddr = ""
  nodeConfig.P2P.NoDial = true
  nodeConfig.P2P.NoDiscovery = true
  //TODO: setup command line args?
  nodeConfig.IPCPath = "naive-sequencer.ipc"// TODO: learn more about ipc
  nodeConfig.HTTPHost = "localhost"
  nodeConfig.HTTPPort = 5055 //TODO: Hardcode
  nodeConfig.HTTPModules = append(nodeConfig.HTTPModules, "personal") //TODO: Hardcode
  nodeConfig.HTTPModules = append(nodeConfig.HTTPModules, "naive")
  //nodeConfig.HTTPModules = append(nodeConfig.HTTPModules, "eth")

  // Mimicing eth/backend.go:New
  node, err := node.New(&nodeConfig)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Node created", node.DataDir())

  am := node.AccountManager()
  am.AddBackend(keystore.NewKeyStore(node.KeyStoreDir(), keystore.StandardScryptN, keystore.StandardScryptP))

  //TODO: chainDb more research on args
  // Handles, Persistent Chain Dir, & Ancient from nitro/cmd/conf/database.go
  // Caching from arbnode/execution/blockchain.go DatabaseCache
  // Namespace is prefix for metrics
  // Open rawdb from geth/core with ancients freezer & configs from arbitrum chainDb ( Disk based db )
  chainDb, err := node.OpenDatabaseWithFreezer("l2-chain", 2048, 512, "", "naive-l2/chaindb", false)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Chain DB created", chainDb)

  file, err := os.Open(nodeConfig.DataDir + "/genesis.json") //TODO: Hardcode
  if err != nil {
    utils.Fatalf("Failed to read genesis file: %v", err)
  }
  defer file.Close()
  genesis := new(core.Genesis)
  if err := json.NewDecoder(file).Decode(genesis); err != nil {
    utils.Fatalf("invalid genesis file: %v", err)
  }
  triedb := trie.NewDatabaseWithConfig(chainDb, &trie.Config{
    Preimages: true,  
  })
  _, hash, err := core.SetupGenesisBlock(chainDb, triedb, genesis)
  if err != nil {
    utils.Fatalf("Failed to setup genesis block: %v", err)
  }
  log.Println("Genesis block created", hash)

  backends := node.AccountManager().Backends(keystore.KeyStoreType)
  if len(backends) == 0 {
    log.Fatal("no key store backends found")
  }
  ks := backends[0].(*keystore.KeyStore)
  address := ks.Accounts()[0].Address

  //TODO: engine look into args
  //TODO: Genesis state & accounts / genesis.json setup
  config := ethconfig.Defaults
  config.Miner.Etherbase = address //TODO: Hardcode
  config.Miner.GasCeil = 300000000 //TODO: Hardcode
  ethashConfig := config.Ethash
  ethashConfig.NotifyFull = config.Miner.NotifyFull
  cliqueConfig, err := core.LoadCliqueConfig(chainDb, genesis)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Ethash config created", ethashConfig)
  log.Println("Clique config created", cliqueConfig)
  engine := ethconfig.CreateConsensusEngine(node, &ethashConfig, cliqueConfig, config.Miner.Notify, config.Miner.Noverify, chainDb)
  log.Println("Engine created", engine)

  //TODO: l2blockchain more research on args
  var l2BlockChain *core.BlockChain
  // Override the chain config with provided settings.
  var overrides core.ChainOverrides
  //if config.OverrideShanghai != nil {
  //  overrides.OverrideShanghai = config.OverrideShanghai
  //}
  vmConfig := vm.Config{
    EnablePreimageRecording: false,
  }
  txLookupLimi := uint64(31536000) // 1 year at 1 block per second
  l2BlockChain, err = core.NewBlockChain(chainDb, DefaultCacheConfigFor(node, false), genesis, &overrides, engine, vmConfig, shouldPreserveFalse, &txLookupLimi)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("L2 BlockChain created", l2BlockChain)

  //TODO: naivedb more research on args & what is this for
  naiveDb, err := node.OpenDatabase("naivedata", 0, 0, "", false)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Naive DB created", naiveDb)

  //valNode create & start
  ////TODO: close dbs & stop blockchain defers
  fatalErrChan := make(chan error, 10)

  naiveNode, err := NewNode(naiveDb, node, chainDb, l2BlockChain, engine, &config)
  if err != nil {
    fatalErrChan <- err
  }
  log.Println("Naive Node created", naiveNode)

  err = naiveNode.Start()
  if err != nil {
    fatalErrChan <- err
  }
  log.Println("Naive Node started", naiveNode)

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
