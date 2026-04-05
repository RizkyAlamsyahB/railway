// Package websocket provides a WebSocket Hub for real-time push notifications.
//
// Architecture:
//   - One Hub per server process, shared across all connections.
//   - Each authenticated user can have multiple concurrent Client connections
//     (e.g. multiple browser tabs).
//   - When a chat message is saved via REST, the sender's handler calls
//     Hub.SendToUser(recipientID, message) to push it in real-time.
//   - Client → Server: only ping/pong keepalive (no message sending over WS).
//   - Server → Client: JSON-encoded WSMessage payloads.
package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Timing constants for connection health.
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// MessageType identifies the kind of push payload.
type MessageType string

const (
	TypeChatMessage         MessageType = "chat_message"
	TypeTicketUpdate        MessageType = "ticket_update"
	TypeNotificationCreated MessageType = "notification.created"
)

// WSMessage is the envelope sent to the client over WebSocket.
type WSMessage struct {
	Type           MessageType     `json:"type"`
	ConversationID string          `json:"conversation_id,omitempty"`
	TicketID       string          `json:"ticket_id,omitempty"`
	Data           json.RawMessage `json:"data"`
}

// BroadcastTarget is an outbound message destined for a specific user.
type BroadcastTarget struct {
	RecipientID uuid.UUID
	Payload     []byte
}

// Client represents a single WebSocket connection for one user.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID uuid.UUID
}

// Hub manages all active WebSocket clients.
type Hub struct {
	// clients maps userID → set of active connections for that user.
	clients    map[uuid.UUID]map[*Client]struct{}
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastTarget
}

// NewHub creates and starts a new Hub. Call once at startup.
func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[uuid.UUID]map[*Client]struct{}),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan BroadcastTarget, 512),
	}
	go h.run()
	return h
}

// run is the Hub's event loop. Must run in its own goroutine.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; !ok {
				h.clients[client.UserID] = make(map[*Client]struct{})
			}
			h.clients[client.UserID][client] = struct{}{}
			h.mu.Unlock()
			log.Printf("[ws] user %s connected (%d connections)", client.UserID, len(h.clients[client.UserID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.UserID]; ok {
				delete(conns, client)
				if len(conns) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			close(client.send)
			log.Printf("[ws] user %s disconnected", client.UserID)

		case target := <-h.broadcast:
			h.mu.RLock()
			conns := h.clients[target.RecipientID]
			h.mu.RUnlock()
			for c := range conns {
				select {
				case c.send <- target.Payload:
				default:
					// slow client — drop and disconnect
					h.unregister <- c
				}
			}
		}
	}
}

// SendToUser encodes msg as JSON and queues it for delivery to all connections
// owned by recipientID. Safe to call from any goroutine.
func (h *Hub) SendToUser(recipientID uuid.UUID, msg WSMessage) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[ws] marshal error: %v", err)
		return
	}
	h.broadcast <- BroadcastTarget{RecipientID: recipientID, Payload: payload}
}

// readPump drains inbound frames (we only care about pong control frames).
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		// We only read to detect disconnects; discard any client messages.
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

// writePump fans out queued messages and sends periodic pings.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			// Flush any queued messages in the same write frame.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte("\n"))
				_, _ = w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
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

// ServeClient registers the client and starts its read/write pumps.
// Call this after the HTTP upgrade.
func (h *Hub) ServeClient(conn *websocket.Conn, userID uuid.UUID) {
	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
	}
	h.register <- client
	go client.writePump()
	client.readPump() // blocks until disconnect
}
