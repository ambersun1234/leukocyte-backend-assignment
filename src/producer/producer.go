package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"leukocyte/src/logger"
	rabbitMQ "leukocyte/src/message_queue/rabbitmq"
	"leukocyte/src/producer/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	// initialize producer
	if err := logger.NewLogger(); err != nil {
		panic(err)
	}

	mq := rabbitMQ.NewRabbitMQ(ctx, logger.Entry, "amqp://rabbitmq:rabbitmq@localhost:5672/")
	s := service.NewProducer(ctx, logger.Entry, mq)

	s.Start()

	// graceful shutdown
	<-ctx.Done()
	logger.Entry.Info("Shutting down producer...")
	s.Stop()
	mq.Close()
}
