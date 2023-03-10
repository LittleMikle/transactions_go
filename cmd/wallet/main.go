package main

import (
	"context"
	"github.com/LittleMikle/transactions_go/internal/postgres"
	"github.com/LittleMikle/transactions_go/internal/rabbit"
)

func main() {
	s := newServer()
	db, err := postgres.NewConn()
	if err != nil {
		s.echo.Logger.Fatal(err)
	}
	defer db.Close(context.Background())

	amqp, err := rabbit.NewConn()
	if err != nil {
		s.echo.Logger.Fatal(err)
	}
	defer amqp.Close()

	s.db = db
	s.amqp = amqp

	s.echo.Logger.Fatal(s.echo.Start(":8080"))
}
