package service

import (
	"bytes"
	"chat/internal/domain/models"
	"chat/internal/storage"
	"chat/internal/storage/redis"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessagerService struct {
	storage        storage.Messager
	messageChannel chan models.Message
}

func NewService(messageChannel chan models.Message, url string) *MessagerService {
	messageStorage, err := redis.New(url)
	if err != nil {
		log.Panicln(err)
	}
	return &MessagerService{
		messageChannel: messageChannel,
		storage:        messageStorage}
}

func (s *MessagerService) SaveMessage(content []byte, author string) (models.Message, error) {
	content = bytes.TrimSpace(bytes.Replace(content, newline, space, -1))
	message := models.Message{Text: content, Author: author}
	s.messageChannel <- message
	err := s.storage.SaveMessage(string(content), author)
	if err != nil {
		return models.Message{}, err
	}
	
	return message, nil
}
