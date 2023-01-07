package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/puttiwatWan/assessment/body"
	"github.com/stretchr/testify/assert"
)

type mockDatabase struct {
	CloseFn    func() error
	ExecFn     func(query string, args ...interface{}) (sql.Result, error)
	QueryRowFn func(query string, args ...interface{}) *sql.Row
	PrepareFn  func(query string) (*sql.Stmt, error)
}

func (m *mockDatabase) Close() error {
	return m.CloseFn()
}

func (m *mockDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.ExecFn(query, args)
}

func (m *mockDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.QueryRowFn(query, args)
}

func (m *mockDatabase) Prepare(query string) (*sql.Stmt, error) {
	return m.PrepareFn(query)
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

type mockSqlStmt struct {
	QueryFn func(args ...interface{}) (*sql.Rows, error)
}

func (m *mockSqlStmt) Query(args ...interface{}) (*sql.Rows, error) {
	return m.QueryFn(args...)
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

func TestGetExpenseByIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"id", "title", "amount", "note", "tags"}
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses WHERE id = ").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows(columns).AddRow("1", "test title", 70, "test note", pq.Array([]string{"test tag"})))

	mockDb := DBClient{client: db}

	expense, err := mockDb.GetExpenseById("1")

	assert.NoError(t, err)
	assert.Equal(t, "1", expense.Id)
	assert.Equal(t, "test title", expense.Title)
	assert.Equal(t, float64(70), expense.Amount)
	assert.Equal(t, "test note", expense.Note)
	assert.Equal(t, []string{"test tag"}, expense.Tags)
}

func TestGetExpenseByIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses WHERE id = ").
		WithArgs("1").
		WillReturnError(fmt.Errorf("boom"))

	mockDb := DBClient{client: db}

	_, err = mockDb.GetExpenseById("1")

	assert.Error(t, err)
	assert.EqualError(t, err, "boom")
}
