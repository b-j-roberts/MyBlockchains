package main

import (
	"flag"
	"naive-cryptocurrency-pow/src/p2p"
	"naive-cryptocurrency-pow/src/rpc"

	golog "github.com/ipfs/go-log"
)

func main() {
  golog.SetAllLoggers(golog.LevelInfo) // Change to DEBUG for extra info

  // Parse options from the command line
  listenPort := flag.Int("port", 8985, "Port to listen for p2p mempool connections")
  //TODO: Same peer as ledger peer?
  peerFilename := flag.String("peer", "config/peer-list-mempool.txt", "Peer path filename to dial for mempool connection")
  seed := flag.Int64("seed", 0, "Seed for id generation")
  rpcPort := flag.String("rpc", "8987", "RPC port / prom metrics port")
  secio := flag.Bool("secio", false, "enable secio")
  flag.Parse()

  //if *listenPort == 0 {
  //  log.Fatal("You must specify a port to bind the client on. Use the 'port' flag")
  //}

  p2p.DialMempoolPeer(*listenPort, *secio, *seed, *peerFilename)
  //TODO: Remove read data and only write data on new transactions till found by x peers
  rpc.RpcTransactionSetup(*rpcPort)

  select {} // hang forever
}

