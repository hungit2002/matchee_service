package rabbitmq

import (
	"matchee/services/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

func NewConnection(cfg *config.Config) (*amqp091.Connection, *amqp091.Channel, error) {
	conn, err := amqp091.Dial(cfg.RabbitURL)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
}
