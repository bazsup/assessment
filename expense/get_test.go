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

func TestGetExpenseByID(t *testing.T) {
	t.Run("Get Expense success", func(t *testing.T) {
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
		ctx.SetParam("1")

		expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(want.ID, want.Title, want.Amount, want.Note, pq.Array(&want.Tags))

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.ExpectQuery().WithArgs("1").WillReturnRows(expenseMockRows)

		// Act
		err := expense.GetOneByIDHandler(ctx, database)

		var exp expense.Expense
		ctx.DecodeResponse(&exp)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, ctx.status)

			assert.Equal(t, want.ID, exp.ID)
			assert.Equal(t, want.Title, exp.Title)
			assert.Equal(t, want.Amount, exp.Amount)
			assert.Equal(t, want.Note, exp.Note)
			assert.Equal(t, want.Tags, exp.Tags)
		}
	})

	t.Run("Prepare statment error should returns internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := NewTestCtx(nil)
		ctx.SetParam("1")

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.WillReturnError(fmt.Errorf("error prepare statement"))

		// Act
		err := expense.GetOneByIDHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Scan query error should returns internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := NewTestCtx(nil)
		ctx.SetParam("1")

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.ExpectQuery().WithArgs("1").WillReturnError(fmt.Errorf("error query"))

		// Act
		err := expense.GetOneByIDHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Query success but not found should returns status not found", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := NewTestCtx(nil)
		ctx.SetParam("1")

		database, mock, sqlErr := sqlmock.New()
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.ExpectQuery().WithArgs("1").WillReturnError(sql.ErrNoRows)

		// Act
		err := expense.GetOneByIDHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
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
