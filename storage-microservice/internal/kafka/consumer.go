package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	service "hw3/internal/services"
	"log"
	"os"
	"strings"
	"sync"
)

type MessageHandler func(message *sarama.ConsumerMessage) error

type Consumer struct {
	handler MessageHandler
}

func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func NewConsumer(handler MessageHandler) *Consumer {
	return &Consumer{
		handler: handler,
	}
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("message channel was closed")
			}
			err := consumer.handler(message)
			if err != nil {
				fmt.Println(err)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			session.Commit()
			return nil
		}
	}
}

func InitConsumerConfig() *sarama.Config {
	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

var (
	group  = "1"
	topics = "messages"
)

func RunConsumer(ctx context.Context, wg *sync.WaitGroup, brokers string, service *service.MessagerService) (sarama.ConsumerGroup, error) {
	consumer := NewConsumer(func(message *sarama.ConsumerMessage) error {
		_, err := service.SaveMessage(message.Value)
		if err != nil {
			return err
		}
		log.Printf("Message claimed: value = %s, time = %s, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		return nil
	})
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, InitConsumerConfig())
	if err != nil {
		log.Fatalln(err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, strings.Split(topics, ","), consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
			}
			if ctx.Err() != nil {
				return
			}

		}
	}()
	return consumerGroup, nil
}
