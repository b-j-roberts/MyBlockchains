package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/trie"

	l2core "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/src/core"
)

func CreateNaiveNode(dataDir string, httpHost string, httpPort int, httpModules string) (*l2core.Node, error) {
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

  naiveNode, err := l2core.NewNode(node, chainDb, l2BlockChain, engine, ethConfig, nil) //TODO: nil when or when bridge not needed?
  if err != nil {
    return nil, fmt.Errorf("failed to create naive node: %v", err)
  }


  return naiveNode, nil
}

func main() { os.Exit(mainImpl()) }

func mainImpl() int {
  log.Println("Starting rpc...")

  osHomeDir, err := os.UserHomeDir()
  dataDir := flag.String("datadir", osHomeDir + "/naive-rpc-data", "data directory for the database and keystore")
  //TODO: to url not host port
  httpHost := flag.String("httphost", "localhost", "HTTP-RPC server listening interface")
  httpPort := flag.Int("httpport", 5056, "HTTP-RPC server listening port")
  httpModules := flag.String("httpmodules", "personal,naive", "Comma separated list of API modules to enable on the HTTP-RPC interface")
  flag.Parse()

  naiveNode, err := CreateNaiveNode(*dataDir, *httpHost, *httpPort, *httpModules)
  if err != nil {
    utils.Fatalf("Failed to create naive rpc node: %v", err)
  }

  ////TODO: close dbs & stop blockchain defers
  fatalErrChan := make(chan error, 10)

  err = naiveNode.Start()
  if err != nil {
    fatalErrChan <- err
  }
  log.Println("Naive Rpc Node started", naiveNode)

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
