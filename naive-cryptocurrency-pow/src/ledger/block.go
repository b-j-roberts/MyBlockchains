package ledger

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"time"

	unsafe "unsafe"
)

// TODO: For now a block contains one datapoint, in the future change to multiple transactions
// TODO: match closer https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch07.html
type BlockHeader struct {
  Version int
  PrevBlockHash uint32
  TransactionHash uint32
  Timestamp int64
  Difficulty uint32 // TODO: How decided & not changed to dupe system to be faster
  Nonce int
}

type Block struct {
  BlockSize uint
  BlockHeader BlockHeader
  Transactions []Transaction
}

func CreateUnminedBlock(prevBlock Block, transactions []Transaction, blockheight uint32) (Block, error) {
  for _, tx := range transactions {
    if !VerifyTransaction(tx) {
      log.Printf("Transaction invalid")
      return Block{}, errors.New("Transaction invalid")
    }
  }

  version := 1
  difficulty := computeDifficulty(blockheight)

  prevBlockHash := CalculateBlockHeaderHash(prevBlock.BlockHeader)

  // Transaction hash
  transactionsHashString := ""
  for _, transaction := range transactions {
    transactionsHashString += string(CalculateTransactionHash(transaction))
  }
  txHash := DoubleHash(transactionsHashString)

  now := time.Now()

  //TODO: Migrate to json-esk format?
  newBlockHeader := BlockHeader{version, prevBlockHash, txHash, now.Unix(), difficulty, 0}

  newBlock := Block{}
  newBlock.BlockHeader = newBlockHeader
  newBlock.Transactions = transactions
  newBlock.BlockSize = uint(unsafe.Sizeof(newBlock))

  return newBlock, nil
}

func CalculateBlockHeaderHash(blockHeader BlockHeader) uint32 {
  data, err := json.Marshal(blockHeader)
  if err != nil {
    panic(err)
  }

  hasher := sha256.New()
  hasher.Write(data)
  hashed := hasher.Sum(nil)

  return uint32(binary.BigEndian.Uint32(hashed))
}

func DoubleHash(hash string) uint32 {
  hasher := sha256.New()
  hasher.Write([]byte(hash))
  hashed := hasher.Sum(nil)

  hasher = sha256.New()
  hasher.Write(hashed)
  hashed = hasher.Sum(nil)

  return uint32(binary.BigEndian.Uint32(hashed))
}

//TODO: Should I pull in prevBlock from ledger?
func IsBlockValid(newBlock Block, prevBlock Block, blockHeight uint32) bool {
  threshhold := computeDifficulty(blockHeight)

  printBlock(newBlock)

  // Difficulty check
  if CalculateBlockHeaderHash(newBlock.BlockHeader) >= threshhold {
    log.Printf("Header hash %d not less than threshold %d", CalculateBlockHeaderHash(newBlock.BlockHeader), threshhold)
    return false
  }

  // Prev Hash check
  if newBlock.BlockHeader.PrevBlockHash != CalculateBlockHeaderHash(prevBlock.BlockHeader) {
    log.Printf("Prev Header hash %d not correct %d", newBlock.BlockHeader.PrevBlockHash, CalculateBlockHeaderHash(prevBlock.BlockHeader))
    return false
  }

  if len(newBlock.Transactions) == 0 || len(newBlock.Transactions) > 100 {
    log.Printf("Transaction count %d not correct, must be 1-100", len(newBlock.Transactions))
    return false
  }

  // Transaction hash check
  transactionsHashString := ""
  for _, transaction := range newBlock.Transactions {
    transactionsHashString += string(CalculateTransactionHash(transaction))
  }
  if DoubleHash(transactionsHashString) != newBlock.BlockHeader.TransactionHash {
    log.Printf("Transaction hash %d not correct %d", newBlock.BlockHeader.TransactionHash, DoubleHash(transactionsHashString))
    return false
  }

  // Transaction check
  for _, transaction := range newBlock.Transactions {
    if !VerifyTransaction(transaction) {
      log.Printf("Transaction invalid")
      return false
    }
  }

  return true
}

func VerifyBlockchain(chain []Block) bool {
  log.Println("Verifying the Blockchain is valid\n")
  for i := 1; i < len(chain); i++ {
    if !IsBlockValid(chain[i], chain[i-1], uint32(i)) {
      log.Printf("Block %d is invalid", i)
      return false
    }
  }

  //TODO: Verify account balances and nonces correct
  return true
}

func createGenesisBlock(transaction Transaction) Block {
  version := 1
  difficulty := computeDifficulty(0)

  prevBlockHash := uint32(0)
  txHash := CalculateTransactionHash(transaction)
  now := time.Now()

  genesisBlockHeader := BlockHeader{version, prevBlockHash, txHash, now.Unix(), difficulty, 0}

  //TODO: Use mining function
  genesisBlockHeader.Nonce = 0
  //genesisBlockHeader.Nonce = mineBlockNonce(genesisBlockHeader)
  genesisBlock := Block{}
  genesisBlock.BlockHeader = genesisBlockHeader
  genesisBlock.Transactions = []Transaction{transaction}
  genesisBlock.BlockSize = uint(unsafe.Sizeof(genesisBlock))

  return genesisBlock
}

//TODO: to member String
func printBlock(block Block) {
  PrintBlockHeader(block.BlockHeader)
  for _, transaction := range block.Transactions {
    PrintTransaction(transaction)
  }
}

func PrintBlockHeader(blockHeader BlockHeader) {
  MaxUint := ^uint32(0)
  log.Printf("Printing blockHeader")
  log.Printf("    Version %d", blockHeader.Version)
  log.Printf("    Prev Block Hash : 0x%08x", blockHeader.PrevBlockHash)
  log.Printf("    Transaction Hash : 0x%08x", blockHeader.TransactionHash)
  log.Printf("    Timestamp : %s", time.Unix(blockHeader.Timestamp, 0))
  log.Printf("    Difficulty : 0x%08x", MaxUint - blockHeader.Difficulty)
  log.Printf("    Nonce %d", blockHeader.Nonce)
}
