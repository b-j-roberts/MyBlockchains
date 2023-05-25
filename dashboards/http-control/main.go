package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
  indexHtmlFile = "index.html"
  logsHtmlFile = "logs.html"
  logsHtmlTemplate = "logs.html.template"

  logsDir = "logs"
  gethL1MinerLogsFile = logsDir + "/geth-l1-miner.log"
  gethL1RpcLogsFile = logsDir + "/geth-l1-rpc.log"
  gethL2SequencerLogsFile = logsDir + "/geth-l2-sequencer.log"
  smartContractExporterLogsFile = logsDir + "/smart-contract-exporter.log"
  proverLogsFile = logsDir + "/prover.log"
)

func templateReplaceLogs(logBytes []byte, templateName string, templateBytes []byte) string {
  output := strings.Replace(string(templateBytes), "{{" + templateName + "}}", string(logBytes), -1)
  return output
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Request received for path: ", r.URL.Path)
  switch r.URL.Path {
  case "/":
    tmpl := template.Must(template.ParseFiles(indexHtmlFile))
    tmpl.Execute(w, nil)

  case "/logs":
    //TODO: Read page & generate logs.html
    templateBytes, err := ioutil.ReadFile(logsHtmlTemplate)
    if err != nil {
      http.Error(w, "Error reading logs.html.template.", 500)
      return
    }

    gethL1LogsBytes, err := ioutil.ReadFile(gethL1MinerLogsFile)
    if err != nil {
      gethL1LogsBytes = []byte("No Geth L1 miner logs found.")
    }
    gethL1LogsBytes = []byte(strings.Replace(string(gethL1LogsBytes), "\n", "<br>", -1))
    output := templateReplaceLogs(gethL1LogsBytes, "GETH_L1_MINER_LOGS", templateBytes)

    gethL1RpcLogsBytes, err := ioutil.ReadFile(gethL1RpcLogsFile)
    if err != nil {
      gethL1RpcLogsBytes = []byte("No Geth L1 RPC logs found.")
    }
    gethL1RpcLogsBytes = []byte(strings.Replace(string(gethL1RpcLogsBytes), "\n", "<br>", -1))
    output = templateReplaceLogs(gethL1RpcLogsBytes, "GETH_L1_RPC_LOGS", []byte(output))

    gethL2SequencerLogsBytes, err := ioutil.ReadFile(gethL2SequencerLogsFile)
    if err != nil {
      gethL2SequencerLogsBytes = []byte("No Geth L2 sequencer logs found.")
    }
    gethL2SequencerLogsBytes = []byte(strings.Replace(string(gethL2SequencerLogsBytes), "\n", "<br>", -1))
    output = templateReplaceLogs(gethL2SequencerLogsBytes, "GETH_L2_SEQUENCER_LOGS", []byte(output))

    smartContractExporterLogsBytes, err := ioutil.ReadFile(smartContractExporterLogsFile)
    if err != nil {
      smartContractExporterLogsBytes = []byte("No smart contract exporter logs found.")
    }
    smartContractExporterLogsBytes = []byte(strings.Replace(string(smartContractExporterLogsBytes), "\n", "<br>", -1))
    output = templateReplaceLogs(smartContractExporterLogsBytes, "SC_METRICS_EXPORTER_LOGS", []byte(output))

    proverLogsBytes, err := ioutil.ReadFile(proverLogsFile)
    if err != nil {
      proverLogsBytes = []byte("No prover logs found.")
    }
    proverLogsBytes = []byte(strings.Replace(string(proverLogsBytes), "\n", "<br>", -1))
    output = templateReplaceLogs(proverLogsBytes, "PROVER_LOGS", []byte(output))

    err = ioutil.WriteFile(logsHtmlFile, []byte(output), 0644)
    if err != nil {
      http.Error(w, "Error writing logs.html.", 500)
      return
    }

    tmpl := template.Must(template.ParseFiles(logsHtmlFile))
    tmpl.Execute(w, nil)

  case "/exec-a":
    cmd := exec.Command("touch", "hello-a")
    err := cmd.Run()
    if err != nil {
      http.Error(w, "Error executing command A.", 500)
      return
    }

    fmt.Fprintf(w, "Command A executed successfully!")

  case "/exec-b":
    cmd := exec.Command("touch", "hello-b")
    err := cmd.Run()
    if err != nil {
      http.Error(w, "Error executing command B.", 500)
      return
    }

    fmt.Fprintf(w, "Command B executed successfully!")

  default:
    fmt.Fprintf(w, "Sorry, we couldn't find that page.")
  }
}

func main() {
  http.HandleFunc("/", httpHandler)

  fmt.Println("Http Controller is listening on port 7000...")
  http.ListenAndServe(":7000", nil)
}
