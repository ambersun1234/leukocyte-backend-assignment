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
}

func NewConsumer(logger *zap.Logger, c container.Container, queue queue.Queue) *Consumer {
	return &Consumer{
		logger:       logger,
		container:    c,
		messageQueue: queue,
	}
}

func (c *Consumer) Worker(data string) error {
	c.logger.Info(fmt.Sprintf("Received message: %s", data))

	obj := &types.JobObject{}
	if err := json.Unmarshal([]byte(data), obj); err != nil {
		c.logger.Error("Failed to unmarshal job", zap.Error(err))
		return err
	}

	if err := c.container.Schedule(*obj); err != nil {
		c.logger.Error("Failed to schedule job", zap.Error(err))
		return err
	}

	return nil
}

func (c *Consumer) Start() error {
	if err := c.messageQueue.Consume("test", c.Worker); err != nil {
		c.logger.Error("Failed to consume message", zap.Error(err))

		return err
	}

	return nil
}