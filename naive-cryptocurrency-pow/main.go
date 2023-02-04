// References :
// https://github.com/nosequeldeebee/blockchain-tutorial/tree/master/p2p
// https://github.com/nosequeldeebee/blockchain-tutorial/blob/master/proof-work/main.go
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html

package main

import (
  "bufio"
  "bytes"
  "context"
  crypt "crypto"
	"crypto/rand"
  "crypto/rsa"
	"crypto/sha256"
  "crypto/x509"
  "encoding/binary"
  "encoding/gob" //TODO: Replace?
  "encoding/pem"
  "encoding/hex"
  "encoding/json"
  "io"
  "io/ioutil"
  "flag"
  "fmt"
  "net/http"
  "log"
  "os"
  mrand "math/rand"
  "math/big"
  "strings"
  "strconv"
  "sync"
  "time"
  //valid "github.com/asaskevich/govalidator"

  unsafe "unsafe"

  //"github.com/davecgh/go-spew/spew"
	golog "github.com/ipfs/go-log"

  libp2p "github.com/libp2p/go-libp2p"
  crypto "github.com/libp2p/go-libp2p/core/crypto"
  host "github.com/libp2p/go-libp2p/core/host"
  net "github.com/libp2p/go-libp2p/core/network"
  peer "github.com/libp2p/go-libp2p/core/peer"
  pstore "github.com/libp2p/go-libp2p/core/peerstore"


  ma "github.com/multiformats/go-multiaddr"

  "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/btcsuite/btcutil/base58"
)

// TODO: For now a block contains one datapoint, in the future change to multiple transactions
// TODO: match closer https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html 
type BlockHeader struct {
  Version int // TODO: Why version
  PrevBlockHash uint32
  TransactionHash uint32
  Timestamp int64
  Difficulty uint32 // TODO: How decided & not changed to dupe system to be faster
  Nonce int
}

type Transaction struct {
  FromAddress uint64
  ToAddress uint64
  Amount uint64
  Fee uint64
  Signature []byte
}

type Block struct {
  BlockSize uint
  BlockHeader BlockHeader
  Transaction Transaction
}

type Ledger struct {
  AccountMap map[uint64]uint64
  Blockchain []Block
}

var Mempool []Transaction

var mutex = &sync.Mutex{}
var lock sync.Mutex

var TheLedger Ledger

// Block Header hash -- Obtained by using sha256 twice on blockheader ( double sha256 )
//    Not stored on chain, but computed/checked each time a node receives the new block

// TODO: Store block height & hashes in indexed db somewhere?
//       How does main btc client do this? can you make calls based on blockheight?
//       NOTE : Blockheight not necessarily unique block


// When receiving block : Validate - check prev hash, ensure block header hash correct

//TODO: Mempool, transaction keys, multi transaction, blocktime & dynamic difficulty, cleanup / refactor, account balance check on snapshot, all verifications, multipeers, miner fees, airdrop keys, CI/CD and full deploy pipeline...

// Account id / node id  / private key
// Sign transactions
// Mempool ( with Ledger for now )

// Kube

func calculateTransactionHash(transaction Transaction) uint32 {
  //var network bytes.Buffer
  //enc := gob.NewEncoder(&network)
  //err := enc.Encode(transaction)
  //if err != nil {
  //  panic(err)
  //}
  //data_arr :=  network.Bytes()
  data, err := json.Marshal(transaction)
  if err != nil {
    panic(err)
  }

  hasher := sha256.New()
  hasher.Write(data)
  hashed := hasher.Sum(nil)

  log.Printf("Calculated Transaction Hash : 0x%08x for Transaction 0x%08x -> 0x%08x : %d", uint32(binary.BigEndian.Uint32(hashed)), transaction.FromAddress, transaction.ToAddress, transaction.Amount)
  return uint32(binary.BigEndian.Uint32(hashed))
}

func calculateBlockHeaderHash(blockHeader BlockHeader) uint32 {
  //var network bytes.Buffer
  //enc := gob.NewEncoder(&network)
  //err := enc.Encode(blockHeader)
  //if err != nil {
  //  panic(err)
  //}
  //data_arr :=  network.Bytes()
  data, err := json.Marshal(blockHeader)
  if err != nil {
    panic(err)
  }
  
  hasher := sha256.New()
  hasher.Write(data)
  hashed := hasher.Sum(nil)

  log.Printf("Calculated Header Hash : 0x%08x", uint32(binary.BigEndian.Uint32(hashed)))
  return uint32(binary.BigEndian.Uint32(hashed))
}

// PublicKeyToAddress converts the public key to an address using Base58Check encoding
func PublicKeyToAddress(pk *rsa.PublicKey) (uint64, error) {
	// Hash the public key using SHA256
  hexString := hex.EncodeToString(pk.N.Bytes())
	hash := sha256.Sum256([]byte(hexString))

	// Add a checksum to the end of the hash
  checksum := hash[:4] //TODO: Is this a checksum?
  encoded := append(hash[:], checksum...)

	// Encode the result using Base58Check
	a := base58.Encode(encoded)

  // Encode the address as a hexadecimal string
	hexString = hex.EncodeToString([]byte(a))

	// Parse the hexadecimal string to a uint64
	address, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, err
	}

	return address, nil
}

// AddressToPublicKey converts an address to the RSA public key
func AddressToPublicKey(address uint64) (*rsa.PublicKey, error) {
	// Decode the address using Base58Check encoding
	decoded := base58.Decode(fmt.Sprintf("%x", address))

	// Extract the hash part of the decoded address
	hash := decoded[:len(decoded)-4]

	// Convert the hash to a hexadecimal string
	hexString := hex.EncodeToString(hash)

	// Convert the hexadecimal string back to a byte array
	b := []byte(hexString)

	// Create a new RSA public key from the byte array
	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(b),
		E: 65537,
	}

	return publicKey, nil
}

func signTransaction(tx *Transaction) {
  accountNumber := "12345"
  privateKeyFile, err := os.Open("private_" + accountNumber + ".pem")
	if err != nil {
		fmt.Println("Error opening private key file:", err)
		return
	}
	defer privateKeyFile.Close()

	// privateKeyDecoder := gob.NewDecoder(privateKeyFile)
	// privateKey := &rsa.PrivateKey{}
  // err = privateKeyDecoder.Decode(privateKey)
  // if err != nil {
  //   fmt.Println("Error decoding private key file:", err)
  //   return
  // }

  privateKeyPEM, err := ioutil.ReadAll(privateKeyFile)
  if err != nil {
    fmt.Println("Error reading private key file:", err)
    return
  }
  
  privateKeyBlock, _ := pem.Decode(privateKeyPEM)
  if privateKeyBlock == nil {
    fmt.Println("Error decoding private key PEM")
    return
  }
  
  privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
  if err != nil {
    fmt.Println("Error parsing private key:", err)
    return
  }

  txHash := sha256.Sum256([]byte(fmt.Sprintf("%d%d%d%d", tx.FromAddress, tx.ToAddress, tx.Amount, tx.Fee)))
  // Sign the hash with the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypt.SHA256, txHash[:])
	if err != nil {
		fmt.Println("Error signing transaction:", err)
		return
	}

	// Save the signature to the transaction
	tx.Signature = signature

	fmt.Println("Transaction:", tx)
	fmt.Println("Signature:", hex.EncodeToString(tx.Signature))
}

func verifyTransactionFromSender(tx Transaction) {
  // Load the public key -- TODO: convert from FromAddress
	accountNumber := string(12345)
	publicKeyFile, err := os.Open("public_" + accountNumber + ".pem")
	if err != nil {
		fmt.Println("Error opening public key file:", err)
		return
	}
	defer publicKeyFile.Close()

	publicKeyDecoder := gob.NewDecoder(publicKeyFile)
	publicKey := &rsa.PublicKey{}
	publicKeyDecoder.Decode(publicKey)

	// Hash the transaction data
	txHash := sha256.Sum256([]byte(fmt.Sprintf("%d%d%d%d", tx.FromAddress, tx.ToAddress, tx.Amount, tx.Fee)))

	// Verify the signature
	err = rsa.VerifyPKCS1v15(publicKey, crypt.SHA256, txHash[:], tx.Signature)
	if err != nil {
		fmt.Println("Error verifying signature:", err)
		return
	}

	fmt.Println("Signature verified.")
}

func mineBlockNonce(blockHeader BlockHeader) int {
  MaxUint := ^uint32(0)
  threshhold := MaxUint - blockHeader.Difficulty
  for {
    blockHash := calculateBlockHeaderHash(blockHeader)
    if(blockHash < threshhold) {
      log.Printf("Block mined with hash %08x for block header :", blockHash)
      printBlockHeader(blockHeader)
      break
    } else {
      log.Printf("Nonce %d failed with hash 0x%08x", blockHeader.Nonce, blockHash)
    }
    blockHeader.Nonce += 1
  }
  blocksMined.Inc()
  return blockHeader.Nonce
}

// create a new block using previous block's hash
func generateBlock(prevBlock Block, transaction Transaction) Block {

  now := time.Now()

  BaseDifficulty := uint32(4026531840)

  newBlockHeader := BlockHeader{}
  newBlockHeader = BlockHeader{1, calculateBlockHeaderHash(prevBlock.BlockHeader), calculateTransactionHash(transaction), now.Unix(), BaseDifficulty, 0}
  newBlockHeader.Nonce = mineBlockNonce(newBlockHeader)
  log.Printf("New Block Mined with Header Hash %08x", calculateBlockHeaderHash(newBlockHeader))

  newBlock := Block{}
  newBlock = Block{0, newBlockHeader, transaction}
  newBlock.BlockSize = uint(unsafe.Sizeof(newBlock))

  return newBlock
}

//TODO: Verboseness
func isBlockValid(newBlock Block, prevBlock Block, idx int) bool {
  MaxUint := ^uint32(0)
  threshhold := MaxUint - prevBlock.BlockHeader.Difficulty

  log.Printf("Checking if block %d is Valid\n", idx)

  printBlock(newBlock)

  // Difficulty check
  if calculateBlockHeaderHash(newBlock.BlockHeader) >= threshhold {
    log.Printf("Header hash %d not less than threshold %d", calculateBlockHeaderHash(newBlock.BlockHeader), threshhold)
    return false
  }

  // Prev Hash check
  if newBlock.BlockHeader.PrevBlockHash != calculateBlockHeaderHash(prevBlock.BlockHeader) {
    log.Printf("Prev\n")
    printBlock(prevBlock)
    log.Printf("Prev Header hash %d not correct %d", newBlock.BlockHeader.PrevBlockHash, calculateBlockHeaderHash(prevBlock.BlockHeader))
    return false
  }

  // Transaction hash check
  if calculateTransactionHash(newBlock.Transaction) != newBlock.BlockHeader.TransactionHash {
    log.Printf("Transaction hash %d not correct %d", calculateTransactionHash(newBlock.Transaction), newBlock.BlockHeader.TransactionHash)
    return false
  }

  // Transaction check

  return true
}

func verifyBlockchain(chain []Block) bool {
  log.Println("Checking if chain is Valid\n")
  for i := 1; i < len(chain); i++ {
    if !isBlockValid(chain[i], chain[i-1], i) {
      log.Printf("Block %d is invalid", i)
      return false
    }
  }
  return true
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

  // If the seed is zero, use real cryptographic randomness. Otherwise, use a
  // deterministic randomness source to make generated keys stay the same
  // across multiple runs
  var r io.Reader
  if randseed == 0 {
    r = rand.Reader
  } else {
    r = mrand.New(mrand.NewSource(randseed))
  }

  // Generate a key pair for this host. We will use it
  // to obtain a valid host ID.
  priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
  if err != nil {
    return nil, err
  }

  opts := []libp2p.Option{
    libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
    libp2p.Identity(priv),
  }

  basicHost, err := libp2p.New(opts...)
  if err != nil {
    return nil, err
  }

  // Build host multiaddress
  hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

  // Now we can build a full multiaddress to reach this host
  // by encapsulating both addresses:
  addrs := basicHost.Addrs()
  var addr ma.Multiaddr
  // select the address starting with "ip4"
  for _, i := range addrs {
    if strings.HasPrefix(i.String(), "/ip4") {
      addr = i
      break
    }
  }
  fullAddr := addr.Encapsulate(hostAddr)
  log.Printf("I am %s\n", fullAddr)
  if secio {
    log.Printf("Now run \"go run main.go --port %d --peer %s --secio\" on a different terminal\n", listenPort+1, fullAddr)
  } else {
    log.Printf("Now run \"go run main.go --port %d --peer %s\" on a different terminal\n", listenPort+1, fullAddr)
  }

  return basicHost, nil
}

func handleStream(s net.Stream) {

  log.Println("Got a new stream!")
  inboundPeers.Inc()

  // Create a buffer stream for non blocking read and write.
  rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

  go readData(rw)
  go writeData(rw)

  // stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {

  for {
    str, err := rw.ReadString('\n')
    if err != nil {
      log.Fatal(err)
    }

    if str == "" {
      return
    }
    if str != "\n" {

      var ledger Ledger
      if err := json.Unmarshal([]byte(str), &ledger); err != nil {
        log.Fatal(err)
      }

      mutex.Lock()
      //TODO: Verify blockchain diff instead?
      //TODO: Can someone overwrite all chain?
      chain := ledger.Blockchain
      if len(chain) > len(TheLedger.Blockchain) {
        if verifyBlockchain(chain) {
          TheLedger = ledger
          blockheight.Set(float64(len(TheLedger.Blockchain)))

          storeChainToFile(TheLedger)
        }
      }
      mutex.Unlock()
    }
  }
}

func isTransactionValid(transaction Transaction) bool {
  // Enough funds
  if transaction.Amount + transaction.Fee > TheLedger.AccountMap[transaction.FromAddress] {
    log.Printf("Transacting more than owned %d", TheLedger.AccountMap[transaction.FromAddress])
    return false
  }

  return true
}

func writeData(rw *bufio.ReadWriter) {

  go func() {
    for {
      time.Sleep(5 * time.Second)
      mutex.Lock()
      bytes, err := json.Marshal(TheLedger)
      if err != nil {
        log.Println(err)
      }
      mutex.Unlock()

      mutex.Lock()
      rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
      rw.Flush()
      mutex.Unlock()

    }
  }()

  stdReader := bufio.NewReader(os.Stdin)

  for {
    log.Print("\n> ")
    sendData, err := stdReader.ReadString('\n')
    if err != nil {
      log.Fatal(err)
    }

    sendData = strings.Replace(sendData, "\n", "", -1)

    if sendData == "print" {
      printChain(TheLedger.Blockchain)
      printAccounts(5)
      continue
    } else if sendData == "exit" || sendData == "quit" {
      log.Println("Exiting...")
      os.Exit(0)
      continue
    } else if sendData == "help" || sendData == "h" {
      log.Println("Options include : help, print, exit, <int>transaction-data")
      continue
    }

    //TODO: Check if valid
    values := strings.Split(sendData, ",")

    // Assign the values to individual variables
    var from, to, amount, fee uint64
    from, err = strconv.ParseUint(values[0], 10, 64)
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

    transaction := Transaction{from, to, amount, fee, make([]byte, 0)}
    signTransaction(&transaction)
    Mempool = append(Mempool, transaction)
    printMempool()
    newBlock := generateBlock(TheLedger.Blockchain[len(TheLedger.Blockchain)-1], transaction)

    if isBlockValid(newBlock, TheLedger.Blockchain[len(TheLedger.Blockchain)-1], len(TheLedger.Blockchain)) && isTransactionValid(transaction) {
      mutex.Lock()
      TheLedger.Blockchain = append(TheLedger.Blockchain, newBlock)
      blockheight.Set(float64(len(TheLedger.Blockchain)))

      MinerAddress := uint64(0)
      TheLedger.AccountMap[transaction.FromAddress] -= transaction.Fee
      TheLedger.AccountMap[MinerAddress] += transaction.Fee
      TheLedger.AccountMap[transaction.FromAddress] -= transaction.Amount
      TheLedger.AccountMap[transaction.ToAddress] += transaction.Amount

      mutex.Unlock()
    }

    bytes, err := json.Marshal(TheLedger)
    if err != nil {
      log.Println(err)
    }

    printChain(TheLedger.Blockchain)
    printAccounts(5)

    mutex.Lock()
    rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
    rw.Flush()
    mutex.Unlock()
  }

}

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

func printChain(blockchain []Block) {
  MaxUint := ^uint32(0)

  log.Println("Blockchain :")
  for i := 0; i < len(blockchain); i++ {
    log.Printf("Block %d", i)
    log.Printf("    Block Header Hash : 0x%08x", calculateBlockHeaderHash(blockchain[i].BlockHeader))
    log.Printf("    Prev Block Hash : 0x%08x", blockchain[i].BlockHeader.PrevBlockHash)
    log.Printf("    Timestamp : %s", time.Unix(blockchain[i].BlockHeader.Timestamp, 0))
    log.Printf("    Difficulty : 0x%08x", MaxUint - blockchain[i].BlockHeader.Difficulty)
    log.Printf("    Nonce %d", blockchain[i].BlockHeader.Nonce)
    log.Printf("    Transaction %d -> %d : %d", blockchain[i].Transaction.FromAddress, blockchain[i].Transaction.ToAddress, blockchain[i].Transaction.Amount)
    log.Printf("    Transaction Hash : 0x%08x", calculateTransactionHash(blockchain[i].Transaction))
  }
}

func printBlock(block Block) {
  MaxUint := ^uint32(0)
  log.Printf("    Block Header Hash : 0x%08x", calculateBlockHeaderHash(block.BlockHeader))
  log.Printf("    Prev Block Hash : 0x%08x", block.BlockHeader.PrevBlockHash)
  log.Printf("    Timestamp : %s", time.Unix(block.BlockHeader.Timestamp, 0))
  log.Printf("    Difficulty : 0x%08x", MaxUint - block.BlockHeader.Difficulty)
  log.Printf("    Nonce %d", block.BlockHeader.Nonce)
  log.Printf("    Transaction %d -> %d : %d", block.Transaction.FromAddress, block.Transaction.ToAddress, block.Transaction.Amount)
  log.Printf("    Transaction Hash : 0x%08x", calculateTransactionHash(block.Transaction))
  //TODO: Transaction signature
}

func printMempool() {
  for i, t := range Mempool {
    fmt.Printf("Transaction %d:\n", i+1)
    fmt.Printf("  From Address: %d\n", t.FromAddress)
    fmt.Printf("  To Address: %d\n", t.ToAddress)
    fmt.Printf("  Amount: %d\n", t.Amount)
    fmt.Printf("  Fee: %d\n", t.Fee)
    fmt.Printf("  Signature: %x\n", t.Signature)
  }
}

func printBlockHeader(blockHeader BlockHeader) {
  MaxUint := ^uint32(0)
  log.Printf("Printing blockHeader")
  log.Printf("    Version %d", blockHeader.Version)
  log.Printf("    Prev Block Hash : 0x%08x", blockHeader.PrevBlockHash)
  log.Printf("    Transaction Hash : 0x%08x", blockHeader.TransactionHash)
  log.Printf("    Timestamp : %s", time.Unix(blockHeader.Timestamp, 0))
  log.Printf("    Difficulty : 0x%08x", MaxUint - blockHeader.Difficulty)
  log.Printf("    Nonce %d", blockHeader.Nonce)
}

func printAccounts(count uint64) {
  log.Println("Accounts :")
  for i := uint64(0); i < count; i++ {
    log.Printf("Account %d : Balance %d", i, TheLedger.AccountMap[i])
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
func Load(path string) error {
  log.Print("Loading Ledger...")
  snapshotLoading.Set(1.0)
  lock.Lock()
  defer lock.Unlock()
  f, err := os.Open(path)
  if err != nil {
    return err
  }
  defer f.Close()
  var checkLedger Ledger
  if err:= Unmarshal(f, &checkLedger); err != nil {
    return err
  }

  log.Print("Verifying Ledger...")
  printChain(checkLedger.Blockchain)
  for i := uint64(0); i < 5; i++ {
    log.Printf("Account %d : Balance %d", i, checkLedger.AccountMap[i])
  }

  snapshotLoading.Set(2.0)
  if verifyBlockchain(checkLedger.Blockchain) {
    //TODO: FOr all spots and refactor
    log.Print("Verified!")
    TheLedger = checkLedger
    blockheight.Set(float64(len(TheLedger.Blockchain)))
  }

  //TODO: Load initial state cause not valid

  snapshotLoading.Set(3.0)
  return nil
}

func storeChainToFile(ledger Ledger) {
  log.Println("Saving snapshot to file")
  if err := Save("./snapshot.tmp", ledger); err != nil {
    log.Fatalln(err)
  }
}

func loadChainFromFile(snapPath string) {
  log.Println("Loading snapshot from file")
  if err := Load(snapPath); err != nil {
    log.Fatalln(err)

  }
}

var (
	// Create a Prometheus metric to track the number of requests
	requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests",
	})
	// Create a Prometheus metric to track the duration of requests
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Duration of requests",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 20),
	})
	// Create a Prometheus metric with block height
	blockheight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "block_height",
		Help: "Blockchain current height",
	})
	// Create a Prometheus metric to track number of blocks mined
	blocksMined = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "blocks_mined",
		Help: "Total number of blocks mined",
	})
	// Create a Prometheus metric with status of loading snapshot
	snapshotLoading = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "snapshot_loading",
    Help: "Snapshot loading status : 0 - Pre / No Snapshot, 1 - Loading Snapshot, 2 - Verifying Snapshot, 3 - Done",
	})
	// Create a Prometheus metric to track the number of inbound peers
	inboundPeers = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "inbound_peers",
    Help: "Inbound Peers to node",
	})
  // Total value transfered? TransactionCount?
)

func promSetup() {
  // Register the metrics with Prometheus
	prometheus.MustRegister(requests)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(blockheight)
	prometheus.MustRegister(blocksMined)
	prometheus.MustRegister(snapshotLoading)
	prometheus.MustRegister(inboundPeers)
}

func airdrop(accounts map[uint64]uint64, airdropFile string) {
  // One Trillion
  TotalSupply := uint64(1000000000000)

  if airdropFile != "" {
    f, err := os.Open(airdropFile)
    if err != nil {
      return
    }
    var total float64
    defer f.Close()
    scanner := bufio.NewScanner(f)
	  for scanner.Scan() {
	  	line := scanner.Text()
	  	parts := strings.Split(line, ",")
	  	id, err := strconv.Atoi(parts[0])
	  	if err != nil {
	  		fmt.Println(err)
	  		return
	  	}
	  	value, err := strconv.ParseFloat(parts[1], 64)
	  	if err != nil {
	  		fmt.Println(err)
	  		return
	  	}
      total += value
	  	fmt.Println("ID:", id, "Value:", value)
      accounts[uint64(id)] = uint64(float64(TotalSupply) * value)
	  }

	  if err := scanner.Err(); err != nil {
	  	fmt.Println(err)
	  }
    if total != 1.0 {
	  	fmt.Println("Error: the total of all values is not 1.0")
	  } else {
	  	fmt.Println("Success: the total of all values is 1.0")
	  }
  } else {
    for i := uint64(0); i < 5; i++ {
      accounts[i] = TotalSupply / 5
    }
  }
}

//TODO: Multi peer & learn p2p
func main() {
	golog.SetAllLoggers(golog.LevelInfo) // Change to DEBUG for extra info

  // Genesis Account Map
  accounts := make(map[uint64]uint64)

  // Genesis Block Creation
  MaxUint := ^uint32(0)
  BaseDifficulty := uint32(4026531840)
  log.Printf("Hashes must be less than 0x%08x", MaxUint - BaseDifficulty)

  // Parse options from the command line
  listenPort := flag.Int("port", 0, "Port to listen for connections")
  peerToCall := flag.String("peer", "", "Peer port / path to dial")
  rpcPort := flag.String("rpc", "", "RPC port / prom metrics port")
  secio := flag.Bool("secio", false, "enable secio")
  seed := flag.Int64("seed", 0, "Seed for id generation")
  snapPath := flag.String("snap", "", "Load blockchain from snapshot")
  airdropFile := flag.String("air", "", "Airdrop addresses")
  keyFile := flag.String("key", "", "Private Key")
  flag.Parse()

  // Airdrop to accounts
  airdrop(accounts, *airdropFile)
  
  //TOOD: Do I need to provide a public or just private?
  if *keyFile == "" {
    generatePublicPrivateKey()
    //TODO: Save filename
    //TODO: Maybe just use account number and abstract pub key away? allow backup tho
  }

  if *snapPath != "" {
    loadChainFromFile(*snapPath)

    printChain(TheLedger.Blockchain)
    printAccounts(5)
  } else {
    now := time.Now()

    genesisBlockTransaction := Transaction{0, 1, 0, 0, make([]byte, 0)}
    signTransaction(&genesisBlockTransaction) //TODO: Sign from what address?
    genesisBlockHeader := BlockHeader{}
    genesisBlockHeader = BlockHeader{1, 0, calculateTransactionHash(genesisBlockTransaction), now.Unix(), BaseDifficulty, 0}
    genesisBlockHeader.Nonce = mineBlockNonce(genesisBlockHeader)
    log.Printf("Genesis Block Header Hash 0x%08x", calculateBlockHeaderHash(genesisBlockHeader))

    genesisBlock := Block{}
    genesisBlock = Block{0, genesisBlockHeader, genesisBlockTransaction}
    genesisBlock.BlockSize = uint(unsafe.Sizeof(genesisBlock))

    //TODO: Ledger transfer value and refactor
    TheLedger.Blockchain = append(TheLedger.Blockchain, genesisBlock)
    blockheight.Set(float64(len(TheLedger.Blockchain)))

    TheLedger.AccountMap = accounts
    TheLedger.AccountMap[genesisBlockTransaction.FromAddress] -= genesisBlockTransaction.Amount
    TheLedger.AccountMap[genesisBlockTransaction.ToAddress] += genesisBlockTransaction.Amount

    printChain(TheLedger.Blockchain)
    printAccounts(5)
  }

  if *listenPort == 0 {
    log.Fatal("You must specify a port to bind the client on. Use the 'port' flag")
  }

  // Setup P2P
  ha, err := makeBasicHost(*listenPort, *secio, *seed)
  if err != nil {
    log.Fatal(err)
  }

  if *peerToCall == "" {
    log.Println("listening for connections")
    // Set a stream handler on host A. /p2p/1.0.0 is
    // a user-defined protocol name.
    ha.SetStreamHandler("/p2p/1.0.0", handleStream)

  } else {
    ha.SetStreamHandler("/p2p/1.0.0", handleStream)

    // The following code extracts target's peer ID from the
    // given multiaddress
    ipfsaddr, err := ma.NewMultiaddr(*peerToCall)
    if err != nil {
      log.Fatalln(err)
    }

    pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
    if err != nil {
      log.Fatalln(err)
    }

    peerid, err := peer.Decode(pid)
    if err != nil {
      log.Fatalln(err)
    }

    // Decapsulate the /ipfs/<peerID> part from the target
    // /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
    targetPeerAddr, _ := ma.NewMultiaddr(
      fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
    targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

    // We have a peer ID and a targetAddr so we add it to the peerstore
    // so LibP2P knows how to contact it
    ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

    log.Println("opening stream")
    // make a new stream from host B to host A
    // it should be handled on host A by the handler we set above because
    // we use the same /p2p/1.0.0 protocol
    s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
    if err != nil {
      log.Fatalln(err)
    }
    // Create a buffered stream so that read and writes are non blocking.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

    // Create a thread to read and write data.
    go writeData(rw)
    go readData(rw)

  }

  promSetup()
  // Create a new HTTP server to handle metrics requests
  // TODO: To jsonrpc
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		requests.Inc()
    w.Write([]byte("The exposed endpoints are : metrics, help, add, block\n"))
  })
	http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		requests.Inc()
		defer requestDuration.Observe(time.Since(time.Now()).Seconds())

    value := r.URL.Query().Get("value")
    i, err := strconv.Atoi(value)
    if err != nil {
      log.Fatal(err)
    }
    if i >= len(TheLedger.Blockchain) {

      w.Write([]byte(fmt.Sprintf("Blockchain is of length : %d so cannot index position : %d\n", len(TheLedger.Blockchain), i)))
    } else {
      message := fmt.Sprintf("Block %d\n", i)
      message += fmt.Sprintf("    Block Header Hash : 0x%08x\n", calculateBlockHeaderHash(TheLedger.Blockchain[i].BlockHeader))
      message += fmt.Sprintf("    Prev Block Hash : 0x%08x\n", TheLedger.Blockchain[i].BlockHeader.PrevBlockHash)
      message += fmt.Sprintf("    Timestamp : %s\n", time.Unix(TheLedger.Blockchain[i].BlockHeader.Timestamp, 0))
      message += fmt.Sprintf("    Difficulty : 0x%08x\n", MaxUint - TheLedger.Blockchain[i].BlockHeader.Difficulty)
      message += fmt.Sprintf("    Nonce %d\n", TheLedger.Blockchain[i].BlockHeader.Nonce)
      message += fmt.Sprintf("    Transaction %d -> %d : %d", TheLedger.Blockchain[i].Transaction.FromAddress, TheLedger.Blockchain[i].Transaction.ToAddress, TheLedger.Blockchain[i].Transaction.Amount)
      w.Write([]byte(message))
    }
  })
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		// Track the number of requests
		requests.Inc()
		// Track the duration of the request
		defer requestDuration.Observe(time.Since(time.Now()).Seconds())
    value := r.URL.Query().Get("value")
    values := strings.Split(value, ",")

    // Assign the values to individual variables
    var from, to, amount, fee uint64
    from, err = strconv.ParseUint(values[0], 10, 64)
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
    transaction := Transaction{from, to, amount, fee, make([]byte, 0)}
    signTransaction(&transaction)
    Mempool = append(Mempool, transaction)
    printMempool()
    newBlock := generateBlock(TheLedger.Blockchain[len(TheLedger.Blockchain)-1], transaction)

    if isBlockValid(newBlock, TheLedger.Blockchain[len(TheLedger.Blockchain)-1], len(TheLedger.Blockchain)) {
      mutex.Lock()
      TheLedger.Blockchain = append(TheLedger.Blockchain, newBlock)
      blockheight.Set(float64(len(TheLedger.Blockchain)))

      MinerAddress := uint64(0)
      TheLedger.AccountMap[transaction.FromAddress] -= transaction.Fee
      TheLedger.AccountMap[MinerAddress] += transaction.Fee
      TheLedger.AccountMap[transaction.FromAddress] -= transaction.Amount
      TheLedger.AccountMap[transaction.ToAddress] += transaction.Amount
      mutex.Unlock()
    }

    printChain(TheLedger.Blockchain)
    printAccounts(5)
	})

  log.Printf("Serving RPC at %s", *rpcPort)
	http.ListenAndServe(":" + *rpcPort, nil)

  select {} // hang forever
}

//TODO: Test suite
func generatePublicPrivateKey() {
	// Account number for the client
	accountNumber := 12345

	// Generate a new RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the public key in the right format
	publicKey := privateKey.PublicKey

	// Save the private key to a PEM file
	privateKeyFile, err := os.Create(fmt.Sprintf("private_%d.pem", accountNumber))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		fmt.Println(err)
		return
	}

	// Save the public key to a PEM file
	publicKeyFile, err := os.Create(fmt.Sprintf("public_%d.pem", accountNumber))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer publicKeyFile.Close()

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := pem.Encode(publicKeyFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyPEM}); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("RSA Key pair generated for account number %d\n", accountNumber)
	fmt.Printf("Private key saved to private_%d.pem\n", accountNumber)
	fmt.Printf("Public key saved to public_%d.pem\n", accountNumber)
}
