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

	connectionStr string
	ctx           context.Context
	ch            *amqp.Channel
	conn          *amqp.Connection
}

func NewRabbitMQ(ctx context.Context, logger *zap.Logger, connectionStr string) *RabbitMQ {
	return &RabbitMQ{
		logger:        logger,
		connectionStr: connectionStr,
		ctx:           ctx,
		ch:            nil,
		conn:          nil,
	}
}

func (mq *RabbitMQ) Connect() error {
	if mq.ch != nil || mq.conn != nil {
		// if already connected, do nothing
		return nil
	}

	conn, err := amqp.Dial(mq.connectionStr)
	if err != nil {
		mq.logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		mq.logger.Fatal("Failed to open a channel", zap.Error(err))
		return err
	}

	mq.ch = ch
	mq.conn = conn

	return nil
}

func (mq *RabbitMQ) Close() error {
	if err := mq.ch.Close(); err != nil {
		mq.logger.Error("Failed to close channel", zap.Error(err))
		return err
	}

	if err := mq.conn.Close(); err != nil {
		mq.logger.Error("Failed to close connection", zap.Error(err))
		return err
	}

	mq.ch = nil
	mq.conn = nil

	return nil
}

func (mq *RabbitMQ) Publish(key, data string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := retry.Do(func() error {
		if err := mq.Connect(); err != nil {
			mq.logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
			return err
		}

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
	if err := mq.Connect(); err != nil {
		mq.logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))

		return err
	}

	queue, err := mq.ch.Consume(key, "", true, false, false, false, nil)
	if err != nil {
		mq.logger.Error("Failed to consume from queue", zap.Error(err))

		return err
	}

	for {
		select {
		case <-mq.ctx.Done():
			return nil

		case msg := <-queue:
			mq.logger.Debug("Received message", zap.String("body", string(msg.Body)))

			if err := callback(string(msg.Body)); err != nil {
				mq.logger.Error("Failed to process message", zap.Error(err))
			}
		}
	}
}
