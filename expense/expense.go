package expense

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

type handler struct {
	DB *sql.DB
}

func NewExpense(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) CreateExpense(c echo.Context) error {
	return CreateExpenseHandler(c, h.DB)
}
