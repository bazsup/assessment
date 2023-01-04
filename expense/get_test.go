//go:build unit
// +build unit

package expense_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bazsup/assessment/expense"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupGetOneExpense(t *testing.T) (*TestCtx, *TestStore) {
	t.Parallel()

	ctx := NewTestCtx(nil)
	store := NewTestStore()
	return ctx, store
}

func TestGetExpenseByID(t *testing.T) {
	t.Run("Get Expense success", func(t *testing.T) {
		ctx, store := setupGetOneExpense(t)

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
		ctx, store := setupGetOneExpense(t)

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
		ctx, store := setupGetOneExpense(t)

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
		ctx, store := setupGetOneExpense(t)

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
		t.Parallel()

		// Arrange
		want := expense.Expense{
			ID:     1,
			Title:  "test-title",
			Amount: 39000,
			Note:   "test-note",
			Tags:   []string{"tag1", "tag2"},
		}
		ctx := NewTestCtx(nil)

		expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(want.ID, want.Title, want.Amount, want.Note, pq.Array(&want.Tags))

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses")
		get.ExpectQuery().WillReturnRows(expenseMockRows)

		// Act
		err := expense.GetAllExpensesHandler(ctx, database)

		var expenses []expense.Expense
		ctx.DecodeResponse(&expenses)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, ctx.status)

			assert.Equal(t, 1, len(expenses))
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare statement error should returns status internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := NewTestCtx(nil)

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses")
		get.WillReturnError(fmt.Errorf("prepare stmt error"))

		// Act
		err := expense.GetAllExpensesHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Database Query error should returns status internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := NewTestCtx(nil)

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses")
		get.ExpectQuery().WillReturnError(fmt.Errorf("query error"))

		// Act
		err := expense.GetAllExpensesHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Scan for entity error should returns status internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		want := expense.Expense{
			ID:     1,
			Title:  "test-title",
			Amount: 39000,
			Note:   "test-note",
			Tags:   []string{"tag1", "tag2"},
		}
		ctx := NewTestCtx(nil)

		expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow("invalid", want.Title, want.Amount, want.Note, pq.Array(&want.Tags))
		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses")
		get.ExpectQuery().WillReturnRows(expenseMockRows)

		// Act
		err := expense.GetAllExpensesHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})
}
