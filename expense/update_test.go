//go:build unit
// +build unit

package expense_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bazsup/assessment/expense"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpense(t *testing.T) {
	t.Run("Update Expense success", func(t *testing.T) {
		t.Parallel()

		// Arrange
		want := expense.Expense{
			ID:     1,
			Title:  "updated-title",
			Amount: 40000,
			Note:   "updated-note",
			Tags:   []string{"updated-tag"},
		}
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)
		ctx := NewTestCtx(reqBody)
		ctx.SetParam("1")

		database, mock, sqlErr := sqlmock.New()
		update := mock.ExpectPrepare("UPDATE .+ SET .+ WHERE id = .+")

		update.
			ExpectExec().
			WithArgs(want.ID, want.Title, want.Amount, want.Note, pq.Array(&want.Tags)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		err := expense.UpdateExpense(ctx, database)

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
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Invalid update expense request should returns status bad request", func(t *testing.T) {
		t.Parallel()

		// Arrange
		reqBody := bytes.NewBufferString(`invalid request`)
		ctx := NewTestCtx(reqBody)
		ctx.SetBindErr(fmt.Errorf("bind error"))

		database, _, _ := sqlmock.New()

		// Act
		err := expense.UpdateExpense(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Invalid expense id param should returns status not found", func(t *testing.T) {
		t.Parallel()

		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)
		ctx := NewTestCtx(reqBody)
		ctx.SetParam("invalid param")

		database, _, _ := sqlmock.New()

		// Act
		err := expense.UpdateExpense(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, ctx.status)

			assert.Equal(t, "expense not found", errRes.Message)
		}
	})

	t.Run("SQL Prepare statement error should returns status internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)
		ctx := NewTestCtx(reqBody)
		ctx.SetParam("1")

		database, mock, sqlErr := sqlmock.New()
		update := mock.ExpectPrepare("UPDATE .+ SET .+ WHERE id = .+")
		update.WillReturnError(fmt.Errorf("prepare statment error"))

		// Act
		err := expense.UpdateExpense(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SQL Execute error should returns status internal server error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)
		ctx := NewTestCtx(reqBody)
		ctx.SetParam("1")

		database, mock, sqlErr := sqlmock.New()
		update := mock.ExpectPrepare("UPDATE .+ SET .+ WHERE id = .+")
		update.ExpectExec().WillReturnError(fmt.Errorf("execute statment error"))

		// Act
		err := expense.UpdateExpense(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		assert.NoError(t, sqlErr)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
