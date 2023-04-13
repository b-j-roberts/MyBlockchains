// References :
// https://github.com/nosequeldeebee/blockchain-tutorial/tree/master/p2p
// https://github.com/nosequeldeebee/blockchain-tutorial/blob/master/proof-work/main.go
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html

package main

import (
	"flag"

	golog "github.com/ipfs/go-log"

	"naive-cryptocurrency-pow/src/ledger"
	"naive-cryptocurrency-pow/src/metrics"
	"naive-cryptocurrency-pow/src/p2p"
	"naive-cryptocurrency-pow/src/rpc"
	"naive-cryptocurrency-pow/src/signer"
)

func main() {
	golog.SetAllLoggers(golog.LevelInfo) // Change to DEBUG for extra info

  //TODO: Requires
  // Parse options from the command line
  listenMempoolPort := flag.Int("mem-port", 8985, "Port to listen for Mempool connections")
  listenLedgerPort := flag.Int("blk-port", 8986, "Port to listen for Ledger connections")
  mempoolPeerFilename := flag.String("mem-peer", "config/peer-list-mempool.txt", "Peer path filename to dial for mempool connection")
  peerToCallLedger := flag.String("blk-peer", "", "Peer port / path to dial")
  rpcPort := flag.String("rpc", "8987", "RPC port / prom metrics port")
  secio := flag.Bool("secio", false, "enable secio")
  seed := flag.Int64("seed", 0, "Seed for id generation")
  snapPath := flag.String("snap", "", "Load blockchain from snapshot")
  airdropFile := flag.String("air", "config/airdrop_config.csv", "Airdrop addresses")
  accountsDirectory := flag.String("accounts", "accounts", "Directory containing account dirs")
  accountNumber := flag.Int("account-id", 12345, "Account number to use from account directory.")
  flag.Parse()

  if !signer.VerifyKeyPairExists(*accountsDirectory, *accountNumber) {
    signer.GeneratePublicPrivateKey(*accountsDirectory, *accountNumber)
  }

//  if *listenMempoolPort == 0 {
//    log.Fatal("You must specify a mempool port to bind the client on. Use the 'port' flag")
//  }
//  if *listenLedgerPort == 0 {
//    log.Fatal("You must specify a ledger port to bind the client on. Use the 'port' flag")
//  }

  ledger.TheLedger.Create(*snapPath, *airdropFile)
  ledger.PrintChain(ledger.TheLedger.Blockchain)
  ledger.PrintAccounts(5)

  p2p.DialPeer(*listenLedgerPort, *secio, *seed, *peerToCallLedger)
  p2p.DialMempoolPeer(*listenMempoolPort, *secio, *seed, *mempoolPeerFilename)
  metrics.PromSetup()
  rpc.RpcSetup(*rpcPort)

  select {} // hang forever
}
