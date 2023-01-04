//go:build unit
// +build unit

package expense_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bazsup/assessment/expense"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func setupCreateExpense(t *testing.T, reqBody *bytes.Buffer) (*TestCtx, *TestStore) {
	t.Parallel()

	ctx := NewTestCtx(reqBody)
	store := NewTestStore()
	return ctx, store
}

func TestCreateExpense(t *testing.T) {
	t.Run("Create Expense success", func(t *testing.T) {
		// Arrange
		validReqBody := bytes.NewBufferString(`{
			"title": "test-title",
			"amount": 39000,
			"note": "test-note",
			"tags": ["test-tag1", "test-tag2"]
		}`)
		ctx, store := setupCreateExpense(t, validReqBody)

		store.CreateExpenseWillReturn(1, nil)

		// Act
		err := expense.CreateExpenseHandler(ctx, store)

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
		invalidReqBody := bytes.NewBufferString(`xx`)
		ctx, store := setupCreateExpense(t, invalidReqBody)

		ctx.SetBindErr(fmt.Errorf("bind error"))

		// Act
		err := expense.CreateExpenseHandler(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})

	t.Run("Create Expense should return status code internal server error", func(t *testing.T) {
		// Arrange
		validReqBody := bytes.NewBufferString(`{
			"title": "test-title",
			"amount": 39000,
			"note": "test-note",
			"tags": ["test-tag1", "test-tag2"]
		}`)
		ctx, store := setupCreateExpense(t, validReqBody)

		store.CreateExpenseWillReturn(0, fmt.Errorf("fail to create expense"))

		// Act
		err := expense.CreateExpenseHandler(ctx, store)

		var errRes expense.Err
		ctx.DecodeResponse(&errRes)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, ctx.status)
			assert.NotEmpty(t, errRes.Message)
		}
	})
}

type TestStore struct {
	ctr  *CreateExpenseTestResult
	gotr *GetOneExpenseTestResult
	gatr *GetAllExpensesTestResult
}

func NewTestStore() *TestStore {
	return &TestStore{}
}

func (s *TestStore) CreateExpense(exp expense.Expense) (int, error) {
	return s.ctr.id, s.ctr.err
}

func (s *TestStore) CreateExpenseWillReturn(id int, err error) {
	s.ctr = &CreateExpenseTestResult{id, err}
}

func (s *TestStore) GetExpenseByID(id int) (*expense.Expense, error) {
	return s.gotr.exp, s.gotr.err
}

func (s *TestStore) GetExpenseByIDWillReturn(exp *expense.Expense, err error) {
	s.gotr = &GetOneExpenseTestResult{exp, err}
}

func (s *TestStore) GetAllExpenses() ([]*expense.Expense, error) {
	return s.gatr.exp, s.gatr.err
}

func (s *TestStore) GetAllExpensesWillReturn(expenses []*expense.Expense, err error) {
	s.gatr = &GetAllExpensesTestResult{expenses, err}
}

type CreateExpenseTestResult struct {
	id  int
	err error
}

type GetOneExpenseTestResult struct {
	exp *expense.Expense
	err error
}

type GetAllExpensesTestResult struct {
	exp []*expense.Expense
	err error
}

type TestCtx struct {
	req     *bytes.Buffer
	status  int
	v       []byte
	bindErr error
	param   string
}

func NewTestCtx(req *bytes.Buffer) *TestCtx {
	return &TestCtx{req: req}
}

func (c *TestCtx) SetParam(value string) {
	c.param = value
}

func (c *TestCtx) Param(name string) string {
	return c.param
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
