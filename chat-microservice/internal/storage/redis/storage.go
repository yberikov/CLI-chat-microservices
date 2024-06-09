package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"chat/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

type MessageStorage struct {
	client *redis.Client
}

func New(addr string) (*MessageStorage, error) {
	// TODO redis configuration
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return &MessageStorage{
		client: rdb,
	}, nil
}

func (c *MessageStorage) FillCacheFromService(serviceURL string) error {
	resp, err := http.Get(serviceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get messages: %s", resp.Status)
	}

	var messages []models.Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return err
	}
	slices.Reverse(messages)
	for _, message := range messages {
		if err := c.SaveMessage(message); err != nil {
			return err
		}
	}

	return nil
}

func (c *MessageStorage) SaveMessage(message models.Message) error {
	ctx := context.Background()
	// Generate a new incrementing key
	newKey, err := c.client.Incr(ctx, "message_key").Result()
	if err != nil {
		return err
	}

	// Serialize the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Store the message with the new key
	messageKey := fmt.Sprintf("message:%d", newKey)
	err = c.client.Set(ctx, messageKey, jsonData, 0).Err()
	if err != nil {
		return err
	}

	// Add the new key to the list of message keys
	err = c.client.LPush(ctx, "message_keys", messageKey).Err()
	if err != nil {
		return err
	}

	// Trim the list to only keep the last 10 keys
	err = c.client.LTrim(ctx, "message_keys", 0, 9).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *MessageStorage) GetMessages() ([]models.Message, error) {
	ctx := context.Background()
	// Retrieve the last 10 message keys
	keys, err := c.client.LRange(ctx, "message_keys", 0, 9).Result()
	if err != nil {
		return nil, err
	}

	var messages []models.Message
	for _, key := range keys {
		jsonData, err := c.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var msg models.Message
		err = json.Unmarshal([]byte(jsonData), &msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (c *MessageStorage) Close() error {
	return c.client.Close()
}
