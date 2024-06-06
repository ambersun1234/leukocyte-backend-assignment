package rabbitmq

import (
	"context"
	"time"

	"leukocyte/src/types"

	"github.com/avast/retry-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	logger *zap.Logger

	ch   *amqp.Channel
	conn *amqp.Connection
}

func NewRabbitMQ(logger *zap.Logger, connectionStr string) *RabbitMQ {
	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("Failed to open a channel", zap.Error(err))
	}

	return &RabbitMQ{
		logger: logger,
		ch:     ch,
		conn:   conn,
	}
}

func (mq *RabbitMQ) Publish(key, data string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := retry.Do(func() error {
		return mq.ch.PublishWithContext(
			ctx, "", key, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(data)},
		)
	}, retry.Attempts(3))

	if err != nil {
		mq.logger.Error("Failed to publish message", zap.Error(err))

		return err
	}

	mq.logger.Debug("Published message", zap.String("body", data))

	return nil
}

func (mq *RabbitMQ) Consume(key string, callback types.CallbackFunc) error {
	queue, err := mq.ch.Consume(key, "", true, false, false, false, nil)
	if err != nil {
		mq.logger.Error("Failed to consume from queue", zap.Error(err))

		return err
	}

	for msg := range queue {
		mq.logger.Debug("Received message", zap.String("body", string(msg.Body)))

		if err := callback(string(msg.Body)); err != nil {
			mq.logger.Error("Failed to process message", zap.Error(err))

			continue
		}
	}

	return nil
}
