package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type MessageStorage struct {
	client *redis.Client
}

func New(url string) (*MessageStorage, error) {
	//TODO redis configuration
	rdb := redis.NewClient(&redis.Options{
		Addr: url,
	})

	status, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalln("Redis connection was refused")
	}
	fmt.Println(status)
	return &MessageStorage{
		client: rdb,
	}, nil
}

func (c *MessageStorage) SaveMessage(msg string, author string) error {
	c.client.Set(context.TODO(), msg, author, time.Hour)
	return nil
}

func (c *MessageStorage) Close() error {
	return c.client.Close()
}
