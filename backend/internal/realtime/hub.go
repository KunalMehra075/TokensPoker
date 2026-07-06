package realtime

import (
	"log/slog"
	"sync"
)

// Hub tracks connected clients grouped by room and fans out events. It is the
// concrete Broadcaster implementation backed by WebSockets.
type Hub struct {
	log *slog.Logger

	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{} // roomID -> set of clients
}

// NewHub builds an empty hub.
func NewHub(log *slog.Logger) *Hub {
	return &Hub{
		log:   log,
		rooms: make(map[string]map[*Client]struct{}),
	}
}

// joinRoom registers a client in a room and returns the resulting presence list.
func (h *Hub) joinRoom(roomID string, c *Client) []MemberPresence {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]struct{})
	}
	h.rooms[roomID][c] = struct{}{}
	c.rooms[roomID] = struct{}{}
	return h.presenceLocked(roomID)
}

// leaveRoom removes a client from a single room and returns remaining presence.
func (h *Hub) leaveRoom(roomID string, c *Client) []MemberPresence {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.removeFromRoomLocked(roomID, c)
	delete(c.rooms, roomID)
	return h.presenceLocked(roomID)
}

// dropClient removes a client from every room it was in, returning the affected
// rooms with their updated presence lists so callers can notify peers.
func (h *Hub) dropClient(c *Client) map[string][]MemberPresence {
	h.mu.Lock()
	defer h.mu.Unlock()
	affected := make(map[string][]MemberPresence, len(c.rooms))
	for roomID := range c.rooms {
		h.removeFromRoomLocked(roomID, c)
		affected[roomID] = h.presenceLocked(roomID)
	}
	c.rooms = map[string]struct{}{}
	return affected
}

func (h *Hub) removeFromRoomLocked(roomID string, c *Client) {
	if set, ok := h.rooms[roomID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// MemberPresence describes one online participant.
type MemberPresence struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// presenceLocked computes the distinct online users in a room. Caller holds mu.
func (h *Hub) presenceLocked(roomID string) []MemberPresence {
	seen := map[string]struct{}{}
	out := []MemberPresence{}
	for c := range h.rooms[roomID] {
		if _, dup := seen[c.userID]; dup {
			continue
		}
		seen[c.userID] = struct{}{}
		out = append(out, MemberPresence{UserID: c.userID, Name: c.name, Email: c.email})
	}
	return out
}

// EmitToRoom sends an event to every client currently in a room. Satisfies the
// Broadcaster interface used by services.
func (h *Hub) EmitToRoom(roomID, event string, payload any) {
	env := Envelope{Event: event, RoomID: roomID, Payload: payload}
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.rooms[roomID]))
	for c := range h.rooms[roomID] {
		clients = append(clients, c)
	}
	h.mu.RUnlock()
	for _, c := range clients {
		c.enqueue(env)
	}
}
