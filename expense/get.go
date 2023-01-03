package expense

import (
	"database/sql"
	"net/http"

	"github.com/bazsup/assessment/router"
	"github.com/lib/pq"
)

func GetOneByIDHandler(c router.RouterCtx, database *sql.DB) error {
	id := c.Param("id")
	stmt, err := database.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statement:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	exp := Expense{}
	err = row.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, exp)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}

func GetAllExpensesHandler(c router.RouterCtx, database *sql.DB) error {
	stmt, _ := database.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	// TODO "Prepare statement error should returns status internal server error"

	rows, _ := stmt.Query()
	// TODO "Database Query error should returns status internal server error"

	expenses := []Expense{}
	for rows.Next() {
		var exp Expense
		// TODO "Scan for entity error should returns status internal server error"
		rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))

		expenses = append(expenses, exp)
	}

	return c.JSON(http.StatusOK, expenses)
}
