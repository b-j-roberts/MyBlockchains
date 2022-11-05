// References :
// https://github.com/nosequeldeebee/blockchain-tutorial/tree/master/p2p
// https://github.com/nosequeldeebee/blockchain-tutorial/blob/master/proof-work/main.go
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html

package main

import (
  "bufio"
  "bytes"
  "context"
	"crypto/rand"
	"crypto/sha256"
  "encoding/binary"
  "encoding/gob"
  "encoding/json"
  "io"
  "flag"
  "fmt"
  "log"
  "os"
  mrand "math/rand"
  "strings"
  "strconv"
  "sync"
  "time"
  valid "github.com/asaskevich/govalidator"

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

type Block struct {
  BlockSize uint
  BlockHeader BlockHeader
  TransactionData int
}

var mutex = &sync.Mutex{}
var lock sync.Mutex

// Block Header hash -- Obtained by using sha256 twice on blockheader ( double sha256 )
//    Not stored on chain, but computed/checked each time a node receives the new block

// TODO: Store block height & hashes in indexed db somewhere?
//       How does main btc client do this? can you make calls based on blockheight?
//       NOTE : Blockheight not necessarily unique block


// When receiving block : Validate - check prev hash, ensure block header hash correct

// TODO: Double SHA256
func calculateTransactionHash(data int) uint32 {
  data_arr :=  make([]byte, 4)
  binary.BigEndian.PutUint32(data_arr, uint32(data))
  
  hasher := sha256.New()
  hasher.Write(data_arr)
  hashed := hasher.Sum(nil)

  return uint32(binary.BigEndian.Uint32(hashed))
}

func calculateBlockHeaderHash(blockHeader BlockHeader) uint32 {
  var network bytes.Buffer
  enc := gob.NewEncoder(&network)
  err := enc.Encode(blockHeader)
  if err != nil {
    panic(err)
  }

  data_arr :=  network.Bytes()
  
  hasher := sha256.New()
  hasher.Write(data_arr)
  hashed := hasher.Sum(nil)

  return uint32(binary.BigEndian.Uint32(hashed))
}

func mineBlockNonce(blockHeader BlockHeader) int {
  MaxUint := ^uint32(0)
  threshhold := MaxUint - blockHeader.Difficulty
  for {
    blockHash := calculateBlockHeaderHash(blockHeader)
    if(blockHash < threshhold) {
      log.Printf("Block mined with hash %08x", blockHash)
      break
    } else {
      log.Printf("Nonce %d failed with hash %08x", blockHeader.Nonce, blockHash)
    }
    blockHeader.Nonce += 1
  }
  return blockHeader.Nonce
}

// create a new block using previous block's hash
func generateBlock(prevBlock Block, transactionData int) Block {

  now := time.Now()

  BaseDifficulty := uint32(4026531840)

  newBlockHeader := BlockHeader{}
  newBlockHeader = BlockHeader{1, calculateBlockHeaderHash(prevBlock.BlockHeader), calculateTransactionHash(transactionData), now.Unix(), BaseDifficulty, 0}
  newBlockHeader.Nonce = mineBlockNonce(newBlockHeader)
  log.Printf("New Block Mined with Header Hash %08x", calculateBlockHeaderHash(newBlockHeader))

  newBlock := Block{}
  newBlock = Block{0, newBlockHeader, transactionData}
  newBlock.BlockSize = uint(unsafe.Sizeof(newBlock))

  return newBlock
}

//TODO: Verboseness
func isBlockValid(newBlock Block, prevBlock Block, idx int) bool {
  MaxUint := ^uint32(0)
  threshhold := MaxUint - prevBlock.BlockHeader.Difficulty

  log.Printf("Checking if block %d is Valid\n", idx)
  // Difficulty check
  if calculateBlockHeaderHash(newBlock.BlockHeader) >= threshhold {
    return false
  }

  // Prev Hash check
  if newBlock.BlockHeader.PrevBlockHash != calculateBlockHeaderHash(prevBlock.BlockHeader) {
    return false
  }

  // Transaction hash check
  if calculateTransactionHash(newBlock.TransactionData) != newBlock.BlockHeader.TransactionHash {
    return false
  }

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

      chain := make([]Block, 0)
      if err := json.Unmarshal([]byte(str), &chain); err != nil {
        log.Fatal(err)
      }

      mutex.Lock()
      //TODO: Verify blockchain diff instead?
      //TODO: Can someone overwrite all chain?
      if len(chain) > len(Blockchain) {
        if verifyBlockchain(chain) {
          Blockchain = chain

          storeChainToFile(Blockchain)
        }
      }
      mutex.Unlock()
    }
  }
}

func writeData(rw *bufio.ReadWriter) {

  go func() {
    for {
      time.Sleep(5 * time.Second)
      mutex.Lock()
      bytes, err := json.Marshal(Blockchain)
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

    if !valid.IsInt(sendData) {
      if sendData == "print" {
        printChain(Blockchain)
      } else if sendData == "exit" || sendData == "quit" {
        log.Println("Exiting...")
        os.Exit(0)
      } else if sendData == "help" || sendData == "h" {
        log.Println("Options include : help, print, exit, <int>transaction-data")
      }

      continue
    }

    data, err := strconv.Atoi(sendData)
    if err != nil {
      log.Fatal(err)
    }
    newBlock := generateBlock(Blockchain[len(Blockchain)-1], data)

    if isBlockValid(newBlock, Blockchain[len(Blockchain)-1], len(Blockchain)) {
      mutex.Lock()
      Blockchain = append(Blockchain, newBlock)
      mutex.Unlock()
    }

    bytes, err := json.Marshal(Blockchain)
    if err != nil {
      log.Println(err)
    }

    printChain(Blockchain)

    mutex.Lock()
    rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
    rw.Flush()
    mutex.Unlock()
  }

}

var Blockchain []Block

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
    log.Printf("    Block Hash : 0x%08x", calculateBlockHeaderHash(blockchain[i].BlockHeader))
    log.Printf("    Prev Block Hash : 0x%08x", blockchain[i].BlockHeader.PrevBlockHash)
    log.Printf("    Timestamp : %s", time.Unix(blockchain[i].BlockHeader.Timestamp, 0))
    log.Printf("    Difficulty : 0x%08x", MaxUint - blockchain[i].BlockHeader.Difficulty)
    log.Printf("    Nonce %d", blockchain[i].BlockHeader.Nonce)
    log.Printf("    Data %d", blockchain[i].TransactionData)
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
  lock.Lock()
  defer lock.Unlock()
  f, err := os.Open(path)
  if err != nil {
    return err
  }
  defer f.Close()
  var checkChain []Block
  if err:= Unmarshal(f, &checkChain); err != nil {
    return err
  }

  if verifyBlockchain(checkChain) {
    Blockchain = checkChain
  }

  return nil
}

func storeChainToFile(blockchain []Block) {
  log.Println("Saving snapshot to file")
  if err := Save("./snapshot.tmp", blockchain); err != nil {
    log.Fatalln(err)
  }
}

func loadChainFromFile(snapPath string) {
  log.Println("Loading snapshot from file")
  if err := Load(snapPath); err != nil {
    log.Fatalln(err)
  }
}

//TODO: Multi peer & learn p2p
func main() {
	golog.SetAllLoggers(golog.LevelInfo) // Change to DEBUG for extra info

  // Genesis Block Creation
  MaxUint := ^uint32(0)
  BaseDifficulty := uint32(4026531840)
  log.Printf("Hashes must be less than 0x%08x", MaxUint - BaseDifficulty)

  // Parse options from the command line
  listenPort := flag.Int("port", 0, "Port to listen for connections")
  peerToCall := flag.String("peer", "", "Peer port / path to dial")
  secio := flag.Bool("secio", false, "enable secio")
  seed := flag.Int64("seed", 0, "Seed for id generation")
  snapPath := flag.String("snap", "", "Load blockchain from snapshot")
  flag.Parse()

  if *snapPath != "" {
    loadChainFromFile(*snapPath)

    printChain(Blockchain)
  } else {
    now := time.Now()

    genesisBlockTransactionData := 42
    genesisBlockHeader := BlockHeader{}
    genesisBlockHeader = BlockHeader{1, 0, calculateTransactionHash(genesisBlockTransactionData), now.Unix(), BaseDifficulty, 0}
    genesisBlockHeader.Nonce = mineBlockNonce(genesisBlockHeader)
    log.Printf("Genesis Block Header Hash 0x%08x", calculateBlockHeaderHash(genesisBlockHeader))

    genesisBlock := Block{}
    genesisBlock = Block{0, genesisBlockHeader, genesisBlockTransactionData}
    genesisBlock.BlockSize = uint(unsafe.Sizeof(genesisBlock))

    Blockchain = append(Blockchain, genesisBlock)

    printChain(Blockchain)
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

    select {} // hang forever
    /**** This is where the listener code ends ****/
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

    select {} // hang forever

  }

}

//TODO: Test suite
