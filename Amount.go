package main

import (
	ETH_balance "Amount/ETH"
	Fantom_balance "Amount/Fantom"
	Matic_balance "Amount/Matic"
	"flag"
	"fmt"
)

func main() {
	// Define flags for Ethereum and Fantom addresses
	var (
		ethAddressFlag    = flag.String("address_ETH", "", "Ethereum address")
		fantomAddressFlag = flag.String("address_Fantom", "", "Fantom address")
		maticAddressFlag  = flag.String("address_Matic", "", "Matic address")
	)

	// Parse command-line flags
	flag.Parse()

	// Retrieve the values of the flags
	ethAddress := *ethAddressFlag
	fantomAddress := *fantomAddressFlag
	maticAddress := *maticAddressFlag

	// Check if Ethereum address flag is provided
	if ethAddress == "" {
		fmt.Println("Error: Ethereum address is required")
		return
	}

	// Check if Fantom address flag is provided
	if fantomAddress == "" {
		fmt.Println("Error: Fantom address is required")
		return
	}

	if maticAddress == "" {
		fmt.Println("Error: Matic address is required")
		return
	}

	// Call functions with these addresses
	ETH_balance.ETH_balance(ethAddress)
	Fantom_balance.Fantom_balance(fantomAddress)
	Matic_balance.Matic_balance(maticAddress)
}
