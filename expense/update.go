package expense

import (
	"net/http"
	"strconv"

	"github.com/bazsup/assessment/router"
)

func UpdateExpense(c router.RouterCtx, store storer) error {
	var exp Expense
	if err := c.Bind(&exp); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	}
	exp.ID = id

	if err = store.UpdateExpense(exp); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, exp)
}
