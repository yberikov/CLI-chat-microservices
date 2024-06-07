package storage

type Messager interface {
	SaveMessage(msg string, author string) error
}
