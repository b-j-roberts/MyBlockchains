package main

import (
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
)

// Don't preserve reorg'd out blocks
func ShouldPreserveFalse(_ *types.Header) bool {
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

func NodeConfig(dataDir string, httpHost string, httpPort int, httpModules string) *node.Config {
  //TODO: Learn more about node config + default config
  nodeConfig := node.DefaultConfig                          
  nodeConfig.DataDir = dataDir                             
  //nodeConfig.P2P.ListenAddr = ""                            
  //nodeConfig.P2P.NoDial = true                              
  //nodeConfig.P2P.NoDiscovery = true
  nodeConfig.P2P.ListenAddr = ":30313"
  nodeConfig.IPCPath = "naive-sequencer.ipc"// TODO: learn more about ipc
  nodeConfig.HTTPHost = httpHost
  nodeConfig.HTTPPort = httpPort
  nodeConfig.HTTPCors = []string{"*"}
  log.Println("httpModules: ", httpModules)
  nodeConfig.HTTPModules = append(nodeConfig.HTTPModules, strings.Split(httpModules, ",")...)

  log.Println("Node config: ", nodeConfig)

  return &nodeConfig
}

func EthConfig(address common.Address) *ethconfig.Config {
  //TODO: Learn more about eth config + default config
  config := ethconfig.Defaults
  config.Miner.Etherbase = address
  config.Miner.GasCeil = 300000000 // 10x default
  config.Ethash.NotifyFull = config.Miner.NotifyFull

  return &config
}
