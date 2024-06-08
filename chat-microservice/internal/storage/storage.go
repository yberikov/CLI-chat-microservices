package storage

import "chat/internal/domain/models"

type Messager interface {
	SaveMessage(models.Message) error
	GetMessages() ([]models.Message, error)
}
