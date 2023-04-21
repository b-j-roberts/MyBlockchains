package mempool

import (
	"fmt"
	"sync"

	"naive-cryptocurrency-pow/src/ledger"
)

type Mempool struct {
  AvailableTransactions []ledger.Transaction
  Mutex                 sync.Mutex
}

var TheMempool Mempool

func MempoolContainsTransaction(tx ledger.Transaction) bool {
  for _, t := range TheMempool.AvailableTransactions {
    if ledger.CompareTransactions(t, tx) {
      return true
    }
  }
  return false
}

func (mempool *Mempool) AddTransaction(transaction ledger.Transaction, check bool) {
  if check == false || !MempoolContainsTransaction(transaction) {
    mempool.AvailableTransactions = append(mempool.AvailableTransactions, transaction)
  }
}

func (mempool *Mempool) String() string {
  var ret string
  for i, t := range mempool.AvailableTransactions {
    ret += fmt.Sprintf("Transaction %d:\n", i+1)
    ret += fmt.Sprintf("  From Address: %d\n", t.FromAddress)
    ret += fmt.Sprintf("  To Address: %d\n", t.ToAddress)
    ret += fmt.Sprintf("  Amount: %d\n", t.Amount)
    ret += fmt.Sprintf("  Fee: %d\n", t.Fee)
    ret += fmt.Sprintf("  Signature: %x\n", t.Signature)
  }

  return ret
}

func (mempool *Mempool) PrintMempool() {

  fmt.Println("Mempool:")
  for i, t := range mempool.AvailableTransactions {
    fmt.Printf("Transaction %d:\n", i+1)
    ledger.PrintTransaction(t)
  }
}
