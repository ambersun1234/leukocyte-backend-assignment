package service

import (
	"encoding/json"
	"fmt"

	"leukocyte/src/container"
	queue "leukocyte/src/message_queue"
	"leukocyte/src/types"

	"go.uber.org/zap"
)

type Consumer struct {
	logger *zap.Logger

	container    container.Container
	messageQueue queue.Queue
	routingKey   string
}

func NewConsumer(
	logger *zap.Logger, c container.Container,
	queue queue.Queue, routingKey string) *Consumer {
	return &Consumer{
		logger:       logger,
		container:    c,
		messageQueue: queue,
		routingKey:   routingKey,
	}
}

func (c *Consumer) enqueue(data string) error {
	if err := c.messageQueue.Publish(c.routingKey, data); err != nil {
		c.logger.Error("Failed to re enqueue message", zap.Error(err))
		return err
	}

	return nil
}

func (c *Consumer) Worker(data string) error {
	c.logger.Info(fmt.Sprintf("Received message: %s", data))

	var (
		err error
		obj types.JobObject
	)
	if err = json.Unmarshal([]byte(data), &obj); err != nil {
		c.logger.Error("Failed to unmarshal job", zap.Error(err))
		goto reenqueue
	}

	if err = c.container.Schedule(obj); err != nil {
		c.logger.Error("Failed to schedule job", zap.Error(err))
		goto reenqueue
	}

	return nil

reenqueue:
	if err = c.enqueue(data); err != nil {
		c.logger.Error("Failed to re enqueue message", zap.Error(err))
		return err
	}
	return nil
}

func (c *Consumer) Start() error {
	if err := c.messageQueue.Consume(c.routingKey, c.Worker); err != nil {
		c.logger.Error("Failed to consume message", zap.Error(err))

		return err
	}

	return nil
}
