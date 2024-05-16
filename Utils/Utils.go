package Utils

import (
	"database/sql"

	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnectDB(connectionURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func GetEthClient(nodeURL string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetBalanceAtBlock(client *ethclient.Client, address string, blockNumber *big.Int) (*big.Int, error) {
	addr := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), addr, blockNumber)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func FormatBalance(balance *big.Int) string {
	ether := new(big.Float).SetInt(balance)
	etherQuo := new(big.Float).Quo(ether, big.NewFloat(1e18)) // конвертируем wei в ether
	return etherQuo.Text('f', 18)                             // форматируем баланс с 18 десятичными знаками
}
