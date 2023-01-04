package expense

import (
	"net/http"

	"github.com/bazsup/assessment/router"
)

func CreateExpenseHandler(c router.RouterCtx, store storer) error {
	var exp Expense
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	insertId, err := store.CreateExpense(exp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	exp.ID = insertId

	return c.JSON(http.StatusCreated, exp)
}
