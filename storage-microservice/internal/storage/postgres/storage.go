package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"hw3/internal/domain/models"
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

func (c *MessageStorage) SaveMessage(msg, author string) error {
	_, err := c.db.Exec(context.Background(), "INSERT INTO messages (text, author) VALUES ($1, $2)", msg, author)
	return err
}

func (c *MessageStorage) GetLastMessages(n int) ([]models.Message, error) {
	rows, err := c.db.Query(context.Background(), "SELECT text, author FROM messages ORDER BY id DESC LIMIT $1", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg, author string
		if err := rows.Scan(&msg, &author); err != nil {
			return nil, err
		}
		messages = append(messages, models.Message{Text: []byte(msg), Author: author})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *MessageStorage) Close() error {
	return c.db.Close(context.Background())
}
