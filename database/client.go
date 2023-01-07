package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
	"github.com/puttiwatWan/assessment/body"
)

type DBClient struct {
	client *sql.DB
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

func createTableIfNotExists(db *sql.DB) {
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

func (db *DBClient) CreateExpense(ce body.Expense) (string, error) {
	// insertExpense := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id"
	row := db.client.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", ce.Title, ce.Amount, ce.Note, pq.Array(&ce.Tags))
	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	// _, err := db.client.Exec(insertExpense, ce.Title, ce.Amount, ce.Note, pq.Array(&ce.Tags))
	return id, err
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

func (db *DBClient) UpdateExpenseById(exp body.Expense) (body.Expense, error) {
	// This is for partial update for PATCH operation
	// updateQuery := "UPDATE expenses SET"
	// valueNumber := 0
	// var updateValue []interface{}
	// if exp.Title != "" {
	// 	valueNumber += 1
	// 	updateQuery += " title=$" + strconv.Itoa(valueNumber)
	// 	updateValue = append(updateValue, &exp.Title)
	// }
	// if exp.Amount != 0 {
	// 	valueNumber += 1
	// 	updateQuery += " amount=$" + strconv.Itoa(valueNumber)
	// 	updateValue = append(updateValue, &exp.Amount)
	// }
	// if exp.Note != "" {
	// 	valueNumber += 1
	// 	updateQuery += " note=$" + strconv.Itoa(valueNumber)
	// 	updateValue = append(updateValue, &exp.Title)
	// }
	// if exp.Tags != nil && len(exp.Tags) != 0 {
	// 	valueNumber += 1
	// 	updateQuery += " tags=$" + strconv.Itoa(valueNumber)
	// 	updateValue = append(updateValue, pq.Array(&exp.Tags))
	// }
	// valueNumber += 1
	// updateQuery += " WHERE id=$" + strconv.Itoa(valueNumber) + " RETURNING id, title, amount, note, tags"
	// updateValue = append(updateValue, &exp.Id)

	row := db.client.QueryRow("UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5 RETURNING id, title, amount, note, tags;", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags), exp.Id)

	var expense body.Expense
	err := row.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		return body.Expense{}, err
	}

	return expense, nil
}
