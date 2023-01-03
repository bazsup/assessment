//go:build integration
// +build integration

package expense_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bazsup/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var serverPort = 30001

type teardownFunc = func(t *testing.T)

func setup() teardownFunc {
	eh := echo.New()
	go func(e *echo.Echo) {
		db := expense.InitDB()

		h := expense.NewExpense(db)

		e.POST("/expenses", h.CreateExpense)
		e.GET("/expenses/:id", h.GetExpense)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	return func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := eh.Shutdown(ctx)
		assert.NoError(t, err)
	}
}

func TestITCreateExpense(t *testing.T) {
	// Setup server
	teardown := setup()
	defer teardown(t)

	// Arrange
	reqBody := `{
		"title": "test-title",
		"amount": 39000,
		"note": "test-note",
		"tags": ["test-tag1", "test-tag2"]
	}`

	// Act
	res := request(http.MethodPost, uri("expenses"), strings.NewReader(reqBody))

	var exp expense.Expense
	err := res.Decode(&exp)
	assert.NoError(t, err)
	res.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "test-title", exp.Title)
		assert.Equal(t, float64(39000), exp.Amount)
		assert.Equal(t, "test-note", exp.Note)
		assert.Equal(t, []string{"test-tag1", "test-tag2"}, exp.Tags)
	}
}

func TestITGetExpense(t *testing.T) {
	// Setup server
	teardown := setup()
	defer teardown(t)

	// Arrange
	exp := seedExpense(t)

	// Act
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(exp.ID)), nil)

	var latest expense.Expense
	err := res.Decode(&latest)
	assert.NoError(t, err)
	res.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, exp.ID, latest.ID)
		assert.Equal(t, exp.Title, latest.Title)
		assert.Equal(t, exp.Amount, latest.Amount)
		assert.Equal(t, exp.Note, latest.Note)
		assert.Equal(t, exp.Tags, latest.Tags)
	}
}

func seedExpense(t *testing.T) expense.Expense {
	var c expense.Expense
	body := strings.NewReader(`{
		"title": "test-title",
		"amount": 39000,
		"note": "test-note",
		"tags": ["test-tag1", "test-tag2"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}

	return c
}

func uri(paths ...string) string {
	host := fmt.Sprintf("http://localhost:%d", serverPort)
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)

	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
