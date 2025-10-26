package service

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
)

type NotificationPayload struct {
	Type    string      `json:"type"`
	Time    string      `json:"time"`
	Payload interface{} `json:"payload"`
}

type NotificationPublisher interface {
	Publish(ctx context.Context, channel string, payload interface{}) error
}

type notificationPublisher struct {
	rdb    *redis.Client
	amqpCh *amqp.Channel
}

func NewNotificationPublisher(rdb *redis.Client, amqpCh *amqp.Channel) NotificationPublisher {
	return &notificationPublisher{rdb: rdb, amqpCh: amqpCh}
}

func (p *notificationPublisher) Publish(ctx context.Context, channel string, payload interface{}) error {
	msg := &NotificationPayload{Type: channel, Time: time.Now().Format(time.RFC3339), Payload: payload}
	data, _ := json.Marshal(msg)
	// Redis Pub/Sub
	if p.rdb != nil {
		_ = p.rdb.Publish(ctx, channel, data).Err()
	}
	// RabbitMQ fanout (exchange named as channel)
	if p.amqpCh != nil {
		_ = p.amqpCh.ExchangeDeclare(channel, "fanout", true, false, false, false, nil)
		_ = p.amqpCh.PublishWithContext(ctx, channel, "", false, false, amqp.Publishing{ContentType: "application/json", Body: data})
	}
	return nil
}
