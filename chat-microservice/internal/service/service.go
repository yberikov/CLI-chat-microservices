package service

import (
	"bytes"
	"chat/internal/domain/models"
	"chat/internal/storage"
	"slices"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessagerService struct {
	messager storage.Messager
	kafkaCh  chan models.Message
}

func NewService(kafkaCh chan models.Message, messager storage.Messager) *MessagerService {
	return &MessagerService{
		kafkaCh:  kafkaCh,
		messager: messager}
}

func (s *MessagerService) SaveMessage(content []byte, author string) (models.Message, error) {
	content = bytes.TrimSpace(bytes.Replace(content, newline, space, -1))
	message := models.Message{Text: content, Author: author}
	s.kafkaCh <- message
	err := s.messager.SaveMessage(message)
	if err != nil {
		return models.Message{}, err
	}
	return message, nil
}

func (s *MessagerService) GetMessages() ([]models.Message, error) {
	messages, err := s.messager.GetMessages()
	if err != nil {
		return nil, err
	}
	slices.Reverse(messages)
	return messages, nil
}

func (s *MessagerService) FillCacheFromService(serviceURL string) error {
	err := s.messager.FillCacheFromService(serviceURL)
	if err != nil {
		return err
	}
	return nil
}
