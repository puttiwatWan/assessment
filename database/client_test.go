package database

import (
	"database/sql"
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

	db := DBClient{Client: client}

	err := db.CreateExpense(body.Expense{})

	assert.NoError(t, err)
}
