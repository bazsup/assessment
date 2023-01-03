//go:build unit
// +build unit

package expense_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bazsup/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

// TODO: remove dependency
func init() {
	expense.InitDB()
}

func TestCreateExpense(t *testing.T) {
	t.Run("Create Expense success", func(t *testing.T) {
		// Arrange
		// TODO: grouping this
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString(`{
			"title": "test-title",
			"amount": 39000,
			"note": "test-note",
			"tags": ["test-tag1", "test-tag2"]
		}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Act
		err := expense.CreateExpenseHandler(c)

		var exp expense.Expense
		json.NewDecoder(rec.Body).Decode(&exp)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			assert.NotEqual(t, 0, exp.ID)
			assert.Equal(t, "test-title", exp.Title)
			assert.Equal(t, float64(39000), exp.Amount)
			assert.Equal(t, "test-note", exp.Note)
			assert.Equal(t, []string{"test-tag1", "test-tag2"}, exp.Tags)
		}
	})

	t.Run("Invalid Create Expense Request", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString(`xx`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Act
		err := expense.CreateExpenseHandler(c)

		var errRes expense.Err
		json.NewDecoder(rec.Body).Decode(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.NotEmpty(t, errRes.Message)
		}
	})
}
