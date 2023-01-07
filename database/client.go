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
	Prepare(query string) (*sql.Stmt, error)
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

func (db *DBClient) CloseConnection() {
	db.client.Close()
}

func (db *DBClient) CreateExpense(ce body.Expense) error {
	insertExpense := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)"
	_, err := db.client.Exec(insertExpense, ce.Title, ce.Amount, ce.Note, pq.Array(&ce.Tags))
	return err
}

func (db *DBClient) GetExpenseById(id string) (body.Expense, error) {
	row := db.client.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1", id)

	var expense body.Expense
	err := row.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		return body.Expense{}, err
	}

	return expense, nil
}
