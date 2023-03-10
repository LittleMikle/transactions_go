package rabbit

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// Connects to rabbitmq using os.Getenv params. Uses github.com/streadway/amqp
func NewConn() (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s/", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"))
	return amqp.Dial(dsn)
}
