package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
	"github.com/puttiwatWan/assessment/body"
)

type Database interface {
	Close() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DBClient struct {
	client Database
}

func NewDB() *DBClient {
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTableIfNotExists(db)

	return &DBClient{
		client: db,
	}
}

func createTableIfNotExists(db Database) {
	createTable := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`
	_, err := db.Exec(createTable)
	if err != nil {
		log.Fatal("Create table error", err)
	}
}

func (db *DBClient) CreateExpense(ce body.Expense) error {
	insertExpense := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)"
	_, err := db.client.Exec(insertExpense, ce.Title, ce.Amount, ce.Note, pq.Array(&ce.Tags))
	return err
}

func (db *DBClient) CloseConnection() {
	db.client.Close()
}
