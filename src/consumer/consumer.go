package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"leukocyte/src/config"
	"leukocyte/src/consumer/service"
	"leukocyte/src/logger"
	rabbitMQ "leukocyte/src/message_queue/rabbitmq"
	"leukocyte/src/orchestration/k8s"

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

	// initialize consumer
	if err := logger.NewLogger(); err != nil {
		panic(err)
	}

	if err := config.ReadConfig(); err != nil {
		logger.Entry.Fatal("Failed to read config file", zap.Error(err))
	}

	container := k8s.NewK8s(logger.Entry, config.Cfg.Kubernetes.InCluster, config.Cfg.Kubernetes.ConfigUrl)
	mq := rabbitMQ.NewRabbitMQ(ctx, logger.Entry, config.Cfg.MessageQueue.Url)
	s := service.NewConsumer(logger.Entry, container, mq, config.Cfg.MessageQueue.RoutingKey)

	if err := s.Start(); err != nil {
		logger.Entry.Fatal("Failed to start consumer", zap.Error(err))
	}

	// graceful shutdown
	<-ctx.Done()
	logger.Entry.Info("Shutting down consumer...")
	mq.Close()
}
