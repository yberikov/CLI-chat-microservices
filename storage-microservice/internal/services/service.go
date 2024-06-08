package service

import (
	"bytes"
	"encoding/json"
	"hw3/internal/domain/models"
	"hw3/internal/storage"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessagerService struct {
	storage storage.Messager
}

func NewService(messageStorage storage.Messager) *MessagerService {
	return &MessagerService{
		storage: messageStorage,
	}
}

func (s *MessagerService) SaveMessage(content []byte) (models.Message, error) {
	var msg models.Message
	err := json.Unmarshal(content, &msg)
	if err != nil {
		return models.Message{}, err
	}
	content = bytes.TrimSpace(bytes.Replace(msg.Text, newline, space, -1))
	message := models.Message{Text: content, Author: msg.Author}

	err = s.storage.SaveMessage(string(content), msg.Author)
	if err != nil {
		log.Println("error on: message saved in storage", string(content))
		return models.Message{}, err
	}
	log.Println("Message saved in storage", string(content))
	return message, nil
}

func (s *MessagerService) GetLastMessages(n int) ([]models.Message, error) {
	messages, err := s.storage.GetLastMessages(n)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
