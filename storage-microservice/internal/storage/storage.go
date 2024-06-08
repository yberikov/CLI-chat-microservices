package storage

import "hw3/internal/domain/models"

type Messager interface {
	SaveMessage(msg string, author string) error
	GetLastMessages(n int) ([]models.Message, error)
}
