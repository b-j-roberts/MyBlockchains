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
	"sync"
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

//TODO: Refactor to loop going thru peers rw buffers and querying them for mempool data?

type Peer struct {
  PeerID peer.ID
  PeerAddr ma.Multiaddr
  Reputation uint
}

//TODO: To MempoolHost
type MempoolPeer struct {
//  PeerID peer.ID
//  PeerAddr ma.Multiaddr
  PeerCount int
  Host host.Host
  PeerList []Peer
}

var TheMempoolPeer MempoolPeer

func InitMempool(listenPort int, secio bool, randseed int64, peerFilename string) {
  TheMempoolPeer.PeerList = make([]Peer, 0)
  SetupMempoolHost(listenPort, secio, randseed)
  DialMempoolPeers(peerFilename)
}

func SetupMempoolHost(listenPort int, secio bool, randseed int64) {
  // Setup P2P                                                         
  var err error
  TheMempoolPeer.Host, err = makeBasicMempoolHost(listenPort, secio, randseed)    
  if err != nil {                                                 
    log.Fatal(err)                                                
  }                                                               
                                                                  
  log.Println("listening for connections")    
  // Set a stream handler on host A. /p2p/1.0.0 is
  // a user-defined protocol name.            
  TheMempoolPeer.Host.SetStreamHandler("/p2p/1.0.0", handleMempoolStream)
}

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

  //TODO: Dial mempool peer?
  wg := sync.WaitGroup{}
  wg.Add(2)
  go readMempoolData(rw, &wg)
  go writeMempoolData(rw, &wg)
  wg.Wait()

  // stream 's' will stay open until you close it (or the other side closes it).
}
//TODO: Peer object w/ array of peers connected to with info about them for peer discovery and gossiping, should I read from peers that sub to me?

func readMempoolData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
  defer wg.Done()

  for {
    str, err := rw.ReadString('\n')
    if err != nil {
      log.Println("Error reading from buffer")
      return
    }

    if str == "" {
      log.Println("Empty string received. Closing connection.")
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

func writeMempoolData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
  defer wg.Done()

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
      val, err := rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
      if err != nil {
        log.Println("Error writing to buffer")
        return
      }
      log.Println("Wrote", val, "bytes to buffer")
      err = rw.Flush()
      if err != nil {
        log.Println("Error flushing buffer")
        return
      }

      mempool.TheMempool.Mutex.Unlock()

    }
  }()

}

func DialMempoolPeers(peerFilename string) {
  f, err := os.Open(peerFilename)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  peerCount := 0
  peersToDial := make([]string, 0)
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    // Ignore lines with # or // at the beginning
    if strings.HasPrefix(scanner.Text(), "#") || strings.HasPrefix(scanner.Text(), "//") {
      continue
    }

    peerToCall := scanner.Text()
    peersToDial = append(peersToDial, peerToCall)
    peerCount++
  }


  if peerCount > 0 {
    log.Println("Dialing", peerCount, "peers")
    var wg sync.WaitGroup
    wg.Add(peerCount)
    for _, peerToCall := range peersToDial {
      go DialMempoolPeer(peerToCall, &wg)
    }
    wg.Wait()
  }

  if TheMempoolPeer.PeerCount <= 0 {
    log.Println("WARNING: No peers dialed")
  }
}


//func DialMempoolPeer(ha *host.Host, peerToCall string) {
func DialMempoolPeer(peerToCall string, wg *sync.WaitGroup) {
  defer wg.Done()
  log.Println("Dialing Mempool Peer", peerToCall)

  //TODO: var newPeer Peer
  if peerToCall != "" {
    //TheMempoolPeer.Host.SetStreamHandler("/p2p/1.0.0", handleMempoolStream)

    // The following code extracts target's peer ID from the
    // given multiaddress
    ipfsaddr, err := ma.NewMultiaddr(peerToCall)
    if err != nil {
      log.Println("Error parsing multiaddress for peer")
      return
    }

    pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
    if err != nil {
      log.Println("Error parsing peer ID from multiaddress")
      return
    }

    peerid, err := peer.Decode(pid)
    if err != nil {
      log.Println("Error decoding peer ID")
      return
    }

    // Decapsulate the /ipfs/<peerID> part from the target
    // /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
    targetPeerAddr, _ := ma.NewMultiaddr(
      fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
    targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

    // We have a peer ID and a targetAddr so we add it to the peerstore
    // so LibP2P knows how to contact it
    TheMempoolPeer.Host.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

    log.Println("opening stream")
    // make a new stream from host B to host A
    // it should be handled on host A by the handler we set above because
    // we use the same /p2p/1.0.0 protocol
    s, err := TheMempoolPeer.Host.NewStream(context.Background(), peerid, "/p2p/1.0.0")
    if err != nil {
      log.Println("Error opening stream")
      return
    }
    // Create a buffered stream so that read and writes are non blocking.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
    TheMempoolPeer.PeerCount++

    var wgInner sync.WaitGroup
    wgInner.Add(2)
    // Create a thread to read and write data.
    go writeMempoolData(rw, &wgInner)
    go readMempoolData(rw, &wgInner)
    wgInner.Wait()

    TheMempoolPeer.PeerCount--
  }
}
