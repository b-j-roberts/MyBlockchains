package utils

import (
	"github.com/ethereum/go-ethereum/common"
)

type TokenAddresses struct {
  Erc20Address common.Address
  L2Erc20Address common.Address
  StableErc20Address common.Address
  L2StableErc20Address common.Address
  Erc721Address common.Address
  L2Erc721Address common.Address
  SpecialErc721Address common.Address
  L2SpecialErc721Address common.Address
}

func LoadTokenAddresses(contractAddressDir string) (TokenAddresses, error) {
  var tokenAddresses TokenAddresses
  var err error

  tokenAddresses.Erc20Address, err = ReadContractAddressFromFile(contractAddressDir + "/basic-erc20-address.txt")
  if err != nil {
    return tokenAddresses, err
  }
  tokenAddresses.L2Erc20Address, err = ReadContractAddressFromFile(contractAddressDir + "/l2-basic-erc20-address.txt")
  if err != nil {
    return tokenAddresses, err
  }

  tokenAddresses.StableErc20Address, err = ReadContractAddressFromFile(contractAddressDir + "/stable-erc20-address.txt")
  if err != nil {
    return tokenAddresses, err
  }
  tokenAddresses.L2StableErc20Address, err = ReadContractAddressFromFile(contractAddressDir + "/l2-stable-erc20-address.txt")
  if err != nil {
    return tokenAddresses, err
  }

  tokenAddresses.Erc721Address, err = ReadContractAddressFromFile(contractAddressDir + "/basic-erc721-address.txt")
  if err != nil {
    return tokenAddresses, err
  }
  tokenAddresses.L2Erc721Address, err = ReadContractAddressFromFile(contractAddressDir + "/l2-basic-erc721-address.txt")
  if err != nil {
    return tokenAddresses, err
  }

  tokenAddresses.SpecialErc721Address, err = ReadContractAddressFromFile(contractAddressDir + "/special-erc721-address.txt")
  if err != nil {
    return tokenAddresses, err
  }
  tokenAddresses.L2SpecialErc721Address, err = ReadContractAddressFromFile(contractAddressDir + "/l2-special-erc721-address.txt")
  if err != nil {
    return tokenAddresses, err
  }

  return tokenAddresses, nil
}
