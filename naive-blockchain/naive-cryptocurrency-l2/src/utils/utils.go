package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/metrics"
)

func StartSystemMetrics() {
  go metrics.CollectProcessMetrics(3 * time.Second)
}

func readJsonValue(jsonStr string, key string) string {
  jsonMap := make(map[string]interface{})
  json.Unmarshal([]byte(jsonStr), &jsonMap)
  return jsonMap[key].(string)
}

func ReadContractAddressFromFile(path string) (common.Address, error) {
  file, err := os.Open(path)
  if err != nil {
    return common.Address{}, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  scanner.Scan()
  addressJson := scanner.Text()

  // Parse address out of json under key "address"
  address := readJsonValue(addressJson, "address")
  return common.HexToAddress(address), nil
}
