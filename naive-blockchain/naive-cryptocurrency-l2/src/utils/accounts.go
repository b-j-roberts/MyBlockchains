package utils

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
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

func SetSequencer(address common.Address) {
  // Set sequencer address
  sequencerAddress = address
}

func GetSequencer() common.Address {
  // Get sequencer address
  return sequencerAddress
}
