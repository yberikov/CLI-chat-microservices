package hub

import (
	"context"
	"log/slog"
	"sync"

	"chat/internal/domain/models"
	"chat/internal/service"
)

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan models.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	Service *service.MessagerService
	Cache   []models.Message
	Log     *slog.Logger
}

func NewHub(log *slog.Logger, service *service.MessagerService) *Hub {
	return &Hub{
		Broadcast:  make(chan models.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Service:    service,
		Cache:      []models.Message{},
		Log:        log,
	}
}

func (h *Hub) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case user := <-h.Register:
			h.Log.Info("New user connected:", slog.String("username", user.Username))
			h.Clients[user] = true
		case user := <-h.Unregister:
			if _, ok := h.Clients[user]; ok {
				h.Log.Info("User disconnected:", slog.String("username", user.Username))
				delete(h.Clients, user)
				close(user.Send)
				err := user.Conn.Close()
				if err != nil {
					h.Log.Error("Error on closing websocket connetction")
				}
			}
		case message := <-h.Broadcast:
			h.Log.Info("Message recevied by broadcast:", slog.String("username", message.Author),
				slog.String("message", string(message.Text)))
			for user := range h.Clients {
				select {
				case user.Send <- message:
				default:
					close(user.Send)
					delete(h.Clients, user)
				}
			}
		}
	}
}
