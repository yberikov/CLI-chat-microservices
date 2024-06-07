package hub

import (
	"chat/internal/domain/models"
	"github.com/gorilla/websocket"
	"log/slog"
)

type Client struct {
	Username string
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan models.Message
}

func (c *Client) ReadMessage() {
	defer func() {
		c.Hub.Unregister <- c
	}()

	for {
		messageType, content, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Hub.Log.Error("Error on reading message", slog.String("err", err.Error()))
			}
			break
		}
		if messageType == websocket.CloseMessage {
			break
		}
		message, err := c.Hub.Service.SaveMessage(content, c.Username)
		if err != nil {
			c.Hub.Log.Error("Error on saving message", slog.String("err", err.Error()))
		}
		if len(c.Hub.Cache) == 10 {
			c.Hub.Cache = c.Hub.Cache[1:]
		}
		c.Hub.Cache = append(c.Hub.Cache, message)
		c.Hub.Broadcast <- message
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Hub.Unregister <- c
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				break
			}
			if message.Author != c.Username {
				toSend := message.Author + ": " + string(message.Text)
				err := c.Conn.WriteMessage(1, []byte(toSend))
				if err != nil {
					c.Hub.Log.Error("Error on writing message", slog.String("err", err.Error()))
					break
				}
			}

		}
	}
}
