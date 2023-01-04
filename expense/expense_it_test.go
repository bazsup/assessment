//go:build integration
// +build integration

package expense_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bazsup/assessment/config"
	"github.com/bazsup/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var serverPort = 30001

type teardownFunc = func(t *testing.T)

func setup() teardownFunc {
	eh := echo.New()
	go func(e *echo.Echo) {
		config := config.NewConfig()
		db := expense.InitDB(config.DatabaseUrl)
		store := expense.NewExpenseStore(db)

		expense.NewApp(e, store, config.AuthToken)

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

	byteBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// Assertions
	want := `{"id":\d+,"title":"test-title","amount":39000,"note":"test-note","tags":\["test-tag1","test-tag2"\]}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Regexp(t, want, strings.TrimSpace(string(byteBody)))
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

	byteBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// Assertions
	want := fmt.Sprintf(
		`{"id":%d,"title":"test-title","amount":39000,"note":"test-note","tags":["test-tag1","test-tag2"]}`,
		exp.ID,
	)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, want, strings.TrimSpace(string(byteBody)))
	}
}

func TestITGetAllExpenses(t *testing.T) {
	// Setup server
	teardown := setup()
	defer teardown(t)

	// Arrange
	seedExpense(t)

	// Act
	res := request(http.MethodGet, uri("expenses"), nil)

	var body []interface{}
	err := res.Decode(&body)

	// Assertions
	if assert.NoError(t, err) {
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(body), 0)
	}
}

func TestITUpdateExpense(t *testing.T) {
	// Setup server
	teardown := setup()
	defer teardown(t)

	// Arrange
	exp := seedExpense(t)
	reqBody := `{
		"title": "updated-title",
		"amount": 40000,
		"note": "updated-note",
		"tags": ["updated-tag"]
	}`

	// Act
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(exp.ID)), strings.NewReader(reqBody))

	byteBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// Assertions
	want := fmt.Sprintf(
		`{"id":%d,"title":"updated-title","amount":40000,"note":"updated-note","tags":["updated-tag"]}`,
		exp.ID,
	)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)

		assert.Equal(t, want, strings.TrimSpace(string(byteBody)))
	}
}

func TestITAuthTokenRequired(t *testing.T) {
	// Setup server
	teardown := setup()
	defer teardown(t)

	// Arrange
	req, _ := http.NewRequest(http.MethodGet, uri("expenses"), nil)
	req.Header.Add("Authorization", "November 10, 2009wrong_token")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	// Act
	res, err := client.Do(req)
	res.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
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

var conf = config.NewConfig()

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)

	req.Header.Add("Authorization", conf.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
