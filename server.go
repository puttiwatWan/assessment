package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/puttiwatWan/assessment/body"
	"github.com/puttiwatWan/assessment/database"
	"github.com/puttiwatWan/assessment/ports"
)

var db ports.DBOperations

const IdQueryParam = "id"

type Err struct {
	Message string `json:"message"`
}

func setUpDB() {
	db = database.NewDB()
}

func createExpenseHandler(c echo.Context) error {
	var in body.Expense
	err := c.Bind(&in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	err = db.CreateExpense(in)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, in)
}

func getExpensesByIdHandler(c echo.Context) error {
	id := c.Param(IdQueryParam)

	expense, err := db.GetExpenseById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, expense)
}

func main() {
	port := os.Getenv("PORT")

	setUpDB()
	defer db.CloseConnection()

	e := echo.New()
	e.POST("/expenses", createExpenseHandler)
	e.GET("/expenses/:"+IdQueryParam, getExpensesByIdHandler)

	// Start server
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Shutdown gracefully")
}
