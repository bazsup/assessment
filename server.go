package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bazsup/assessment/config"
	"github.com/bazsup/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.NewConfig()
	db := expense.InitDB(config.DatabaseUrl)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	store := expense.NewExpenseStore(db)
	expense.NewApp(e, store, config.AuthToken)

	go func() {
		if err := e.Start(config.Port); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
		log.Println("bye bye!")
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
