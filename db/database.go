package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	utils "github.com/brongh/go-telegram-financial-bot/utils"
)

type Expense struct {
	Id int
	UserId int
	Title string
	Amount float64
	ExpenseDate string
}

type User struct {
	Id int
	Username string
	TgId int
}

var DB *sql.DB

func GetDB() *sql.DB {
	return DB
}

func DbInit() error {
	config := utils.ReadConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", 
		config.Host, config.Port, config.DbUser, config.DbPassword, config.DbName)

	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		return err
	}
	err = setupSchema(DB)
	if err != nil {
		return fmt.Errorf("Error creating tables: %v", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")

	return nil
}

func setupSchema(db *sql.DB) error {
	schemaStatements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			tg_id INTEGER NOT NULL UNIQUE
		);`,
		`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			title VARCHAR(255) NOT NULL,
			amount NUMERIC NOT NULL,
			expense_date DATE NOT NULL
		);`,
		//
	}

	for _, stmt := range schemaStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}