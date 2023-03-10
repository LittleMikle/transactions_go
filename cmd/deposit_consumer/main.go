package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/LittleMikle/transactions_go/internal/postgres"
	"github.com/LittleMikle/transactions_go/internal/rabbit"
	"github.com/LittleMikle/transactions_go/internal/request"
)

func main() {
	db, err := postgres.NewConn()
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer db.Close(context.Background())

	amqp, err := rabbit.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer amqp.Close()

	q, err := rabbit.NewQueue(amqp, rabbit.Deposit)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	deliveries, err := q.Consume()
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}
	go func() {
		for d := range deliveries {
			r := &request.Deposit{}
			if err = json.Unmarshal(d.Body, &r); err != nil {
				log.Fatalf("failed to bind delivery data %#v to binding %s: %v", d.Body, rabbit.Deposit, err)
			}

			commandTag, err := db.Exec(context.Background(), "UPDATE wallets SET amount = amount + $1 WHERE id = $2", r.Amount, r.Receiver)
			if err != nil {
				log.Fatalf("deposit operation failed: %v", err)
			}
			if commandTag.RowsAffected() != 1 {
				log.Fatalf("wallet %v not found!", r.Receiver)
			}
		}
	}()
	<-forever
}
