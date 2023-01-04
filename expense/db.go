package expense

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

func InitDB() *sql.DB {
	var db *sql.DB

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	return db
}

type ExpenseStore struct {
	*sql.DB
}

func NewExpenseStore(db *sql.DB) *ExpenseStore {
	return &ExpenseStore{db}
}

func (e *ExpenseStore) CreateExpense(exp Expense) (int, error) {
	row := e.DB.QueryRow(
		"INSERT INTO expenses ( title, amount, note, tags ) VALUES ( $1, $2, $3, $4 ) RETURNING id",
		exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags))
	err := row.Scan(&exp.ID)
	return exp.ID, err
}
