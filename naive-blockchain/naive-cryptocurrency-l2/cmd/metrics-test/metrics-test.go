package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
  // metrics test value map types
  TEST_EQUAL = iota
  TEST_GREATER
  TEST_LESS
)

func printTestCompare(op int) string {
  switch op {
  case TEST_EQUAL:
    return "=="
  case TEST_GREATER:
    return ">"
  case TEST_LESS:
    return "<"
  default:
    return "??"
  }
}

type MetricsTest struct {
  MetricName string
  TestType   int
  TestValue  float64
}

//TODO: Add a "Metric Compare" type that can compare two metrics values & use it on tests marked with MTC
func (metricsTest *MetricsTest) Test() (bool, error) {
  value, err := getMetricValue(metricsTest.MetricName)
  if err != nil {
    return false, err
  }

  switch metricsTest.TestType {
  case TEST_EQUAL:
    return value == metricsTest.TestValue, nil
  case TEST_GREATER:
    return value > metricsTest.TestValue, nil
  case TEST_LESS:
    return value < metricsTest.TestValue, nil
  }

  return false, nil
}

func (metricsTest *MetricsTest) PrintTest() string {
  return metricsTest.MetricName + "  " + printTestCompare(metricsTest.TestType) + "  " + strconv.FormatFloat(metricsTest.TestValue, 'f', -1, 64)
}

var metricsTests = []MetricsTest{
  {"batch_count{job='l2-smart-contract-exporter'}", TEST_GREATER, 3},
  {"last_confirmed_batch{job='l2-smart-contract-exporter'}", TEST_GREATER, 3}, //MTC

  {"bridge_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 4980000000000000},
  {"l2_burn_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 20000000000000},
  {"deposit_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 5},
  {"withdrawal_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 2},
  {"l2_deposit_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 5}, //MTC
  {"l2_withdrawal_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 2}, //MTC

  {"l1_token_deposit_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 14},
  {"l1_token_withdrawal_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 7},
  {"l2_token_deposit_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 14}, //MTC
  {"l2_token_withdrawal_nonce{job='l2-smart-contract-exporter'}", TEST_EQUAL, 7}, //MTC

  {"l1_basic_token_bridge_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 497000},
  {"l1_basic_token_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 99999999503000},
  {"l2_basic_token_supply{job='l2-smart-contract-exporter'}", TEST_EQUAL, 497000}, //MTC
  {"l2_basic_token_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 497000}, //MTC
  {"l1_stable_token_bridge_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 497000},
  {"l1_stable_token_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 99999999503000},
  {"l2_stable_token_supply{job='l2-smart-contract-exporter'}", TEST_EQUAL, 497000}, //MTC
  {"l2_stable_token_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 497000}, //MTC

  {"l1_basic_nft_bridge_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 2},
  {"l1_basic_nft_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 3},
  {"l2_basic_nft_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 2}, //MTC
  {"l1_special_nft_bridge_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 1},
  {"l1_special_nft_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 0},
  {"l2_special_nft_sequencer_balance{job='l2-smart-contract-exporter'}", TEST_EQUAL, 1}, //MTC

  {"prover_total_proofs{job='l2-prover'}", TEST_GREATER, 2}, //MTC
  {"prover_total_proofs_verified{job='l2-prover'}", TEST_GREATER, 2}, //MTC

  {"chain_head_block{job='geth-l1-miner'}", TEST_GREATER, 10},
  {"txpool_invalid{job='geth-l1-miner'}", TEST_EQUAL, 0},
  {"txpool_pending{job='geth-l1-miner'}", TEST_EQUAL, 0},
  {"txpool_queued{job='geth-l1-miner'}", TEST_EQUAL, 0},
  {"txpool_valid{job='geth-l1-miner'}", TEST_EQUAL, 62},

  {"chain_head_block{job='geth-l2-sequencer'}", TEST_GREATER, 10},
  {"txpool_invalid{job='geth-l2-sequencer'}", TEST_EQUAL, 0},
  {"txpool_pending{job='geth-l2-sequencer'}", TEST_EQUAL, 0},
  {"txpool_queued{job='geth-l2-sequencer'}", TEST_EQUAL, 0},
  {"txpool_valid{job='geth-l2-sequencer'}", TEST_EQUAL, 60},
}

func getMetricValue(metricName string) (float64, error) {
  registry := prometheus.NewRegistry()
  registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
  registry.MustRegister(prometheus.NewGoCollector())

  response, err := http.Get("http://localhost:9090/api/v1/query?query=" + metricName)
  if err != nil {
    return 0, err
  }
  defer response.Body.Close()

  result := struct {
    Data struct {
      Result []struct {
        Value []interface{} `json:"value"`
      } `json:"result"`
    } `json:"data"`
  }{}

  if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
    return 0, err
  }

  if len(result.Data.Result) == 0 {
    return 0, nil
  }

  value, _ := strconv.ParseFloat(result.Data.Result[0].Value[1].(string), 64)
  return value, nil
}

func main() {
  log.Println("Starting tests...")

  for _, metricsTest := range metricsTests {
    result, err := metricsTest.Test()
    if err != nil {
      log.Printf("Test failed: %s -- Error : %s", metricsTest.MetricName, err)
      continue
    }

    if !result {
      log.Printf("Test failed: %s -- No Result", metricsTest.MetricName)
      continue
    }

    log.Printf("Test passed: %s", metricsTest.MetricName + "  " + printTestCompare(metricsTest.TestType) + "  " + strconv.FormatFloat(metricsTest.TestValue, 'f', -1, 64))
  }

  log.Println("All tests done!")
}
