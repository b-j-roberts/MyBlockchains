package ledger

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"naive-cryptocurrency-pow/src/metrics"
)

type Ledger struct {
  AccountBalances map[uint64]uint64
  AccountNonces map[uint64]uint
  Blockchain []Block

  Mutex sync.Mutex //TODO: Move mutex out?
}

var TheLedger Ledger
var lock sync.Mutex

// Marshal is a function that marshals the object into an
// io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
  b, err := json.MarshalIndent(v, "", "\t")
  if err != nil {
    return nil, err
  }
  return bytes.NewReader(b), nil
}

func (ledger *Ledger) Create(snapshotPath string, airdropFile string) {
  if snapshotPath != "" {

    ledger.Mutex.Lock()
    //TODO: error check and genesis if not
    //TODO: Faster if from pointer?
    var err error
    (*ledger), err = loadChainFromFile(snapshotPath)
    if err != nil {
      log.Fatalf("Error loading chain from file", err)
    }
    ledger.Mutex.Unlock()
  } else {
    //TODO: Sign from a specific address and hardcode signature
    transaction := createGenesisTransaction()
    genesisBlock := createGenesisBlock(transaction)

    ledger.Mutex.Lock()
    //TODO: Ledger transfer value and refactor
    ledger.Blockchain = append(ledger.Blockchain, genesisBlock)
    metrics.Blockheight.Set(float64(len(ledger.Blockchain))) //TODO: Move to metrics

    // Genesis Account Map
    accounts:= Airdrop(airdropFile)
    ledger.AccountBalances = accounts
    ledger.AccountBalances[transaction.FromAddress] -= transaction.Amount
    ledger.AccountBalances[transaction.ToAddress] += transaction.Amount
    ledger.AccountNonces = make(map[uint64]uint)

    ledger.Mutex.Unlock()
  }
}

func PrintChain(blockchain []Block) {
  log.Println("Blockchain :")
  for i := 0; i < len(blockchain); i++ {
    log.Printf("Block %d", i)
    printBlock(blockchain[i])
  }
}

func PrintAccounts(count uint64) {
  log.Println("Accounts :")
  for i := uint64(0); i < count; i++ {
    log.Printf("Account %d : Balance %d", i, TheLedger.AccountBalances[i])
  }
}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
  lock.Lock()
  defer lock.Unlock()
  f, err := os.Create(path)
  if err != nil {
    return err
  }
  defer f.Close()
  r, err := Marshal(v)
  if err != nil {
    return err
  }
  _, err = io.Copy(f, r)
  return err
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
  return json.NewDecoder(r).Decode(v)
}

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func Load(path string) (Ledger, error) {
  log.Print("Loading Ledger...")
  metrics.SnapshotLoading.Set(1.0)
  lock.Lock()
  defer lock.Unlock()
  f, err := os.Open(path)
  if err != nil {
    return TheLedger, err
  }
  defer f.Close()
  var checkLedger Ledger
  if err := Unmarshal(f, &checkLedger); err != nil {
    return TheLedger, err
  }

  log.Print("Verifying Ledger...")
  metrics.SnapshotLoading.Set(2.0)
  if VerifyBlockchain(checkLedger.Blockchain) {
    //TODO: FOr all spots and refactor
    log.Print("Verified!")
    TheLedger = checkLedger
    metrics.Blockheight.Set(float64(len(TheLedger.Blockchain)))
  }

  //TODO: Load initial state cause not valid
  log.Print("Loaded!")

  metrics.SnapshotLoading.Set(3.0)
  return TheLedger, nil
}

func StoreChainToFile(ledger Ledger) {
  log.Println("Saving snapshot to file")
  if err := Save("./snapshot.tmp", ledger); err != nil {
    log.Fatalln(err)
  }
}

func loadChainFromFile(snapPath string) (Ledger, error) {
  log.Println("Loading snapshot from file")
  ledger, err := Load(snapPath)
  if err != nil {
    log.Fatalln(err)
  }

  return ledger, nil
}

func GetBlock(height uint32) Block {
  return TheLedger.Blockchain[height]
}
