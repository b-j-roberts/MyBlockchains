package utils

import (
	"bytes"
	"compress/flate"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
)

func CreateTransactOpts(account accounts.Account, chainID *big.Int) (*bind.TransactOpts, error) {
  // Create a temporary keystore
  osHomeDir := os.Getenv("HOME")
  // Create transactor directory if it doesn't exist
  if _, err := os.Stat(osHomeDir + "/.transactor"); os.IsNotExist(err) {
    os.Mkdir(osHomeDir + "/.transactor", 0700)
  }

  keystore := keystore.NewKeyStore(osHomeDir + "/.transactor", keystore.StandardScryptN, keystore.StandardScryptP)
  // Read password from environment variable
  keystore.Unlock(account, os.Getenv("ACCOUNT_PASS"))
  return bind.NewKeyStoreTransactorWithChainID(keystore, account, chainID)
}

func CompressTransactionData(data []byte) ([]byte, error) {
  //TODO: Use better compression algorithm

  // Compress data using compress/flate
  var buf bytes.Buffer
  compressor, err := flate.NewWriter(&buf, flate.BestCompression)
  if err != nil {
    return nil, err
  }

  if _, err := compressor.Write(data); err != nil {
    return nil, err
  }

  if err := compressor.Close(); err != nil {
    return nil, err
  }

  return buf.Bytes(), nil
}

func ReceiptLogsWithEvent(receipt *types.Receipt, eventSignature []byte) []*types.Log {
  var receipt_logs []*types.Log
  for _, receipt_log := range receipt.Logs {
    if bytes.Equal(receipt_log.Topics[0].Bytes(), eventSignature) {
      receipt_logs = append(receipt_logs, receipt_log)
    }
  }
  return receipt_logs
}
