package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mrand "math/rand"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	net "github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	pstore "github.com/libp2p/go-libp2p/core/peerstore"

	"naive-cryptocurrency-pow/src/ledger"
	"naive-cryptocurrency-pow/src/mempool"
	"naive-cryptocurrency-pow/src/metrics"
)

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
  metrics.InboundPeers.Inc()

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
      var inLedger ledger.Ledger
      if err := json.Unmarshal([]byte(str), &inLedger); err != nil {
        log.Fatal(err)
      }

      ledger.TheLedger.Mutex.Lock()
      //TODO: Verify blockchain diff instead?
      //TODO: Can someone overwrite all chain?
      chain := inLedger.Blockchain
      if len(chain) > len(ledger.TheLedger.Blockchain) {
        if ledger.VerifyBlockchain(chain) {
          ledger.TheLedger = inLedger
          metrics.Blockheight.Set(float64(len(ledger.TheLedger.Blockchain)))

          ledger.StoreChainToFile(ledger.TheLedger)
        }
      }
      ledger.TheLedger.Mutex.Unlock()
    }
  }
}

func writeData(rw *bufio.ReadWriter) {

  go func() {
    for {
      time.Sleep(5 * time.Second)
      ledger.TheLedger.Mutex.Lock()
      bytes, err := json.Marshal(ledger.TheLedger)
      if err != nil {
        log.Println(err)
      }
      ledger.TheLedger.Mutex.Unlock()

      ledger.TheLedger.Mutex.Lock()
      rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
      rw.Flush()
      ledger.TheLedger.Mutex.Unlock()

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
      log.Println("Printing...")
      ledger.PrintChain(ledger.TheLedger.Blockchain)
      ledger.PrintAccounts(5)
      mempool.TheMempool.PrintMempool()
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
    var from, to, amount, fee, nonce uint64
    var publicKeyFile, privateKeyFile string
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
    nonce, err = strconv.ParseUint(values[4], 10, 64)
    if err != nil {
      log.Fatal(err)
    }
    publicKeyFile = values[5]
    privateKeyFile = values[6]

    transaction, err := ledger.CreateTransaction(from, to, amount, fee, uint(nonce), privateKeyFile, publicKeyFile, true)
    if err != nil {
      log.Println("Error creating transaction")
      continue
    }

    mempool.TheMempool.AddTransaction(transaction, true)
    transactions := []ledger.Transaction{transaction}
    newBlock, err := ledger.CreateUnminedBlock(ledger.TheLedger.Blockchain[len(ledger.TheLedger.Blockchain)-1], transactions, uint32(len(ledger.TheLedger.Blockchain)))
    if err != nil {
      log.Fatal(err)
    }

    if ledger.IsBlockValid(newBlock, ledger.TheLedger.Blockchain[len(ledger.TheLedger.Blockchain)-1], uint32(len(ledger.TheLedger.Blockchain))) && ledger.VerifyTransaction(transaction) {
      ledger.TheLedger.Mutex.Lock()
      ledger.TheLedger.Blockchain = append(ledger.TheLedger.Blockchain, newBlock)
      metrics.Blockheight.Set(float64(len(ledger.TheLedger.Blockchain)))

      MinerAddress := uint64(0)
      ledger.TheLedger.AccountBalances[transaction.FromAddress] -= transaction.Fee
      ledger.TheLedger.AccountBalances[MinerAddress] += transaction.Fee
      ledger.TheLedger.AccountBalances[transaction.FromAddress] -= transaction.Amount
      ledger.TheLedger.AccountBalances[transaction.ToAddress] += transaction.Amount

      ledger.TheLedger.Mutex.Unlock()
    }

    bytes, err := json.Marshal(ledger.TheLedger)
    if err != nil {
      log.Println(err)
    }

    ledger.TheLedger.Mutex.Lock()
    rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
    rw.Flush()
    ledger.TheLedger.Mutex.Unlock()
  }

}

func DialPeer(listenPort int, secio bool, randseed int64, peerToCall string) {
  // Setup P2P
  ha, err := makeBasicHost(listenPort, secio, randseed)
  if err != nil {
    log.Fatal(err)
  }

  if peerToCall == "" {
    log.Println("listening for connections")
    // Set a stream handler on host A. /p2p/1.0.0 is
    // a user-defined protocol name.
    ha.SetStreamHandler("/p2p/1.0.0", handleStream)

  } else {
    ha.SetStreamHandler("/p2p/1.0.0", handleStream)

    // The following code extracts target's peer ID from the
    // given multiaddress
    ipfsaddr, err := ma.NewMultiaddr(peerToCall)
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
}
