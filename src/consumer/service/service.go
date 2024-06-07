package service

import (
	"encoding/json"
	"fmt"

	queue "leukocyte/src/message_queue"
	"leukocyte/src/orchestration"
	"leukocyte/src/types"

	"go.uber.org/zap"
)

type Consumer struct {
	logger *zap.Logger

	orch         orchestration.Orchestration
	messageQueue queue.Queue
	routingKey   string
}

func NewConsumer(
	logger *zap.Logger, o orchestration.Orchestration,
	queue queue.Queue, routingKey string) *Consumer {
	return &Consumer{
		logger:       logger,
		orch:         o,
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

	if err = c.orch.Schedule(obj); err != nil {
		c.logger.Error("Failed to schedule job", zap.Error(err))
		goto reenqueue
	}

	goto success

reenqueue:
	if err = c.enqueue(data); err != nil {
		c.logger.Error("Failed to re enqueue message", zap.Error(err))
		return err
	}

success:
	return nil
}

func (c *Consumer) Start() error {
	if err := c.messageQueue.Consume(c.routingKey, c.Worker); err != nil {
		c.logger.Error("Failed to consume message", zap.Error(err))

		return err
	}

	return nil
}
