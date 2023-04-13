package ledger

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"naive-cryptocurrency-pow/src/signer"

	crypt "crypto"
)

type Transaction struct {
  FromAddress uint64
  ToAddress uint64
  Amount uint64
  Fee uint64
  Nonce uint
  Signature []byte
  SignerPublicKey rsa.PublicKey
  //TODO: Parent
}

//TODO: To member functions
func CreateTransaction(from uint64, to uint64, amount uint64, fee uint64, nonce uint, privateKeyFilename string, publicKeyFilename string, doVerify bool) (Transaction, error) {
  signerPublicKey, err := signer.ParsePublicKey(publicKeyFilename)
  if err != nil {
    return Transaction{}, err
  }

  transaction := Transaction{from, to, amount, fee, nonce, make([]byte, 0), *signerPublicKey}
  transaction.SignTransaction(privateKeyFilename)

  if doVerify && !VerifyTransaction(transaction) {
    return Transaction{}, fmt.Errorf("Transaction is not valid")
  }

  return transaction, nil
}

func VerifyTransaction(tx Transaction) bool {
  //TODO: verifyTransactionNonce only if adding to block, not mempool ( unless it is less than valid)
  if verifyTransactionFromSender(tx) && verifySenderHasFunds(tx) && verifyTransactionNonce(tx) { // TODO: Double check : hash?
    return true
  }

  return false
}

func CalculateTransactionHash(transaction Transaction) uint32 {

  hashed := sha256.Sum256([]byte(fmt.Sprintf("%d%d%d%d%d", transaction.FromAddress, transaction.ToAddress, transaction.Amount, transaction.Fee, transaction.Nonce)))
  uint32Hash := uint32(binary.BigEndian.Uint32(hashed[:]))


  return uint32Hash
}

func PrintTransaction(tx Transaction) {
  log.Printf("Transaction 0x%08x -> 0x%08x : %d", tx.FromAddress, tx.ToAddress, tx.Amount)
  log.Printf("Transaction Hash : 0x%08x", CalculateTransactionHash(tx))
  log.Printf("Transaction Signature : 0x%08x", tx.Signature)
  log.Printf("Transaction Signer Public Key : 0x%08x", tx.SignerPublicKey)
  log.Printf("Transaction Nonce : 0x%08x", tx.Nonce)
  log.Printf("Transaction Fee : 0x%08x", tx.Fee)
}

func CompareTransactions(tx1 Transaction, tx2 Transaction) bool {
  //TODO: Signature and public key?
  if tx1.FromAddress == tx2.FromAddress && tx1.ToAddress == tx2.ToAddress && tx1.Amount == tx2.Amount && tx1.Fee == tx2.Fee && tx1.Nonce == tx2.Nonce {
    return true
  }
  return false
}

func verifyTransactionFromSender(tx Transaction) bool {
  // Hash the transaction data

  txAccountFromPublicKey := signer.PublicKeyToAddress(&tx.SignerPublicKey)
  if txAccountFromPublicKey != tx.FromAddress {
    log.Printf("Transaction from address does not match public key %d != %d", txAccountFromPublicKey, tx.FromAddress)
    return false
  }


  txHash := sha256.Sum256([]byte(fmt.Sprintf("%d%d%d%d%d", tx.FromAddress, tx.ToAddress, tx.Amount, tx.Fee, tx.Nonce)))

  // Verify the signature
  err := rsa.VerifyPKCS1v15(&tx.SignerPublicKey, crypt.SHA256, txHash[:], tx.Signature)
  if err != nil {
    fmt.Println("Error verifying signature:", err)
    return false
  }

  return true
}

func verifySenderHasFunds(transaction Transaction) bool {
  // Enough funds
  if transaction.Amount + transaction.Fee > TheLedger.AccountBalances[transaction.FromAddress] {
    log.Printf("Transacting %d which is more than owned %d", transaction.Amount + transaction.Fee, TheLedger.AccountBalances[transaction.FromAddress])
    log.Printf("Address %d does not have enough funds", transaction.FromAddress)
    return false
  }

  return true
}

func verifyTransactionNonce(transaction Transaction) bool {
  //TODO: THink about this more
  if transaction.Nonce != TheLedger.AccountNonces[transaction.FromAddress] + 1 {
    log.Printf("Transaction nonce %d does not match expected nonce %d", transaction.Nonce, TheLedger.AccountNonces[transaction.FromAddress] + 1)
    return false
  }

  return true
}

//TODO: Sign from a specific address and hardcode signature
func createGenesisTransaction() Transaction {
  accountNumber := "12345"
  transaction, err := CreateTransaction(0, 1, 0, 0, 0, "accounts/account-" + accountNumber + "/private_key.pem", "accounts/account-" + accountNumber + "/public_key.pem", true)
  if err != nil {
    log.Println("Error creating genesis transaction ", err)
  }

  return transaction
}

func (tx *Transaction) SignTransaction(privateKeyFilename string) {
  privateKey, err := signer.ParsePrivateKey(privateKeyFilename)
  if err != nil {
    fmt.Println("Error parsing private key:", err)
    return
  }

  txHash := sha256.Sum256([]byte(fmt.Sprintf("%d%d%d%d%d", tx.FromAddress, tx.ToAddress, tx.Amount, tx.Fee, tx.Nonce)))
  // Sign the hash with the private key
  signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypt.SHA256, txHash[:])
  if err != nil {
    fmt.Println("Error signing transaction:", err)
    return
  }

  // Save the signature to the transaction
  tx.Signature = signature
}
