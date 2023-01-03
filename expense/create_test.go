//go:build unit
// +build unit

package expense_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bazsup/assessment/expense"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateExpense(t *testing.T) {
	t.Run("Create Expense success", func(t *testing.T) {
		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "test-title",
			"amount": 39000,
			"note": "test-note",
			"tags": ["test-tag1", "test-tag2"]
		}`)
		ctx := NewTestCtx(reqBody)
		expenseMockRows := sqlmock.NewRows([]string{"id"}).
			AddRow("1")

		database, mock, err := sqlmock.New()
		mock.ExpectQuery("INSERT INTO expenses (.+) VALUES (.+) RETURNING id").WillReturnRows(expenseMockRows)

		// Act
		err = expense.CreateExpenseHandler(ctx, database)

		var exp expense.Expense
		ctx.DecodeResponse(&exp)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, ctx.status)

			assert.NotEqual(t, 0, exp.ID)
			assert.Equal(t, "test-title", exp.Title)
			assert.Equal(t, float64(39000), exp.Amount)
			assert.Equal(t, "test-note", exp.Note)
			assert.Equal(t, []string{"test-tag1", "test-tag2"}, exp.Tags)
		}
	})

	t.Run("Invalid Create Expense Request", func(t *testing.T) {
		// Arrange
		reqBody := bytes.NewBufferString(`xx`)
		ctx := NewTestCtx(reqBody)
		ctx.SetBindErr(fmt.Errorf("bind error"))
		database, _, _ := sqlmock.New()

		// Act
		err := expense.CreateExpenseHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("SQL error should return status code internal server error", func(t *testing.T) {
		// Arrange
		reqBody := bytes.NewBufferString(`{
			"title": "test-title",
			"amount": 39000,
			"note": "test-note",
			"tags": ["test-tag1", "test-tag2"]
		}`)
		ctx := NewTestCtx(reqBody)
		database, _, _ := sqlmock.New()

		// Act
		err := expense.CreateExpenseHandler(ctx, database)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})
}

type TestCtx struct {
	req     *bytes.Buffer
	status  int
	v       []byte
	bindErr error
}

func NewTestCtx(req *bytes.Buffer) *TestCtx {
	return &TestCtx{req: req}
}

func (c *TestCtx) SetBindErr(err error) {
	c.bindErr = err
}

func (c *TestCtx) Bind(v interface{}) error {
	json.NewDecoder(c.req).Decode(&v)
	return c.bindErr
}

func (c *TestCtx) JSON(code int, v interface{}) error {
	c.status = code
	data, err := json.Marshal(v)
	c.v = data
	return err
}

func (c *TestCtx) DecodeResponse(res interface{}) error {
	return json.Unmarshal(c.v, res)
}
