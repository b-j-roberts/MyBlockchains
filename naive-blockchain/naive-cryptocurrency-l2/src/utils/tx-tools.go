package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func MakeTransactOpts(from common.Address, value uint64) bind.TransactOpts {
  return bind.TransactOpts{
    From: from,
    Value: big.NewInt(int64(value)),
    GasPrice: big.NewInt(200), //TODO: Hardcoded
    GasLimit: uint64(3000000), //TODO: Hardcoded
    Signer: KeystoreSignTx,
  }
}
