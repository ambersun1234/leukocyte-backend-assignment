package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"leukocyte/src/config"
	"leukocyte/src/logger"
	rabbitMQ "leukocyte/src/message_queue/rabbitmq"
	"leukocyte/src/producer/service"

	"go.uber.org/zap"
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

	if err := config.ReadConfig(); err != nil {
		logger.Entry.Fatal("Failed to read config file", zap.Error(err))
	}

	mq := rabbitMQ.NewRabbitMQ(ctx, logger.Entry, config.Cfg.MqURL)
	s := service.NewProducer(ctx, logger.Entry, mq, config.Cfg.MqRoutingKey)

	s.Start()

	// graceful shutdown
	<-ctx.Done()
	logger.Entry.Info("Shutting down producer...")
	s.Stop()
	mq.Close()
}
