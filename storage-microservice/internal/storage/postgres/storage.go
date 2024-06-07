package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type MessageStorage struct {
	db *pgx.Conn
}

// New creates a new SQL connection.
func New(url string) (*MessageStorage, error) {

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	// Verify the connection with a Ping
	if err := conn.Ping(context.Background()); err != nil {
		err := conn.Close(context.Background())
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &MessageStorage{db: conn}, nil
}

func (c *MessageStorage) SaveMessage(msg string, author string) error {
	_, err := c.db.Exec(context.Background(), "INSERT INTO messages (text, author) VALUES ($1, $2)", msg, author)
	return err
}

func (c *MessageStorage) Close() error {
	return c.db.Close(context.Background())
}
