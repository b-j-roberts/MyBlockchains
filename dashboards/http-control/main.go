package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/tidwall/gjson"
)

var (
  httpConsoleDir, _ = os.Getwd()
  indexHtmlFile = httpConsoleDir + "/index.html"
  logsHtmlFile = httpConsoleDir + "/logs.html"
  logsHtmlTemplate = httpConsoleDir + "/logs.html.template"

  logsDir = httpConsoleDir + "/logs"
  gethL1MinerLogsFile = logsDir + "/geth-l1-miner.log"
  gethL1RpcLogsFile = logsDir + "/geth-l1-rpc.log"
  gethL2SequencerLogsFile = logsDir + "/geth-l2-sequencer.log"
  smartContractExporterLogsFile = logsDir + "/smart-contract-exporter.log"
  proverLogsFile = logsDir + "/prover.log"

  l1ContractAddressFile = httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2/contracts/builds/contract-address.txt"
  l1ContractAddress = "0x0"
)

func templateReplaceLogs(logBytes []byte, templateName string, templateBytes []byte) string {
  output := strings.Replace(string(templateBytes), "{{" + templateName + "}}", string(logBytes), -1)
  return output
}

func runMakeCommandEnv(dir string, command string, env string) (string, error) {
  err := os.Chdir(dir)
  if err != nil {
    return "", err
  }

  cmd := exec.Command("make", command)
  cmd.Env = append(os.Environ(), env)
  out, err := cmd.Output()
  if err != nil {
    return string(out), err
  }

  return string(out), nil
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
      log.Println("Error logs.html from template...", logsHtmlTemplate)
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

  case "/geth-l1-launch-miner-local":
    log.Println("Launching L1 miner locally...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network/",
                                  "launch-miner-local", "OUTPUT_FILE=" + gethL1MinerLogsFile)
    if err != nil {
      log.Println("Error launching miner locally: ", out, err)
      http.Error(w, "Error launching miner locally.", 500)
      return
    }

    log.Println("L1 miner launched successfully!")
    fmt.Fprintf(w, "L1 Miner Launched Successfully!")

  case "/geth-l1-launch-rpc-local":
    log.Println("Launching L1 RPC locally...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network/",
                                  "launch-rpc-local", "OUTPUT_FILE=" + gethL1RpcLogsFile)
    if err != nil {
      log.Println("Error launching RPC locally: ", out, err)
      http.Error(w, "Error launching RPC locally.", 500)
      return
    }

    log.Println("L1 RPC launched successfully!")
    fmt.Fprintf(w, "L1 RPC Launched Successfully!")

  case "/l2-build-local":
    log.Println("Building L2 locally...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "all", "")
    if err != nil {
      log.Println("Error building L2 locally: ", out, err)
      http.Error(w, "Error building L2 locally.", 500)
      return
    }

    log.Println("L2 built successfully!")
    fmt.Fprintf(w, "L2 Built Successfully!")

  case "/geth-l2-launch-sequencer-local":
    log.Println("Launching L2 sequencer locally...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "run-sequencer", "OUTPUT_FILE=" + gethL2SequencerLogsFile)
    if err != nil {
      log.Println("Error launching sequencer locally: ", out, err)
      http.Error(w, "Error launching sequencer locally.", 500)
      return
    }

    log.Println("L2 sequencer launched successfully!")
    fmt.Fprintf(w, "L2 Sequencer Launched Successfully!")

  case "/l2-launch-prover-local":
    log.Println("Launching L2 prover locally...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "run-prover", "OUTPUT_FILE=" + proverLogsFile)

    if err != nil {
      log.Println("Error launching prover locally: ", out, err)
      http.Error(w, "Error launching prover locally.", 500)
      return
    }

    log.Println("L2 prover launched successfully!")
    fmt.Fprintf(w, "L2 Prover Launched Successfully!")

  case "/smart-contract-metrics-exporter-local":
    log.Println("Launching smart contract metrics exporter locally...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "run-smart-contract-metrics", "OUTPUT_FILE=" + smartContractExporterLogsFile)
    if err != nil {
      log.Println("Error launching smart contract metrics exporter locally: ", out, err)
      http.Error(w, "Error launching smart contract metrics exporter locally.", 500)
      return
    }

    log.Println("Smart contract metrics exporter launched successfully!")
    fmt.Fprintf(w, "Smart Contract Metrics Exporter Launched Successfully!")

  case "/geth-l1-docker-build":
    log.Println("Building Geth L1 docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network/",
                                  "docker-build", "")
    if err != nil {
      log.Println("Error building Geth L1 docker image: ", out, err)
      http.Error(w, "Error building Geth L1 docker image.", 500)
      return
    }

    log.Println("Geth L1 docker image built successfully!")
    fmt.Fprintf(w, "Geth L1 Docker Image Built Successfully!")

  case "/geth-l1-docker-push":
    log.Println("Pushing Geth L1 docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network/",
                                  "docker-push", "")
    if err != nil {
      log.Println("Error pushing Geth L1 docker image: ", out, err)
      http.Error(w, "Error pushing Geth L1 docker image.", 500)
      return
    }

    log.Println("Geth L1 docker image pushed successfully!")
    fmt.Fprintf(w, "Geth L1 Docker Image Pushed Successfully!")

  case "/geth-l1-miner-docker-run":
    log.Println("Running Geth L1 docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network/",
                                  "docker-run-miner", "OUTPUT_FILE=" + gethL1MinerLogsFile)
    if err != nil {
      log.Println("Error running Geth L1 docker image: ", out, err)
      http.Error(w, "Error running Geth L1 docker image.", 500)
      return
    }

    log.Println("Geth L1 docker image ran successfully!")
    fmt.Fprintf(w, "Geth L1 Docker Image Ran Successfully!")

  case "/geth-l1-rpc-docker-run":
    log.Println("Running Geth L1 RPC docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network/",
                                  "docker-run-rpc", "OUTPUT_FILE=" + gethL1RpcLogsFile)
    if err != nil {
     log.Println("Error running Geth L1 RPC docker image: ", out, err)
      http.Error(w, "Error running Geth L1 RPC docker image.", 500)
      return
    }

    log.Println("Geth L1 RPC docker image ran successfully!")
    fmt.Fprintf(w, "Geth L1 RPC Docker Image Ran Successfully!")

  case "/l2-docker-build":
    log.Println("Building L2 docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "docker-build", "")
    if err != nil {
      log.Println("Error building L2 docker image: ", out, err)
      http.Error(w, "Error building L2 docker image.", 500)
      return
    }

    log.Println("L2 docker image built successfully!")
    fmt.Fprintf(w, "L2 Docker Image Built Successfully!")

  case "/l2-docker-push":
    log.Println("Pushing L2 docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "docker-push", "")
    if err != nil {
      log.Println("Error pushing L2 docker image: ", out, err)
      http.Error(w, "Error pushing L2 docker image.", 500)
      return
    }

    log.Println("L2 docker image pushed successfully!")
    fmt.Fprintf(w, "L2 Docker Image Pushed Successfully!")

  case "/geth-l2-sequencer-docker-run":
    l1ContractAddress := r.URL.Query().Get("l1ContractAddress")
    log.Println("Running Geth L2 sequencer docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "docker-run-sequencer", "OUTPUT_FILE=" + gethL2SequencerLogsFile + " L1_CONTRACT_ADDRESS=" + l1ContractAddress)
    if err != nil {
      log.Println("Error running Geth L2 sequencer docker image: ", out, err)
      http.Error(w, "Error running Geth L2 sequencer docker image.", 500)
      return
    }

    log.Println("Geth L2 sequencer docker image ran successfully!")
    fmt.Fprintf(w, "Geth L2 Sequencer Docker Image Ran Successfully!")

  case "/l2-prover-docker-run":
    l1ContractAddress := r.URL.Query().Get("l1ContractAddress")
    log.Println("Running Geth L2 prover docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "docker-run-prover", "OUTPUT_FILE=" + proverLogsFile + " L1_CONTRACT_ADDRESS=" + l1ContractAddress)
    if err != nil {
      log.Println("Error running Geth L2 prover docker image: ", out, err)
      http.Error(w, "Error running Geth L2 prover docker image.", 500)
      return
    }

    log.Println("Geth L2 prover docker image ran successfully!")
    fmt.Fprintf(w, "Geth L2 Prover Docker Image Ran Successfully!")

  case "/smart-contract-metrics-exporter-docker-run":
    l1ContractAddress := r.URL.Query().Get("l1ContractAddress")
    log.Println("Running smart contract metrics exporter docker image...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "docker-run-metrics-exporter", "OUTPUT_FILE=" + smartContractExporterLogsFile + " L1_CONTRACT_ADDRESS=" + l1ContractAddress)
    if err != nil {
      log.Println("Error running smart contract metrics exporter docker image: ", out, err)
      http.Error(w, "Error running smart contract metrics exporter docker image.", 500)
      return
    }

    log.Println("Smart contract metrics exporter docker image ran successfully!")
    fmt.Fprintf(w, "Smart Contract Metrics Exporter Docker Image Ran Successfully!")

  case "/deploy-l1-contract":
    log.Println("Deploying L1 contract...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "deploy-private-l1", "")
    if err != nil {
      log.Println("Error deploying L1 contract: ", out, err)
      http.Error(w, "Error deploying L1 contract.", 500)
      return
    }

    // Load l1 contract address from output file ( json file w/ .address field )
    l1ContractAddressData, err := ioutil.ReadFile(l1ContractAddressFile)
    if err != nil {
      log.Println("Error reading L1 contract address file: ", err)
      http.Error(w, "Error reading L1 contract address file.", 500)
      return
    }
    l1ContractAddress = gjson.Get(string(l1ContractAddressData), "address").String()

    log.Println("L1 contract address: ", l1ContractAddress)
    log.Println("L1 contract deployed successfully!")
    fmt.Fprintf(w, "L1 Contract Deployed Successfully!")

  case "/connect-geth-l1-nodes-local":
    log.Println("Connecting Geth L1 Miner & RPC nodes...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network",
                                  "connect-local", "")

    if err != nil {
      log.Println("Error connecting Geth L1 Miner & RPC nodes: ", out, err)
      http.Error(w, "Error connecting Geth L1 Miner & RPC nodes.", 500)
      return
    }

    log.Println("Geth L1 Miner & RPC nodes connected successfully!")
    fmt.Fprintf(w, "Geth L1 Miner & RPC Nodes Connected Successfully!")

  case "/connect-geth-l1-nodes-docker":
    log.Println("Connecting Geth L1 Miner & RPC nodes in docker...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network",
                                  "connect-docker", "")
    if err != nil {
      log.Println("Error connecting Geth L1 Miner & RPC nodes in docker: ", out, err)
      http.Error(w, "Error connecting Geth L1 Miner & RPC nodes in docker.", 500)
      return
    }

    log.Println("Geth L1 Miner & RPC nodes connected in docker successfully!")
    fmt.Fprintf(w, "Geth L1 Miner & RPC Nodes Connected in Docker Successfully!")

  case "/clean-all":
    log.Println("Cleaning all...")
    out, err := runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                  "clean", "")
    if err != nil {
      log.Println("Error cleaning all (clean): ", out, err)
      //http.Error(w, "Error cleaning all.", 500)
      //return
    }

    out, err = runMakeCommandEnv(httpConsoleDir + "/../../naive-blockchain/naive-cryptocurrency-l2",
                                 "kube-clean-all", "")
    if err != nil {
      log.Println("Error cleaning all (kube-clean-all): ", out, err)
      //http.Error(w, "Error cleaning all.", 500)
      //return
    }

    out, err = runMakeCommandEnv(httpConsoleDir + "/../../eth-private-network",
                                 "clean", "")
    if err != nil {
      log.Println("Error cleaning all (clean 2): ", out, err)
      //http.Error(w, "Error cleaning all.", 500)
      //return
    }

    // Clean logs folder
    os.RemoveAll(logsDir)

    if _, err := os.Stat(logsDir); os.IsNotExist(err) {
      os.Mkdir(logsDir, 0755)
    }

    log.Println("Cleaned all done!")
    fmt.Fprintf(w, "Cleaned All Done!")

    //TODO: Shutdown command?, one click all, send transaction, store l1 contract address, ...
  default:
    fmt.Fprintf(w, "Sorry, we couldn't find that page.")
  }
}

func main() {
  // Create logs directory if it doesn't exist
  if _, err := os.Stat(logsDir); os.IsNotExist(err) {
    os.Mkdir(logsDir, 0755)
  }

  http.HandleFunc("/", httpHandler)

  fmt.Println("Http Controller is listening on port 7000...")
  http.ListenAndServe(":7000", nil)
}
