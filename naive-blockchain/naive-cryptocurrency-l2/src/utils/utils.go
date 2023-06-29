package utils

import (
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var sequencerAddress common.Address

func AddressFromFile(path string) (common.Address, error) {
  // Read eth address from json file under .address field
  var jsonObject map[string]interface{}

  file, err := os.Open(path)
  if err != nil {
    return common.Address{}, err
  }
  defer file.Close()

  jsonParser := json.NewDecoder(file)
  if err = jsonParser.Decode(&jsonObject); err != nil {
    return common.Address{}, err
  }
  addressString := jsonObject["address"].(string)
  return common.HexToAddress(addressString), nil
}

// Store address -> KeyStoreDir map
var addressKeyStoreDirMap = make(map[common.Address]string)

func StoreKeyStoreDir(address common.Address, dir string) {
  addressKeyStoreDirMap[address] = dir
}

func SetSequencer(address common.Address) {
  // Set sequencer address
  sequencerAddress = address
}

func GetSequencer() common.Address {
  // Get sequencer address
  return sequencerAddress
}

func KeystoreSignTx(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
  // Sign transaction using keystore file
  // keystore := keystore.NewKeyStore("/home/b-j-roberts/workspace/blockchain/my-chains/eth-private-network/data/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
  // Create default keystore
  keystore := keystore.NewKeyStore(addressKeyStoreDirMap[address], keystore.StandardScryptN, keystore.StandardScryptP)
  return keystore.SignTxWithPassphrase(accounts.Account{Address: address}, "password", tx, big.NewInt(505)) //TODO: Hardcode
}
