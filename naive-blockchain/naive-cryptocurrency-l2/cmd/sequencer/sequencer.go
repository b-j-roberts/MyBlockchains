package main

import (
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
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/metrics/exp"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/trie"

	l2core "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/core"
	l2utils "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/utils"
)

func SetupMetrics() {
  metrics.Enabled = true
  metrics.EnabledExpensive = true
  exp.Exp(metrics.DefaultRegistry)
  address := fmt.Sprintf("%s:%d", "localhost", 6160)
  log.Println("Metrics address is", address)
  exp.Setup(address)
}

func StartMetrics() error {
  go metrics.CollectProcessMetrics(3 * time.Second)

  return nil
}

func CreateNaiveNode(genesisFile string, dataDir string, httpHost string, httpPort int, httpModules string, l1Host string, l1Port int, l1ChainID int, miningThreads int,
                     l1ContractAddress common.Address, l1BridgeContractAddress string, posterAddress common.Address, metricsEnabled bool) (*l2core.Sequencer, error) {
  // Function used to create Naive Node mimicing eth/backend.go:New for Ethereum Node object
  if metricsEnabled {
    SetupMetrics()
  }

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
  chainDb, err := node.OpenDatabaseWithFreezer("l2-chain", 2048, 512, "", "naive_l2/chaindb", false)
  if err != nil {
    return nil, fmt.Errorf("failed to open chain database: %v", err)
  }

  // Setup Genesis
  file, err := os.Open(genesisFile)
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
  cliqueConfig.L1Url = fmt.Sprintf("http://%s:%d", l1Host, l1Port)
  cliqueConfig.L1BridgeAddress = l1BridgeContractAddress
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

  naiveNode, err := l2core.NewSequencer(node, chainDb, l2BlockChain, engine, ethConfig, l1ContractAddress, common.HexToAddress(l1BridgeContractAddress), posterAddress, l1Host, l1Port, l1ChainID, miningThreads)
  if err != nil {
    return nil, fmt.Errorf("failed to create naive node: %v", err)
  }


  return naiveNode, nil
}

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Starting sequencer...")
osHomeDir, err := os.UserHomeDir()
  genesisFile := flag.String("genesis", osHomeDir + "/naive-sequencer-data/genesis.json", "genesis file for the L2 blockchain")
  dataDir := flag.String("datadir", osHomeDir + "/naive-sequencer-data", "data directory for the database and keystore")
  httpHost := flag.String("httphost", "localhost", "HTTP-RPC server listening interface")
  httpPort := flag.Int("httpport", 5055, "HTTP-RPC server listening port")
  httpModules := flag.String("httpmodules", "personal,naive", "Comma separated list of API modules to enable on the HTTP-RPC interface")
  l1ContractAddress := flag.String("l1contract", "", "Address of the L1 contract")
  l1BridgeContractAddress := flag.String("l1bridgecontract", "", "Address of the L1 bridge contract")
  sequencerAddress := flag.String("sequencer", "", "Address of the sequencer on L1")
  sequencerKeystore := flag.String("sequencerkeystore", "", "Keystore file for the sequencer on L1")
  l1Host := flag.String("l1host", "localhost", "L1 HTTP-RPC server listening interface")
  l1Port := flag.Int("l1port", 8545, "L1 HTTP-RPC server listening port")
  l1ChainID := flag.Int("l1chainid", 505, "L1 chain ID")
  miningThreads := flag.Int("miningthreads", 4, "Number of threads to use for mining")
  metricsFlag := flag.Bool("metrics", false, "Enable metrics")
  flag.Parse()

  log.Println("Connecting to L1 contract at", *l1Host, *l1Port, "with address", *l1ContractAddress)
  naiveNode, err := CreateNaiveNode(*genesisFile, *dataDir, *httpHost, *httpPort, *httpModules, *l1Host, *l1Port, *l1ChainID, *miningThreads, common.HexToAddress(*l1ContractAddress), *l1BridgeContractAddress, common.HexToAddress(*sequencerAddress), *metricsFlag)
  if err != nil {
    utils.Fatalf("Failed to create naive sequencer node: %v", err)
  }

  l2utils.RegisterAccount(common.HexToAddress(*sequencerAddress), *sequencerKeystore)

  ////TODO: close dbs & stop blockchain defers
  fatalErrChan := make(chan error, 10)

  genesis, err := core.ReadGenesis(naiveNode.L2Node.ChainDb)
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

  if *metricsFlag {
    err = StartMetrics()
    if err != nil {                                                                          
       fatalErrChan <- fmt.Errorf("failed to start metrics: %v", err)                                  
    }
  }

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
