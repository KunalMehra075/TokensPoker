package realtime

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
	sendBuffer     = 32
)

// Client is a single WebSocket connection bound to one authenticated user.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	log    *slog.Logger
	send   chan Envelope
	userID string
	name   string
	email  string
	rooms  map[string]struct{}
}

func (c *Client) enqueue(env Envelope) {
	select {
	case c.send <- env:
	default:
		// Slow consumer: drop the connection rather than block the hub.
		c.log.Warn("dropping slow websocket client", "userId", c.userID)
		close(c.send)
	}
}

// readPump handles inbound control messages (join/leave/ping) and connection
// lifecycle. It runs until the socket closes.
func (c *Client) readPump() {
	defer func() {
		affected := c.hub.dropClient(c)
		for roomID, presence := range affected {
			c.hub.EmitToRoom(roomID, EventPresence, presence)
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var msg inbound
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		switch msg.Type {
		case "join":
			if msg.RoomID == "" {
				continue
			}
			presence := c.hub.joinRoom(msg.RoomID, c)
			c.hub.EmitToRoom(msg.RoomID, EventPresence, presence)
		case "leave":
			if msg.RoomID == "" {
				continue
			}
			presence := c.hub.leaveRoom(msg.RoomID, c)
			c.hub.EmitToRoom(msg.RoomID, EventPresence, presence)
		case "ping":
			// Client-level heartbeat; reply is handled by writePump pings.
		}
	}
}

// writePump serializes outbound events and protocol pings to the socket.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case env, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(env); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
