package expense_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bazsup/assessment/expense"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) (*expense.ExpenseStore, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	expStore := expense.NewExpenseStore(db)

	return expStore, mock
}

func TestDBCreatexpense(t *testing.T) {
	exp := expense.Expense{
		ID:     1,
		Title:  "test-title",
		Amount: 39000,
		Note:   "test-note",
		Tags:   []string{"tag1", "tag2"},
	}

	t.Run("Create Expense Success", func(t *testing.T) {
		expStore, mock := setupDB(t)

		// Arrange
		expenseMockRows := sqlmock.NewRows([]string{"id"}).
			AddRow("1")
		mock.ExpectQuery("INSERT INTO expenses (.+) VALUES (.+) RETURNING id").WillReturnRows(expenseMockRows)

		// Act
		id, err := expStore.CreateExpense(exp)

		// Assertions
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, 1, id)
	})

	t.Run("Create Expense Fail", func(t *testing.T) {
		expStore, mock := setupDB(t)

		// Arrange
		mock.ExpectQuery("INSERT INTO expenses").WillReturnError(fmt.Errorf("get expenses fail"))

		// Act
		_, err := expStore.CreateExpense(exp)

		// Assertions
		assert.NotNil(t, err)
	})
}
