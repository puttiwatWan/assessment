package ports

import "github.com/puttiwatWan/assessment/body"

type DBOperations interface {
	CloseConnection()
	CreateExpense(ce body.Expense) (int, error)
	GetExpenseById(id int) (body.Expense, error)
	UpdateExpenseById(exp body.Expense) (body.Expense, error)
	GetExpenses() ([]body.Expense, error)
}
