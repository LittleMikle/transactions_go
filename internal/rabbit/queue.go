package rabbit

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// List of durable and persistent rabbitmq queues
const (
	Deposit  = "deposit"
	Transfer = "transfer"
)

type queue struct {
	name    string
	channel *amqp.Channel
}

func NewQueue(amqp *amqp.Connection, name string) (*queue, error) {
	ch, err := amqp.Channel()
	if err != nil {
		return nil, errors.New("failed to open RabbitMQ channel")
	}

	if _, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("failed to declare AMQP queue %s", name)
	}

	return &queue{
		channel: ch,
		name:    name,
	}, nil
}

func (q *queue) Publish(body []byte) error {
	if err := q.channel.Publish(
		"",
		q.GetName(),
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish to AMQP channel %s: %#v", q.GetName(), body)
	}
	return nil
}

func (q *queue) Consume() (<-chan amqp.Delivery, error) {
	deliveries, err := q.channel.Consume(
		q.GetName(),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (q *queue) GetName() string {
	return q.name
}

func (q *queue) Close() {
	q.channel.Close()
}
