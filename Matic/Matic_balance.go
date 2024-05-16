// В файле Amount/ETH/balance.go

package Matic_balance

import (
	"Amount/Utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	//"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func Matic_balance(address string) {
	// Open a database connection
	db, err := Utils.ConnectDB("postgres://postgres:fkZ8pkzw@localhost:5432/postgres?sslmode=disable") // Используем функцию ConnectDB из Utils
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Подключаемся к узлу Ethereum
	client, err := Utils.GetEthClient("https://rpc.ankr.com/polygon")
	if err != nil {
		log.Fatal(err)
	}

	// Адрес Ethereum, баланс которого мы хотим отслеживать
	addr := common.HexToAddress(address) // Rename the variable to addr

	// Текущий номер блока
	currentBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Количество блоков в день (приблизительно 7200 блоков в день при среднем времени блока в 12 секунд)
	blocksPerDay := 7200

	// Получаем сегодняшнюю дату
	currentDate := time.Now()

	fmt.Printf("Matic\n")

	// Получаем баланс адреса за каждый из последних 7 дней
	for i := 6; i >= 0; i-- {
		// Номер блока для дня, отстоящего на i дней назад
		targetBlockNumber := currentBlockNumber - uint64(blocksPerDay*(i))

		// Получаем баланс адреса на заданном блоке
		balance, err := Utils.GetBalanceAtBlock(client, addr.String(), big.NewInt(int64(targetBlockNumber))) // Используем функцию GetBalanceAtBlock из Utils
		if err != nil {
			log.Fatal(err)
		}

		// Получаем дату для текущего дня (i дней назад)
		targetDate := currentDate.AddDate(0, 0, -1*(i))

		// Форматируем баланс в ETH с нужным представлением
		formattedBalance := Utils.FormatBalance(balance)

		// Вставляем данные в таблицу PostgreSQL
		insertData(db, addr.String(), targetDate.Format("2006-01-02"), formattedBalance)
	}

	fmt.Printf("\n")
}

// Функция для вставки данных в таблицу PostgreSQL, с предварительной проверкой наличия записи
func insertData(db *sql.DB, address string, date string, balance string) {
	// Проверяем, существует ли запись с таким адресом и датой
	var existingAddress, existingDate, existingBalance string
	err := db.QueryRow(`SELECT "Address", "Date", "Balance" FROM "Matic" WHERE "Address" = $1 AND "Date" = $2`, address, date).Scan(&existingAddress, &existingDate, &existingBalance)
	switch {
	case err == sql.ErrNoRows:
		// Если запись не найдена, выполняем вставку
		stmt, err := db.Prepare(`INSERT INTO "Matic" ("Address", "Date", "Balance") VALUES ($1, $2, $3)`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(address, date, balance)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Новая запись добавлена в базу данных.")
	case err != nil:
		log.Fatal(err)
	default:
		// Если запись уже существует, выводим ее
		fmt.Printf("Запись уже существует:Адрес: %s Дата: %s Баланс: %s Matic\n", existingAddress, existingDate, existingBalance)
	}
}
