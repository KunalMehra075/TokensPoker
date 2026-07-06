package realtime

// Event names broadcast from server to clients. They mirror the domain actions
// described in the architecture doc so the UI can stay event-driven.
const (
	EventPresence      = "presence"       // full list of online members in a room
	EventMemberJoined  = "member_joined"  // a user joined the room (persistent membership)
	EventMemberLeft    = "member_left"    // a user's socket disconnected
	EventTaskCreated   = "task_created"   // a new estimation task opened
	EventVoteReceived  = "vote_received"  // someone cast/changed a vote (value hidden pre-reveal)
	EventVotesRevealed = "votes_revealed" // all votes for the active task are now visible
	EventTaskClosed    = "task_closed"    // a task was closed
	EventFinalDecision = "final_decision" // owner committed a final value
)

// Envelope is the wire format for every server -> client message.
type Envelope struct {
	Event   string `json:"event"`
	RoomID  string `json:"roomId"`
	Payload any    `json:"payload"`
}

// inbound is the wire format for client -> server control messages.
type inbound struct {
	Type   string `json:"type"`   // "join" | "leave" | "ping"
	RoomID string `json:"roomId"` // target room for join/leave
}

// Broadcaster is the realtime port that services depend on. It is intentionally
// transport-agnostic so the WebSocket hub can be swapped without touching
// business logic.
type Broadcaster interface {
	EmitToRoom(roomID, event string, payload any)
}

// NoopBroadcaster is used in tests or when realtime is disabled.
type NoopBroadcaster struct{}

// EmitToRoom does nothing.
func (NoopBroadcaster) EmitToRoom(string, string, any) {}
