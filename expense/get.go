package expense

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/bazsup/assessment/router"
)

func GetOneByIDHandler(c router.RouterCtx, store storer) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	}

	exp, err := store.GetExpenseByID(id)

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, exp)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}

func GetAllExpensesHandler(c router.RouterCtx, storer storer) error {
	expenses, err := storer.GetAllExpenses()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, expenses)
}
