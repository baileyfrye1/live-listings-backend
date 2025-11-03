package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*WSClient]bool

type WSClient struct {
	UserId     int
	UserRole   string
	Manager    *Manager
	Connection *websocket.Conn
	Egress     chan Event
}

func NewWSClient(conn *websocket.Conn, userID int, userRole string, manager *Manager) *WSClient {
	return &WSClient{
		UserId:     userID,
		UserRole:   userRole,
		Manager:    manager,
		Connection: conn,
		Egress:     make(chan Event),
	}
}

func (c *WSClient) ReadMessages(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.Manager.RemoveClient(c)
	}()

	c.Connection.SetReadLimit(1024)
	c.Connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error unmarshalling event:  %v", err)
			continue
		}

		if err := c.Manager.RouteEvent(request, c); err != nil {
			log.Println("Error routing event: ", err)
		}
	}
}

func (c *WSClient) WriteMessages(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.Manager.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				if err := c.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed: ", err)
				}

				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling messages: %v", err)
				continue
			}

			if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("Write message error: ", err)
				return
			}
		}
	}
}

func (c *WSClient) pongHandler(pongMsg string) error {
	return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
}
