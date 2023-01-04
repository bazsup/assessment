package expense_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bazsup/assessment/expense"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) (*expense.ExpenseStore, sqlmock.Sqlmock) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	expStore := expense.NewExpenseStore(db)

	return expStore, mock
}

func TestDBCreatExpense(t *testing.T) {
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

func TestDBGetExpenseByID(t *testing.T) {
	t.Run("Get Expense By ID Success", func(t *testing.T) {
		expStore, mock := setupDB(t)

		// Arrange
		want := expense.Expense{
			ID:     1,
			Title:  "test-title",
			Amount: 39000,
			Note:   "test-note",
			Tags:   []string{"tag1", "tag2"},
		}

		expenseMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(want.ID, want.Title, want.Amount, want.Note, pq.Array(&want.Tags))
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.ExpectQuery().WithArgs(1).WillReturnRows(expenseMockRows)

		// Act
		exp, err := expStore.GetExpenseByID(1)

		// Assertions
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, want.ID, exp.ID)
		assert.Equal(t, want.Title, exp.Title)
		assert.Equal(t, want.Amount, exp.Amount)
		assert.Equal(t, want.Note, exp.Note)
		assert.Equal(t, want.Tags, exp.Tags)
	})

	t.Run("Prepare Statement Error", func(t *testing.T) {
		expStore, mock := setupDB(t)

		// Arrange
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.WillReturnError(fmt.Errorf("error prepare statement"))

		// Act
		exp, err := expStore.GetExpenseByID(1)

		// Assertions
		assert.Nil(t, exp)
		assert.NotNil(t, err)
	})

	t.Run("Error NoRows", func(t *testing.T) {
		expStore, mock := setupDB(t)

		// Arrange
		get := mock.ExpectPrepare("SELECT .+ FROM expenses WHERE id = .+")
		get.ExpectQuery().WithArgs(1).WillReturnError(sql.ErrNoRows)

		// Act
		exp, err := expStore.GetExpenseByID(1)

		// Assertions
		assert.Nil(t, exp)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}