package kafka

import (
	internalConfig "chat/internal/config"
	"chat/internal/domain/models"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

type Producer struct {
	log *slog.Logger
	prd sarama.AsyncProducer
	cfg *internalConfig.Config
	ch  chan models.Message
}

func NewProducer(logger *slog.Logger, cfg *internalConfig.Config, ch chan models.Message) Producer {
	//TODO kafka configuration
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.ClientID = "chat-microservice-1"
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewAsyncProducer(strings.Split(cfg.Brokers, ","), config)
	if err != nil {
		logger.Error("Failed to start Sarama producer:", err)
	}

	return Producer{
		log: logger,
		prd: producer,
		cfg: cfg,
		ch:  ch,
	}
}

func (p *Producer) RunProducing(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-p.ch:
			jsonData, err := json.Marshal(message)
			if err != nil {
				p.log.Error("Failed to marshal message:", err)
				continue
			}
			p.prd.Input() <- &sarama.ProducerMessage{
				Topic: p.cfg.Topic,
				Value: sarama.ByteEncoder(jsonData),
			}
		}

	}

}
