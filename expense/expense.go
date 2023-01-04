package expense

import (
	"net/http"

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
	UpdateExpense(exp Expense) error
}

type CustomMiddleware struct {
	authToken string
}

func NewCustomMiddleware(authToken string) *CustomMiddleware {
	return &CustomMiddleware{authToken: authToken}
}

func (cm *CustomMiddleware) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("Authorization")

		if key != cm.authToken {
			return c.JSON(http.StatusUnauthorized, Err{Message: "Unauthorized"})
		}

		return next(c)
	}
}

func NewApp(e *echo.Echo, s storer, authToken string) {
	h := NewExpense(s)

	cm := NewCustomMiddleware(authToken)
	e.Use(cm.authMiddleware)

	e.POST("/expenses", h.CreateExpense)
	e.GET("/expenses", h.GetAllExpenses)
	e.GET("/expenses/:id", h.GetExpense)
	e.PUT("/expenses/:id", h.UpdateExpense)
}

type handler struct {
	store storer
}

func NewExpense(store storer) *handler {
	return &handler{store}
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
	return UpdateExpense(c, h.store)
}
