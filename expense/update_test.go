//go:build unit
// +build unit

package expense_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/bazsup/assessment/expense"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpense(t *testing.T) {
	t.Run("Update Expense success", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)

		store.UpdateExpenseWillReturn(nil)

		// Act
		ctx.SetParam("1")
		ctx.SetReqBody(reqBody)
		err := expense.UpdateExpense(ctx, store)

		var exp expense.Expense
		ctx.DecodeResponse(&exp)

		// Assertions
		want := expense.Expense{
			ID:     1,
			Title:  "updated-title",
			Amount: 40000,
			Note:   "updated-note",
			Tags:   []string{"updated-tag"},
		}

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, ctx.status)

			assert.Equal(t, want.ID, exp.ID)
			assert.Equal(t, want.Title, exp.Title)
			assert.Equal(t, want.Amount, exp.Amount)
			assert.Equal(t, want.Note, exp.Note)
			assert.Equal(t, want.Tags, exp.Tags)
		}
	})

	t.Run("Invalid update expense request should returns status bad request", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		reqBody := bytes.NewBufferString(`invalid request`)
		ctx.SetBindErr(fmt.Errorf("bind error"))

		// Act
		ctx.SetReqBody(reqBody)
		err := expense.UpdateExpense(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Invalid expense id param should returns status not found", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)

		// Act
		ctx.SetParam("invalid param")
		ctx.SetReqBody(reqBody)
		err := expense.UpdateExpense(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, ctx.status)

			assert.Equal(t, "expense not found", errRes.Message)
		}
	})

	t.Run("Update Expense Fail should returns status internal server error", func(t *testing.T) {
		ctx, store := setupExpense(t)

		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "updated-title",
			"amount": 40000,
			"note": "updated-note",
			"tags": ["updated-tag"]
		}`)

		store.UpdateExpenseWillReturn(fmt.Errorf("can't update expense error"))

		// Act
		ctx.SetParam("1")
		ctx.SetReqBody(reqBody)
		err := expense.UpdateExpense(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)

			assert.NotEmpty(t, errRes.Message)
		}
	})

}
