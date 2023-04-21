package rpc

import (
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

func RpcTransactionSetup(rpcPort string) {
    // Create a new HTTP server to handle metrics requests
  // TODO: To jsonrpc
  http.Handle("/metrics", promhttp.Handler())
  http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
    metrics.RequestCount.Inc()
    w.Write([]byte("The exposed endpoints are : metrics, help, add, block\n"))
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
    //TODO: Do hex to uint64
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
  
    transaction, err := ledger.CreateTransaction(from, to, amount, fee, uint(nonce), privateKeyFile, publicKeyFile, false)
    if err != nil {
      log.Printf("Error creating transaction: %s", err)
    } else {
      mempool.TheMempool.AddTransaction(transaction, true)
      // Print confirmation to the console using log
      log.Printf("Added transaction")
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
