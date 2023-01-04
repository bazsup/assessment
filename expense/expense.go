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

type storer interface {
	CreateExpense(exp Expense) (int, error)
	GetExpenseByID(id int) (*Expense, error)
	GetAllExpenses() ([]*Expense, error)
}

type handler struct {
	DB    *sql.DB
	store storer
}

func NewExpense(db *sql.DB, store storer) *handler {
	return &handler{db, store}
}

func (h *handler) CreateExpense(c echo.Context) error {
	return CreateExpenseHandler(c, h.store)
}

func (h *handler) GetExpense(c echo.Context) error {
	return GetOneByIDHandler(c, h.store)
}

func (h *handler) GetAllExpenses(c echo.Context) error {
	return GetAllExpensesHandler(c, h.store)
}

func (h *handler) UpdateExpense(c echo.Context) error {
	return UpdateExpense(c, h.DB)
}
