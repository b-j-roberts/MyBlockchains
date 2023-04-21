package ledger

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Airdrop(airdropFile string) map[uint64]uint64 {
  // One Trillion
  TotalSupply := uint64(1000000000000)

  accounts := make(map[uint64]uint64)

  if airdropFile != "" {
    f, err := os.Open(airdropFile)
    if err != nil {
      return accounts
    }
    defer f.Close()

    fmt.Println("----------------------------------------")
    fmt.Println("Starting Airdrop:")

    scanner := bufio.NewScanner(f)
    total := uint64(0)
    for scanner.Scan() {
      line := scanner.Text()
      parts := strings.Split(line, ",")

      address, err := strconv.ParseUint(parts[0], 16, 64)
      if err != nil {
        fmt.Println(err)
        return accounts
      }

      percentage, err := strconv.ParseFloat(parts[1], 64)
      if err != nil {
        fmt.Println(err)
        return accounts
      }

      total += uint64(float64(TotalSupply) * percentage)

      accounts[uint64(address)] = uint64(float64(TotalSupply) * percentage)
      fmt.Println("    Airdrop: ", address, accounts[uint64(address)])
    }

    if err := scanner.Err(); err != nil {
      fmt.Println(err)
    }


    fmt.Println("Expected Airdrop Amount:", TotalSupply)
    fmt.Println("Total Airdrop Amount:", total)
    fmt.Println("----------------------------------------")

  } else {
    fmt.Println("----------------------------------------")
    fmt.Println("Starting Airdrop:")

    for i := uint64(0); i < 5; i++ {
      accounts[i] = TotalSupply / 5
      fmt.Println("    Airdrop: ", i, accounts[i])
    }

    fmt.Println("Total Airdrop Amount:", TotalSupply)
    fmt.Println("----------------------------------------")
  }

  return accounts
}
