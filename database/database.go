package database

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type TransactionType int

const (
	TransactionTypeIncome  TransactionType = 1
	TransactionTypeExpense TransactionType = 2
)

func New(ctx context.Context) (*sql.DB, error) {
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("sqlite3", "./"+dbName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(nil)

	if err := goose.Up(db, "./database/migrations"); err != nil {
		return nil, err
	}

	return db, nil
}
