package main

import (
	"context"
	"github.com/IBM/sarama"
	"hw3/internal/config"
	"hw3/internal/kafka"
	service2 "hw3/internal/services"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Storage-microservice started")
	service := service2.NewService(cfg.StoragePath)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	consumer, err := kafka.RunConsumer(ctx, wg, cfg.Brokers, service)
	if err != nil {
		log.Fatalln(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("stopping server")

	cancel()
	wg.Wait()
	err = consumer.Close()
	logger.Info("server stopped")
}
