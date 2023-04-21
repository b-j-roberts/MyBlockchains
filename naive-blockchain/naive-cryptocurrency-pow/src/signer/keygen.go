package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func VerifyKeyPairExists(accountDirectory string, accountNumber int) bool {
  // check if account directory exists
	accountDirPath := filepath.Join(accountDirectory, fmt.Sprintf("account-%d", accountNumber))
	if _, err := os.Stat(accountDirPath); os.IsNotExist(err) {
		return false
	}

  privateKeyPath := filepath.Join(accountDirPath, "private_key.pem")
	privateKey, err := ParsePrivateKey(privateKeyPath)
	if err != nil {
		return false
	}

	// read public key file
	publicKeyPath := filepath.Join(accountDirPath, "public_key.pem")
	publicKey, err := ParsePublicKey(publicKeyPath)
	if err != nil {
		return false
	}

  // verify address file
	addressPath := filepath.Join(accountDirPath, "account_address.txt")
	addressBytes, err := ioutil.ReadFile(addressPath)
	if err != nil {
		return false
	}

  if string(addressBytes) != fmt.Sprintf("%X", PublicKeyToAddress(publicKey)) {
		return false
	}

	// verify that public key matches private key
	if privateKey.PublicKey.N.Cmp(publicKey.N) != 0 {
		return false
	}

  return true
}

func ParsePrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
  privateKeyFile, err := os.Open(privateKeyPath)
  if err != nil {
    fmt.Println("Error opening private key file:", err)
		return nil, err
  }
  defer privateKeyFile.Close()

  privateKeyPEM, err := ioutil.ReadAll(privateKeyFile)
  if err != nil {
    fmt.Println("Error reading private key file:", err)
		return nil, err
  }

  privateKeyBlock, _ := pem.Decode(privateKeyPEM)
  if privateKeyBlock == nil {
    fmt.Println("Error decoding private key PEM")
		return nil, err
  }

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
	//	privateKey, err = x509.ParsePKCS8PrivateKey(privateKeyBytes)
		return nil, err
	}
	return privateKey, nil
}

func ParsePublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
  publicKeyFile, err := os.Open(publicKeyPath)
  if err != nil {
    fmt.Println("Error opening public key file:", err)
    return nil, err
  }
  defer publicKeyFile.Close()

  publicKeyPEM, err := ioutil.ReadAll(publicKeyFile)
  if err != nil {
    fmt.Println("Error reading private key file:", err)
    return nil, err
  }

  publicKeyBlock, _ := pem.Decode(publicKeyPEM)
  if publicKeyBlock == nil {
    fmt.Println("Error decoding public key PEM")
    return nil, err
  }

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		publicKey, err = ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return nil, err
		}
	}
	return publicKey, nil
}

func ParsePKIXPublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	ifc, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return publicKey, nil
}

func privateKeyToFile(pk *rsa.PrivateKey, filename string) error {
  // Save the private key to a PEM file           
  privateKeyFile, err := os.Create(filename)
  if err != nil {                                 
    return err
  }
  defer privateKeyFile.Close()

  privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)}
  if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
    return err
  }

  return nil
}

func publicKeyToFile(pk *rsa.PublicKey, filename string) error {
  // Save the public key to a PEM file
  publicKeyFile, err := os.Create(filename)
  if err != nil {
    return err
  }
  defer publicKeyFile.Close()

  publicKeyPEM, err := x509.MarshalPKIXPublicKey(pk)
  if err != nil {
    return err
  }
  if err := pem.Encode(publicKeyFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyPEM}); err != nil {
    return err
  }

  return nil
}

//TODO: Test suite
//TODO: Seed phrase
func GeneratePublicPrivateKey(accountDirectory string, accountNumber int) {
  // Create the directory if it doesn't exist
  err := os.MkdirAll(accountDirectory, os.ModePerm)
  if err != nil {
      fmt.Println(err)
  }

  accountDirPath := filepath.Join(accountDirectory, fmt.Sprintf("account-%d", accountNumber))

  err = os.MkdirAll(accountDirPath, os.ModePerm)
  if err != nil {
      fmt.Println(err)
  }

  // Generate a new RSA key pair
  privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    fmt.Println(err)
    return
  }

  // Get the public key in the right format
  publicKey := &privateKey.PublicKey

  err = privateKeyToFile(privateKey,  filepath.Join(accountDirPath, "private_key.pem"))
  if err != nil {
    fmt.Println(err)
    return
  }

  err = publicKeyToFile(publicKey,  filepath.Join(accountDirPath, "public_key.pem"))
  if err != nil {
    fmt.Println(err)
    return
  }

  accountAddressFile, err := os.Create(filepath.Join(accountDirPath, "account_address.txt"))
  if err != nil {
    fmt.Println(err)
    return
  }
  defer accountAddressFile.Close()

  fmt.Fprintf(accountAddressFile, "%X", PublicKeyToAddress(publicKey))

  //TODO: Print filename from variables
  fmt.Printf("RSA Key pair generated for account number %d\n", accountNumber)
  fmt.Printf("Private key saved to %s/account-%d/private_key.pem\n", accountDirectory, accountNumber)
  fmt.Printf("Public key saved to %s/account-%d/public_key.pem\n", accountDirectory, accountNumber)
  fmt.Printf("Linked account address: %x\n",  PublicKeyToAddress(publicKey))
}
