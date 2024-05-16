// В файле Amount/ETH/balance.go

package ETH_balance

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"

	"Amount/Utils"

	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func ETH_balance(address_ETH string) {
	// Open a database connection
	db, err := Utils.ConnectDB("postgres://postgres:fkZ8pkzw@localhost:5432/postgres?sslmode=disable") // Используем функцию ConnectDB из Utils
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Подключаемся к узлу Ethereum
	client, err := Utils.GetEthClient("https://mainnet.infura.io/v3/61663dee5a4e40849561d81327665b37") // Используем функцию GetEthClient из Utils
	if err != nil {
		log.Fatal(err)
	}

	// Адрес Ethereum, баланс которого мы хотим отслеживать
	addr := common.HexToAddress(address_ETH) // Rename the variable to addr

	// Текущий номер блока
	currentBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Количество блоков в день (приблизительно 7200 блоков в день при среднем времени блока в 12 секунд)
	blocksPerDay := 7200

	// Получаем сегодняшнюю дату
	currentDate := time.Now()

	fmt.Printf("Etherium\n")

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

func insertData(db *sql.DB, address string, date string, balance string) {
	// Подготовка SQL запроса с условием вставки
	stmt, err := db.Prepare(`
        INSERT INTO "Etherium" ("Address", "Date", "Balance")
        VALUES ($1, $2, $3)
        ON CONFLICT ("Address", "Date") DO UPDATE
        SET "Balance" = EXCLUDED."Balance"
        RETURNING "Address", "Date", "Balance"
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Выполнение запроса
	var existingAddress, existingDate, existingBalance string
	err = stmt.QueryRow(address, date, balance).Scan(&existingAddress, &existingDate, &existingBalance)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка результатов запроса
	if existingAddress == "" {
		fmt.Println("Новая запись добавлена в базу данных.")
	} else {
		fmt.Printf("Запись уже существует: Адрес: %s Дата: %s Баланс: %s ETH\n", existingAddress, existingDate, existingBalance)
	}
}
