package expense

import (
	"database/sql"
	"net/http"

	"github.com/bazsup/assessment/router"
	"github.com/lib/pq"
)

func CreateExpenseHandler(c router.RouterCtx, database *sql.DB) error {
	var exp Expense
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := database.QueryRow(
		"INSERT INTO expenses ( title, amount, note, tags ) VALUES ( $1, $2, $3, $4 ) RETURNING id",
		exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags))
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, exp)
}
