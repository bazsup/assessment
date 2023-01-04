package expense

type storer interface {
	CreateExpense(exp Expense) (int, error)
}