package main

import (
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/LittleMikle/transactions_go/internal/rabbit"
	"github.com/LittleMikle/transactions_go/internal/request"

	"github.com/streadway/amqp"
)

type server struct {
	echo *echo.Echo
	db   *pgx.Conn
	amqp *amqp.Connection
}

func newServer() *server {
	s := &server{echo: echo.New()}
	s.routes()
	return s
}

func (s *server) handleDeposit(c echo.Context) error {
	if err := s.handleToAMQP(c, &request.Deposit{}, rabbit.Deposit); err != nil {
		return err
	}

	return c.String(
		http.StatusOK,
		"success",
	)
}

func (s *server) handleTransfer(c echo.Context) error {
	if err := s.handleToAMQP(c, &request.Transfer{}, rabbit.Transfer); err != nil {
		return err
	}

	return c.String(
		http.StatusOK,
		"success",
	)
}

func (s *server) handleToAMQP(c echo.Context, r request.RequestInterface, queueName string) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := r.Validate(s.db); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err,
		)
	}

	q, err := rabbit.NewQueue(s.amqp, queueName)
	if err != nil {
		s.echo.Logger.Fatal(err)
	}
	defer q.Close()

	body, err := json.Marshal(r)
	if err != nil {
		s.echo.Logger.Fatal(err)
	}

	if err = q.Publish(body); err != nil {
		s.echo.Logger.Fatal(err)
	}

	return nil
}
