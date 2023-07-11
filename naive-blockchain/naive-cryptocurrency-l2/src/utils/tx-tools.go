package utils

import (
	"bytes"
	"compress/flate"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func CreateTransactOpts(account accounts.Account, chainID *big.Int) (*bind.TransactOpts, error) {
  keystore := keystore.NewKeyStore(addressKeyStoreDirMap[account.Address], keystore.StandardScryptN, keystore.StandardScryptP)
  keystore.Unlock(account, "password") //TODO: Hardcoded password
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
