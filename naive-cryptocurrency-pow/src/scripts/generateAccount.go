package main

import (
  "flag"
  "fmt"
  "os"

  golog "github.com/ipfs/go-log"

  "naive-cryptocurrency-pow/src/signer"
)

func main() {
  golog.SetAllLoggers(golog.LevelInfo) // Change to DEBUG for extra info

  accountsDirectory := flag.String("accounts", "", "Directory containing account dirs")
  accountNumber := flag.Int("account-id", 12345, "Account number to use from account directory.")
  overwrite := flag.Bool("overwrite", false, "Overwrite existing account key pair (dangerous)")
  flag.Parse()

  if *accountsDirectory == "" {
		fmt.Fprintln(os.Stderr, "Error: missing required flag --accounts")
		os.Exit(1)
	}

  if *overwrite == false && signer.VerifyKeyPairExists(*accountsDirectory, *accountNumber) {
    fmt.Printf("Account key pair already exists\n")
  } else {
    signer.GeneratePublicPrivateKey(*accountsDirectory, *accountNumber)
  }
}

