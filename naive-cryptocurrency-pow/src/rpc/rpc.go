package rpc

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"naive-cryptocurrency-pow/src/ledger"
	"naive-cryptocurrency-pow/src/mempool"
	"naive-cryptocurrency-pow/src/metrics"
)

func RpcSetup(rpcPort string) {
    // Create a new HTTP server to handle metrics requests
  // TODO: To jsonrpc
  http.Handle("/metrics", promhttp.Handler())
  http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
    metrics.RequestCount.Inc()
    w.Write([]byte("The exposed endpoints are : metrics, help, add, block\n"))
  })
  http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
    metrics.RequestCount.Inc()
    defer metrics.RequestDuration.Observe(time.Since(time.Now()).Seconds())

    value := r.URL.Query().Get("value")
    i, err := strconv.Atoi(value)
    if err != nil {
      log.Fatal(err)
    }
    if i >= len(ledger.TheLedger.Blockchain) {

      w.Write([]byte(fmt.Sprintf("Blockchain is of length : %d so cannot index position : %d\n", len(ledger.TheLedger.Blockchain), i)))
    } else {
      MaxUint := ^uint32(0)
      message := fmt.Sprintf("Block %d\n", i)
      message += fmt.Sprintf("    Block Header Hash : 0x%08x\n", ledger.CalculateBlockHeaderHash(ledger.TheLedger.Blockchain[i].BlockHeader))
      message += fmt.Sprintf("    Prev Block Hash : 0x%08x\n", ledger.TheLedger.Blockchain[i].BlockHeader.PrevBlockHash)
      message += fmt.Sprintf("    Timestamp : %s\n", time.Unix(ledger.TheLedger.Blockchain[i].BlockHeader.Timestamp, 0))
      message += fmt.Sprintf("    Difficulty : 0x%08x\n", MaxUint - ledger.TheLedger.Blockchain[i].BlockHeader.Difficulty)
      message += fmt.Sprintf("    Nonce %d\n", ledger.TheLedger.Blockchain[i].BlockHeader.Nonce)
      for transaction := range ledger.TheLedger.Blockchain[i].Transactions {
        message += fmt.Sprintf("    Transaction %d -> %d : %d", ledger.TheLedger.Blockchain[i].Transactions[transaction].FromAddress, ledger.TheLedger.Blockchain[i].Transactions[transaction].ToAddress, ledger.TheLedger.Blockchain[i].Transactions[transaction].Amount)
      }
      w.Write([]byte(message))
    }
  })
  http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
    // Track the number of requests
    metrics.RequestCount.Inc()
    // Track the duration of the request
    defer metrics.RequestDuration.Observe(time.Since(time.Now()).Seconds())
    value := r.URL.Query().Get("value")
    values := strings.Split(value, ",")

    // Assign the values to individual variables
    var from, to, amount, fee uint64
    var err error
    from, err = strconv.ParseUint(values[0], 16, 64)
    if err != nil {
      log.Fatal(err)
    }
    to, err = strconv.ParseUint(values[1], 10, 64)
    if err != nil {
      log.Fatal(err)
    }
    amount, err = strconv.ParseUint(values[2], 10, 64)
    if err != nil {
      log.Fatal(err)
    }
    fee, err = strconv.ParseUint(values[3], 10, 64)
    if err != nil {
      log.Fatal(err)
    }
    nonce, err := strconv.ParseUint(values[4], 10, 64)
    if err != nil {
      log.Fatal(err)
    }
    publicKeyFile := values[5]
    privateKeyFile := values[6]
  
    transaction, err := ledger.CreateTransaction(from, to, amount, fee, uint(nonce), privateKeyFile, publicKeyFile, true)
    if err != nil {
      log.Printf("Error creating transaction : %s", err)
    } else {
      mempool.TheMempool.AddTransaction(transaction, true)
      log.Printf("Added transaction to mempool")
      mempool.TheMempool.PrintMempool()
    }
    //transactions := []ledger.Transaction{transaction}
    //newBlock := ledger.CreateUnminedBlock(ledger.TheLedger.Blockchain[len(ledger.TheLedger.Blockchain)-1], transactions, uint32(len(ledger.TheLedger.Blockchain)))
    //newBlock.BlockHeader.Nonce = miner.MineBlockNonce(newBlock.BlockHeader)

    // if ledger.IsBlockValid(newBlock, ledger.TheLedger.Blockchain[len(ledger.TheLedger.Blockchain)-1], uint32(len(ledger.TheLedger.Blockchain))) {
    //  ledger.TheLedger.Mutex.Lock()
    //  ledger.TheLedger.Blockchain = append(ledger.TheLedger.Blockchain, newBlock)
    //  metrics.Blockheight.Set(float64(len(ledger.TheLedger.Blockchain)))

    //  MinerAddress := uint64(0)
    //  ledger.TheLedger.AccountNonces[transaction.FromAddress] = transaction.Nonce
    //  ledger.TheLedger.AccountBalances[transaction.FromAddress] -= transaction.Fee
    //  ledger.TheLedger.AccountBalances[MinerAddress] += transaction.Fee
    //  ledger.TheLedger.AccountBalances[transaction.FromAddress] -= transaction.Amount
    //  ledger.TheLedger.AccountBalances[transaction.ToAddress] += transaction.Amount
    //  ledger.TheLedger.Mutex.Unlock()

    //  ledger.StoreChainToFile(ledger.TheLedger) //TODO: remove mutex from ledger & make async
    //}
  })

  log.Printf("Serving RPC at %s", rpcPort)
  http.ListenAndServe(":" + rpcPort, nil)
}
