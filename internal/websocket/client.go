package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn           *websocket.Conn
	send           chan []byte
	hub            *Hub
	ConversationID string
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.Broadcast <- msg

	}
}

func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			c.hub.Unregister <- c
			break
		}
	}
}
