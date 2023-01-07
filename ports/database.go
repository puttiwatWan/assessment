package ports

import "github.com/puttiwatWan/assessment/body"

type DBOperations interface {
	CloseConnection()
	CreateExpense(ce body.Expense) (string, error)
	GetExpenseById(id string) (body.Expense, error)
}
