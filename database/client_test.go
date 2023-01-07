package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/puttiwatWan/assessment/body"
	"github.com/stretchr/testify/assert"
)

type mockDatabase struct {
	CloseFn    func() error
	ExecFn     func(query string, args ...interface{}) (sql.Result, error)
	QueryRowFn func(query string, args ...interface{}) *sql.Row
}

func (m *mockDatabase) Close() error {
	return m.CloseFn()
}

func (m *mockDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.ExecFn(query, args)
}

func (m *mockDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.QueryRow(query, args)
}

type mockSqlResult struct {
	LastInsertIdFn func() (int64, error)
	RowsAffectedFn func() (int64, error)
}

func (m *mockSqlResult) LastInsertId() (int64, error) {
	return m.LastInsertIdFn()
}

func (m *mockSqlResult) RowsAffected() (int64, error) {
	return m.RowsAffectedFn()
}

func TestCreateExpenseSuccess(t *testing.T) {
	client := &mockDatabase{}
	client.ExecFn = func(query string, args ...interface{}) (sql.Result, error) {
		return &mockSqlResult{}, nil
	}
	db := DBClient{client: client}

	err := db.CreateExpense(body.Expense{})

	assert.NoError(t, err)
}

func TestCreateExpenseError(t *testing.T) {
	client := &mockDatabase{}
	client.ExecFn = func(query string, args ...interface{}) (sql.Result, error) {
		return nil, fmt.Errorf("boom")
	}
	db := DBClient{client: client}

	err := db.CreateExpense(body.Expense{
		Title:  "Spaghetti",
		Amount: 70,
		Note:   "food",
		Tags:   []string{"daily"},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "boom")
}
