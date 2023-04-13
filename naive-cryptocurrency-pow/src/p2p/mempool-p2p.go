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

	"naive-cryptocurrency-pow/src/mempool"
	"naive-cryptocurrency-pow/src/metrics"
)

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicMempoolHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

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
    libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
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
  log.Printf("Mempool located at %s\n", fullAddr)

  return basicHost, nil
}

func handleMempoolStream(s net.Stream) {

  log.Println("Got a new stream!")
  metrics.InboundPeers.Inc()

  // Create a buffer stream for non blocking read and write.
  rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

  go readMempoolData(rw)
  go writeMempoolData(rw)

  // stream 's' will stay open until you close it (or the other side closes it).
}

func readMempoolData(rw *bufio.ReadWriter) {

  for {
    str, err := rw.ReadString('\n')
    if err != nil {
      log.Fatal(err)
    }

    if str == "" {
      return
    }
    if str != "\n" {
      var inMemPool mempool.Mempool
      if err := json.Unmarshal([]byte(str), &inMemPool); err != nil {
        log.Fatal(err)
      }

      log.Println("Received Mempool Data")
      inMemPool.PrintMempool()

      mempool.TheMempool.Mutex.Lock()

      for _, tx := range inMemPool.AvailableTransactions {
        //TODO: Remove transactions that are already in the ledger
          //if ledger.VerifyTransaction(tx) {
          mempool.TheMempool.AddTransaction(tx, true)
          //}
      }

      mempool.TheMempool.Mutex.Unlock()
    }
  }
}

func writeMempoolData(rw *bufio.ReadWriter) {

  go func() {
    for {
      time.Sleep(5 * time.Second)
      mempool.TheMempool.Mutex.Lock()
      bytes, err := json.Marshal(mempool.TheMempool)
      if err != nil {
        log.Println(err)
      }
      mempool.TheMempool.Mutex.Unlock()

      mempool.TheMempool.Mutex.Lock()
      rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
      rw.Flush()
      mempool.TheMempool.Mutex.Unlock()

    }
  }()

}

func DialMempoolPeer(listenPort int, secio bool, randseed int64, peerFilename string) {
  // Setup P2P
  ha, err := makeBasicMempoolHost(listenPort, secio, randseed)
  if err != nil {
    log.Fatal(err)
  }

  file, err := os.Open(peerFilename)
  if err != nil {
    log.Println("No peer file found")
    log.Fatal(err)
  }
  defer file.Close()

  var peerToCall string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    peerToCall = scanner.Text()
  }

  if peerToCall == "" {
    log.Println("listening for connections")
    // Set a stream handler on host A. /p2p/1.0.0 is
    // a user-defined protocol name.
    ha.SetStreamHandler("/p2p/1.0.0", handleMempoolStream)

  } else {
    ha.SetStreamHandler("/p2p/1.0.0", handleMempoolStream)

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
    go writeMempoolData(rw)
    go readMempoolData(rw)

  }
}
