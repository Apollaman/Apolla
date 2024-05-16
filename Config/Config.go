package config

import (
	"flag"
)

// ETH возвращает адрес Ethereum в виде строки
func ETH() string {
	var address_ETH string
	flag.StringVar(&address_ETH, "address_ETH", "", "Ethereum address")
	flag.Parse()

	return address_ETH
}

func Fantom() string {
	var address_Fantom string
	flag.StringVar(&address_Fantom, "address_Fantom", "", "Fantom address")
	flag.Parse()

	return address_Fantom

}
