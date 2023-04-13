// References :
// https://github.com/nosequeldeebee/blockchain-tutorial/tree/master/p2p
// https://github.com/nosequeldeebee/blockchain-tutorial/blob/master/proof-work/main.go
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html

package main

import (
	"flag"
	"log"
	"time"

	golog "github.com/ipfs/go-log"

	"naive-cryptocurrency-pow/src/ledger"
	"naive-cryptocurrency-pow/src/mempool"
	"naive-cryptocurrency-pow/src/metrics"
	"naive-cryptocurrency-pow/src/miner"
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

  ledger.TheLedger.Create(*snapPath, *airdropFile)
  ledger.PrintChain(ledger.TheLedger.Blockchain)
  ledger.PrintAccounts(5)

  p2p.DialPeer(*listenLedgerPort, *secio, *seed, *peerToCallLedger)
  p2p.DialMempoolPeer(*listenMempoolPort, *secio, *seed, *mempoolPeerFilename)
  metrics.PromSetup()

  log.Println("RPC port: ", *rpcPort)

  
  //TODO: Waitgroup
  go rpc.RpcSetup(*rpcPort)

  for {
    log.Println("Mining...")
    if len(mempool.TheMempool.AvailableTransactions) == 0 {
      log.Println("No transactions to mine")
      time.Sleep(1 * time.Second)
      continue
    }

    //TODO: Add a way for mining to stop if new block is received
    transactions := []ledger.Transaction{mempool.TheMempool.AvailableTransactions[0]}
    prevBlock := ledger.TheLedger.Blockchain[len(ledger.TheLedger.Blockchain) - 1]
    createdBlock, err := ledger.CreateUnminedBlock(prevBlock, transactions, uint32(len(ledger.TheLedger.Blockchain)))
    if err != nil {
      log.Println("Error creating block: ", err)
      //TODO: Only remove the transaction that failed
      mempool.TheMempool.AvailableTransactions = mempool.TheMempool.AvailableTransactions[1:]
      continue
    }

    log.Println("Created block: ", createdBlock)

    createdBlock.BlockHeader.Nonce = miner.MineBlockNonce(createdBlock.BlockHeader)
    log.Println("Mined block: ", createdBlock)

    if !ledger.IsBlockValid(createdBlock, prevBlock, uint32(len(ledger.TheLedger.Blockchain))) {
      log.Fatal("Block is not valid")
      mempool.TheMempool.AvailableTransactions = mempool.TheMempool.AvailableTransactions[1:]
      //TODO: Move this logic to transaction check not block check
    }
  
    ledger.TheLedger.Blockchain = append(ledger.TheLedger.Blockchain, createdBlock)
    mempool.TheMempool.AvailableTransactions = mempool.TheMempool.AvailableTransactions[1:]

    ledger.PrintChain(ledger.TheLedger.Blockchain)

    time.Sleep(1 * time.Second)
  }
}
