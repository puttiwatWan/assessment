// go:build integration

package database

import (
	"testing"

	"github.com/puttiwatWan/assessment/body"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateExpenseSuccess(t *testing.T) {
	db := NewDB()

	id, err := db.CreateExpense(body.Expense{
		Title:  "test title",
		Amount: 70,
		Note:   "test note",
		Tags:   []string{"test tag"},
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestIntegrationGetExpenseByIdSuccess(t *testing.T) {
	db := NewDB()

	expense, err := db.GetExpenseById(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, expense.Id)
	assert.Equal(t, "test title", expense.Title)
	assert.Equal(t, float64(70), expense.Amount)
	assert.Equal(t, "test note", expense.Note)
	assert.Equal(t, []string{"test tag"}, expense.Tags)
}

func TestIntegrationUpdateExpenseByIdSuccess(t *testing.T) {
	db := NewDB()

	expense, err := db.UpdateExpenseById(body.Expense{
		Id:     1,
		Title:  "test title2",
		Amount: 60,
		Note:   "test note2",
		Tags:   []string{"test tag2"},
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, expense.Id)
	assert.Equal(t, "test title2", expense.Title)
	assert.Equal(t, float64(60), expense.Amount)
	assert.Equal(t, "test note2", expense.Note)
	assert.Equal(t, []string{"test tag2"}, expense.Tags)
}

func TestIntegrationGetExpensesSuccess(t *testing.T) {
	db := NewDB()

	expenses, err := db.GetExpenses()

	assert.NoError(t, err)
	assert.Greater(t, len(expenses), 0)
}
