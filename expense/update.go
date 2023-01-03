package expense

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/bazsup/assessment/router"
	"github.com/lib/pq"
)

func UpdateExpense(c router.RouterCtx, database *sql.DB) error {
	var exp Expense
	if err := c.Bind(&exp); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	}
	exp.ID = id

	stmt, err := database.Prepare(`
	UPDATE expenses
	SET title = $2, amount = $3, note = $4, tags = $5
	WHERE id = $1
	`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statement:" + err.Error()})
	}

	_, err = stmt.Exec(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}
