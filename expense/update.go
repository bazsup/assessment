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
	_ = c.Bind(&exp)                     // TODO "Invalid update expense request should returns status bad request"
	id, _ := strconv.Atoi(c.Param("id")) // TODO "Invalid expense id param should returns status not found"
	exp.ID = id
	stmt, err := database.Prepare(`
	UPDATE expenses
	SET title = $2, amount = $3, note = $4, tags = $5
	WHERE id = $1
	`) // TODO "SQL Prepare statement error should returns status internal server error"
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statement:" + err.Error()})
	}

	// TODO "SQL Execute error should returns status internal server error"
	_, _ = stmt.Exec(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))

	return c.JSON(http.StatusOK, exp)
}
