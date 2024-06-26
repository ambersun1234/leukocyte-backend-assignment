package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	queue "leukocyte/src/message_queue"
	"leukocyte/src/types"

	"go.uber.org/zap"
)

type Producer struct {
	logger *zap.Logger

	ctx          context.Context
	messageQueue queue.Queue
	ticker       *time.Ticker
	routingKey   types.RoutingKey
}

func NewProducer(
	ctx context.Context, logger *zap.Logger,
	messageQueue queue.Queue, routingKey types.RoutingKey) *Producer {
	return &Producer{
		logger:       logger,
		ctx:          ctx,
		messageQueue: messageQueue,
		ticker:       time.NewTicker(10 * time.Second),
		routingKey:   routingKey,
	}
}

func (p *Producer) produce(counter int) error {
	job := &types.JobObject{
		Namespace:     "default",
		Name:          fmt.Sprintf("job-%d", counter),
		Image:         "ubuntu",
		RestartPolicy: "Never",
		Commands:      []string{"echo", fmt.Sprintf("hello world %v !\n", counter)},
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		p.logger.Error("Failed to marshal job", zap.Error(err))
		return err
	}

	p.logger.Info("Publishing message...")
	if err := p.messageQueue.Publish(p.routingKey, string(jobBytes)); err != nil {
		p.logger.Error("Failed to publish message", zap.Error(err))
		return err
	}

	return nil
}

func (p *Producer) Start() {
	counter := 0

	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			if err := p.produce(counter); err != nil {
				p.logger.Error("Failed to produce message", zap.Error(err))
			}

			counter += 1
		}
	}
}

func (p *Producer) Stop() {
	p.ticker.Stop()
}
