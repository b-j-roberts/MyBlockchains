package signer

import (
  "crypto/rsa"
  "crypto/sha256"
  "encoding/hex"
  "fmt"
  "math/big"

  "github.com/btcsuite/btcutil/base58"
)

// PublicKeyToAddress converts the public key to an address using Base58Check encoding
func PublicKeyToAddress(pk *rsa.PublicKey) uint64 {
  //// Hash the public key using SHA256
  ////hexString := hex.EncodeToString(pk.N.Bytes())
  ////hash := sha256.Sum256([]byte(hexString))

  //// Add a checksum to the end of the hash
  ////TODO: Why checksum?
  ////checksum := hash[:4] //TODO: Is this a checksum?
  ////encoded := append(hash[:], checksum...)

  //// Encode the result using Base58Check
  //a := base58.Encode(pk.N.Bytes())

  //// Encode the address as a hexadecimal string
  ////hexString := hex.EncodeToString([]byte(a))

  //// Parse the hexadecimal string to a uint64
  //address, err := strconv.ParseUint(a, 10, 64)
  //if err != nil {
  //  return 0, err
  //}

  //return address, nil
  // Serialize the RSA public key to a string
	publicKeyBytes := []byte(pk.N.String())

	// Calculate the SHA256 hash of the public key
	hash := sha256.Sum256(publicKeyBytes)

	// Convert the hash to a big integer
	hashInt := new(big.Int).SetBytes(hash[:])

	// Truncate the hash to 64 bits
	truncatedHashInt := new(big.Int).And(hashInt, new(big.Int).SetUint64(0xffffffffffffffff))

	// Convert the truncated hash to a uint64
	accountAddress := truncatedHashInt.Uint64()

	return accountAddress
}

func VerifyPublicKeyLinkedToAddress(accountAddress uint64, publicKey *rsa.PublicKey) bool {
  // Serialize the RSA public key to a string
	publicKeyBytes := []byte(publicKey.N.String())

	// Calculate the SHA256 hash of the public key
	hash := sha256.Sum256(publicKeyBytes)

	// Convert the hash to a big integer
	hashInt := new(big.Int).SetBytes(hash[:])

	// Truncate the hash to 64 bits
	truncatedHashInt := new(big.Int).And(hashInt, new(big.Int).SetUint64(0xffffffffffffffff))

	// Compare the truncated hash to the input account address
	return truncatedHashInt.Uint64() == accountAddress
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
