package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/metrics/exp"
	"github.com/ethereum/go-ethereum/node"
)

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

type NodeBaseConfig struct {
  DataDir      string `json:"data-dir"`
  Genesis      string `json:"genesis"`
  Contracts    string `json:"contracts"`
  L2ChainID    int    `json:"l2ChainId"`
  Host         string `json:"host"`
  Port         int    `json:"port"`
  P2PPort      int    `json:"p2pport"`
  Modules      string `json:"modules"`
  L1URL        string `json:"l1Url"`
  L1ChainID    int    `json:"l1ChainId"`
  MiningThreads int   `json:"miningThreads"`
  Metrics      struct {
    Enabled bool   `json:"enabled"`
    Host    string `json:"host"`
    Port    int    `json:"port"`
  } `json:"metrics"`
}

//TODO: Move somewhere else
func SetupMetrics(nodeBaseConfig *NodeBaseConfig) {
  log.Println("Metrics enabled: ", nodeBaseConfig.Metrics.Enabled)
  if !nodeBaseConfig.Metrics.Enabled {
    return
  }

  metrics.Enabled = true
  metrics.EnabledExpensive = true
  exp.Exp(metrics.DefaultRegistry)
  address := fmt.Sprintf("%s:%d", nodeBaseConfig.Metrics.Host, nodeBaseConfig.Metrics.Port)
  log.Println("Metrics address is", address)
  exp.Setup(address)
}

func LoadNodeBaseConfig(configFile string) (*NodeBaseConfig, error) {
  file, err := os.Open(configFile)
  if err != nil {
    return nil, fmt.Errorf("failed to open config file: %v", err)
  }
  defer file.Close()

  config := new(NodeBaseConfig)
  if err := json.NewDecoder(file).Decode(config); err != nil {
    return nil, fmt.Errorf("invalid config file: %v", err)
  }

  return config, nil
}

func NodeConfig(nodeBaseConfig *NodeBaseConfig) *node.Config {
  //TODO: Learn more about node config + default config
  nodeConfig := node.DefaultConfig
  nodeConfig.DataDir = nodeBaseConfig.DataDir
  //nodeConfig.P2P.ListenAddr = ""                            
  //nodeConfig.P2P.NoDial = true                              
  //nodeConfig.P2P.NoDiscovery = true
  nodeConfig.P2P.ListenAddr = ":" + fmt.Sprintf("%d", nodeBaseConfig.P2PPort)
  nodeConfig.IPCPath = "naive-sequencer.ipc"
  nodeConfig.HTTPHost = nodeBaseConfig.Host
  nodeConfig.HTTPPort = nodeBaseConfig.Port
  nodeConfig.HTTPCors = []string{"*"}
  log.Println("httpModules: ", nodeBaseConfig.Modules)
  nodeConfig.HTTPModules = append(nodeConfig.HTTPModules, strings.Split(nodeBaseConfig.Modules, ",")...)

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
