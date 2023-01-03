//go:build integration
// +build integration

package expense_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bazsup/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var serverPort = 30001

func TestITCreateExpense(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		expense.InitDB()

		e.POST("/expenses", expense.CreateExpenseHandler)
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
	// Arrange
	reqBody := `{
		"title": "test-title",
		"amount": 39000,
		"note": "test-note",
		"tags": ["test-tag1", "test-tag2"]
	}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// byteBody, err := ioutil.ReadAll(resp.Body)
	var exp expense.Expense
	json.NewDecoder(resp.Body).Decode(&exp)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "test-title", exp.Title)
		assert.Equal(t, float64(39000), exp.Amount)
		assert.Equal(t, "test-note", exp.Note)
		assert.Equal(t, []string{"test-tag1", "test-tag2"}, exp.Tags)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
