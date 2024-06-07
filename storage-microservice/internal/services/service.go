package service

import (
	"bytes"
	"hw3/internal/domain/models"
	"hw3/internal/storage"
	"hw3/internal/storage/postgres"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessagerService struct {
	storage storage.Messager
}

func NewService(url string) *MessagerService {
	messageStorage, err := postgres.New(url)
	if err != nil {
		panic(err)
	}
	return &MessagerService{
		storage: messageStorage,
	}
}

func (s *MessagerService) SaveMessage(content []byte, author string) (models.Message, error) {
	content = bytes.TrimSpace(bytes.Replace(content, newline, space, -1))
	message := models.Message{Text: content, Author: author}

	err := s.storage.SaveMessage(string(content), author)
	if err != nil {
		log.Println("error on: message saved in storage", string(content))
		return models.Message{}, err
	}
	log.Println("Message saved in storage", string(content))
	return message, nil
}
