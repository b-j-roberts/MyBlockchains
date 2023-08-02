package utils

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	basicerc20 "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/basicerc20"
	l1bridge "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l1bridge"
	l1tokenbridge "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/l1tokenbridge"
	txstorage "github.com/b-j-roberts/MyBlockchains/naive-blockchain/naive-cryptocurrency-l2/contracts/go/txstorage"
)

type L1TransactionConfig struct {
  GasLimit uint64
  GasPrice *big.Int
}

type L1ContractAddressConfig struct {
  TxStorageContractAddress common.Address
  BridgeContractAddress common.Address
  TokenBridgeContractAddress common.Address
}

type L1Contracts struct {
  TxStorageContract *txstorage.Txstorage
  BridgeContract *l1bridge.L1bridge
  TokenBridgeContract *l1tokenbridge.L1tokenbridge
}

func CreateL1ContractAddressConfig(contractAddressDir string) L1ContractAddressConfig {
  // Read contract addresses from files in directory
  txStorageContractAddress, err := ReadContractAddressFromFile(contractAddressDir + "/tx-storage-address.txt")
  if err != nil {
    log.Fatal(err)
  }

  l1BridgeContractAddress, err := ReadContractAddressFromFile(contractAddressDir + "/l1-bridge-address.txt")
  if err != nil {
    log.Fatal(err)
  }

  l1TokenBridgeContractAddress, err := ReadContractAddressFromFile(contractAddressDir + "/l1-token-bridge-address.txt")
  if err != nil {
    log.Fatal(err)
  }

  return L1ContractAddressConfig{
    TxStorageContractAddress: txStorageContractAddress,
    BridgeContractAddress: l1BridgeContractAddress,
    TokenBridgeContractAddress: l1TokenBridgeContractAddress,
  }
}

func CreateL1Contracts(client *ethclient.Client, l1ContractAddressConfig L1ContractAddressConfig) L1Contracts {
  txStorageContract, err := txstorage.NewTxstorage(l1ContractAddressConfig.TxStorageContractAddress, client)
  if err != nil {
    log.Fatal(err)
  }

  l1BridgeContract, err := l1bridge.NewL1bridge(l1ContractAddressConfig.BridgeContractAddress, client)
  if err != nil {
    log.Fatal(err)
  }

  l1TokenBridgeContract, err := l1tokenbridge.NewL1tokenbridge(l1ContractAddressConfig.TokenBridgeContractAddress, client)
  if err != nil {
    log.Fatal(err)
  }

  return L1Contracts{
    TxStorageContract: txStorageContract,
    BridgeContract: l1BridgeContract,
    TokenBridgeContract: l1TokenBridgeContract,
  }
}

type L1Comms struct {
  RpcUrl string
  L1Client *ethclient.Client
  L1ChainID *big.Int
  L1TransactionConfig L1TransactionConfig

  L1ContractAddressConfig L1ContractAddressConfig
  L1Contracts L1Contracts
}

func (l1Comms *L1Comms) CreateL1TransactionOpts(fromAddress common.Address, value *big.Int) (*bind.TransactOpts, error) {
  transactOpts, err := CreateTransactOpts(accounts.Account{Address: fromAddress}, l1Comms.L1ChainID)
  if err != nil {
    return nil, err
  }
  transactOpts.GasLimit = l1Comms.L1TransactionConfig.GasLimit
  transactOpts.GasPrice = l1Comms.L1TransactionConfig.GasPrice
  transactOpts.Value = value

  return transactOpts, nil
}

func NewL1Comms(rpcUrl string, contractAddressDir string, chainID *big.Int, l1TransactionConfig L1TransactionConfig) (*L1Comms, error) {
  rawRpc, err := rpc.Dial(rpcUrl)
  if err != nil {
    return nil, err
  }

  l1Comms := &L1Comms{
    RpcUrl: rpcUrl,
    L1Client: ethclient.NewClient(rawRpc),
    L1ChainID: chainID,
    L1TransactionConfig: l1TransactionConfig,
    L1ContractAddressConfig: CreateL1ContractAddressConfig(contractAddressDir),
  }

  // Create L1 contracts
  l1Comms.L1Contracts = CreateL1Contracts(l1Comms.L1Client, l1Comms.L1ContractAddressConfig)

  return l1Comms, nil
}

func (l1Comms *L1Comms) L2GenesisOnL1(genesis *core.Genesis, posterAddress common.Address) error {
  var genesisHash [32]byte
  copy(genesisHash[:], genesis.ToBlock().Hash().Bytes())

  transactOpts, err := l1Comms.CreateL1TransactionOpts(posterAddress, big.NewInt(0))
  if err != nil {
    log.Fatalf("Failed to create L1 transaction opts: %v", err)
    return err
  }

  _, err = l1Comms.L1Contracts.TxStorageContract.StoreGenesisState(transactOpts, genesisHash)
  if err != nil {
    log.Fatalf("Failed to store genesis state on L1: %v", err)
    return err
  }

  log.Println("Stored genesis state on L1")
  return nil
}

func (l1Comms *L1Comms) PostBatch(transactionByteData []byte, id int64, hash [32]byte, posterAddress common.Address) error {
  log.Println("Posting batch to L1...")

  transactOpts, err := l1Comms.CreateL1TransactionOpts(posterAddress, big.NewInt(0))
  if err != nil {
    log.Fatalf("Failed to create L1 transaction opts: %v", err)
    return err
  }

  _, err = l1Comms.L1Contracts.TxStorageContract.StoreBatch(transactOpts, big.NewInt(id), hash, transactionByteData)
  if err != nil {
    log.Fatalf("Failed to post batch to L1: %v", err)
    return err
  }

  log.Println("Posted batch to L1")
  return nil
}

func (l1Comms *L1Comms) SubmitProof(proof []byte, batchNumber int, proverAddress common.Address) error {
  log.Println("Submitting proof to L1...")

  transactOpts, err := l1Comms.CreateL1TransactionOpts(proverAddress, big.NewInt(0))
  if err != nil {
    log.Fatalf("Failed to create L1 transaction opts: %v", err)
    return err
  }

  _, err = l1Comms.L1Contracts.TxStorageContract.SubmitProof(transactOpts, big.NewInt(int64(batchNumber)), proof)
  if err != nil {
    log.Fatalf("Failed to submit proof to L1: %v", err)
    return err
  }

  log.Println("Submitted proof to L1")
  return nil
}

func (l1BridgeComms *L1Comms) BridgeEthToL2(address common.Address, amount uint64) error {
  log.Println("Bridging ", amount, " ETH to L2 for address ", address.Hex())

  transactOpts, err := l1BridgeComms.CreateL1TransactionOpts(address, big.NewInt(int64(amount)))
  if err != nil {
    log.Println("Failed to create L1 transaction opts", err)
    return err
  }

  bridgeTx, err := l1BridgeComms.L1Contracts.BridgeContract.DepositEth(transactOpts)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}

func (l1BridgeComms *L1Comms) BridgeEthToL1(address common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " ETH to L1 for address ", address.Hex())

  transactOpts, err := l1BridgeComms.CreateL1TransactionOpts(GetSequencer(), big.NewInt(0))
  if err != nil {
    log.Println("Failed to create L1 transaction opts", err)
    return err
  }

  bridgeTx, err := l1BridgeComms.L1Contracts.BridgeContract.WithdrawEth(transactOpts, address, amount)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}

func (l1BridgeComms *L1Comms) BridgeTokenToL2(tokenAddress common.Address, fromAddress common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " tokens of type", tokenAddress.Hex(), " to L2 for address ", fromAddress.Hex(), " on token bridge contract", l1BridgeComms.L1ContractAddressConfig.TokenBridgeContractAddress.Hex())

  transactOpts, err := l1BridgeComms.CreateL1TransactionOpts(fromAddress, big.NewInt(0))
  if err != nil {
    log.Println("Failed to create L1 transaction opts", err)
    return err
  }

  // Using the basic ERC20 contract to approve the token bridge contract to transfer the tokens
  // NOTICE : Using hardcoded basic erc20 contract, but it is just using this to build the tx with the correct interface, so it should work with ERC20 & ERC721 tokens
  erc20Contract, err := basicerc20.NewBasicerc20(tokenAddress, l1BridgeComms.L1Client)
  if err != nil {
    log.Println("Failed to create ERC20 contract", err)
    return err
  }

  approveTx, err := erc20Contract.Approve(transactOpts, l1BridgeComms.L1ContractAddressConfig.TokenBridgeContractAddress, amount)
  if err != nil {
    log.Println("Failed to approve token bridge contract", err)
    return err
  }
  log.Println("Approved token bridge contract: ", approveTx.Hash().Hex())

  bridgeTc, err := l1BridgeComms.L1Contracts.TokenBridgeContract.DepositTokens(transactOpts, tokenAddress, amount)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTc.Hash().Hex())
  return nil
}

func (l1BridgeComms *L1Comms) BridgeTokenToL1(tokenAddress common.Address, toAddress common.Address, amount *big.Int) error {
  log.Println("Bridging ", amount, " tokens to L1 for address ", toAddress.Hex())

  transactOpts, err := l1BridgeComms.CreateL1TransactionOpts(GetSequencer(), big.NewInt(0))
  if err != nil {
    log.Println("Failed to create L1 transaction opts", err)
    return err
  }

  bridgeTx, err := l1BridgeComms.L1Contracts.TokenBridgeContract.WithdrawTokens(transactOpts, tokenAddress, toAddress, amount)
  if err != nil {
    log.Println("Failed to create bridge transaction", err)
    return err
  }

  log.Println("Bridge transaction created: ", bridgeTx.Hash().Hex())
  return nil
}
