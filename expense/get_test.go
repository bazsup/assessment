//go:build unit
// +build unit

package expense_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/bazsup/assessment/expense"
	"github.com/stretchr/testify/assert"
)

func setupExpense(t *testing.T) (*TestCtx, *TestStore) {
	t.Parallel()

	ctx := NewTestCtx(nil)
	store := NewTestStore()
	return ctx, store
}

func TestGetExpenseByID(t *testing.T) {
	t.Run("Get Expense success", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		want := expense.Expense{
			ID:     1,
			Title:  "test-title",
			Amount: 39000,
			Note:   "test-note",
			Tags:   []string{"tag1", "tag2"},
		}
		ctx.SetParam("1")

		store.GetExpenseByIDWillReturn(&want, nil)

		// Act
		err := expense.GetOneByIDHandler(ctx, store)

		var exp expense.Expense
		ctx.DecodeResponse(&exp)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, ctx.status)

			assert.Equal(t, want.ID, exp.ID)
			assert.Equal(t, want.Title, exp.Title)
			assert.Equal(t, want.Amount, exp.Amount)
			assert.Equal(t, want.Note, exp.Note)
			assert.Equal(t, want.Tags, exp.Tags)
		}
	})

	t.Run("Unknown error should returns internal server error", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		ctx.SetParam("1")

		store.GetExpenseByIDWillReturn(nil, fmt.Errorf("unknown error"))

		// Act
		err := expense.GetOneByIDHandler(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Get Expense Error Not found should returns status not found", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		ctx.SetParam("1")

		store.GetExpenseByIDWillReturn(nil, sql.ErrNoRows)

		// Act
		err := expense.GetOneByIDHandler(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, ctx.status)
			assert.Equal(t, "expense not found", errRes.Message)
		}
	})

	t.Run("Get Expense Invalid ID Param should returns status not found", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		ctx.SetParam("invalid")

		// Act
		err := expense.GetOneByIDHandler(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, ctx.status)
			assert.Equal(t, "expense not found", errRes.Message)
		}
	})
}

func TestGetAllExpenses(t *testing.T) {
	t.Run("Get All Expenses success", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		want := expense.Expense{
			ID:     1,
			Title:  "test-title",
			Amount: 39000,
			Note:   "test-note",
			Tags:   []string{"tag1", "tag2"},
		}

		store.GetAllExpensesWillReturn([]*expense.Expense{&want}, nil)

		// Act
		err := expense.GetAllExpensesHandler(ctx, store)

		var expenses []expense.Expense
		ctx.DecodeResponse(&expenses)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, ctx.status)

			assert.Equal(t, 1, len(expenses))
		}
	})

	t.Run("Get All Expenses Fail should returns status internal server error", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		store.GetAllExpensesWillReturn(nil, fmt.Errorf("fail to get all expenses"))

		// Act
		err := expense.GetAllExpensesHandler(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})
}
