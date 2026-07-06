// Single source of truth for brand + environment wiring.

export const BRAND = {
  name: 'FreeTokensPoker',
  shortName: 'FreeTokensPoker',
  domain: 'freetokenspoker.com',
  url: 'https://freetokenspoker.com',
  tagline: 'Planning Poker for the AI era',
  contactEmail: 'hello@freetokenspoker.com',
} as const

// API + socket base URLs. VITE_API_URL points at the Go backend.
export const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
export const WS_URL =
  import.meta.env.VITE_WS_URL ?? API_BASE_URL.replace(/^http/, 'ws') + '/ws'

// Server -> client realtime event names (must match internal/realtime/events.go).
export const SOCKET_EVENTS = {
  presence: 'presence',
  memberJoined: 'member_joined',
  memberLeft: 'member_left',
  taskCreated: 'task_created',
  voteReceived: 'vote_received',
  votesRevealed: 'votes_revealed',
  taskClosed: 'task_closed',
  finalDecision: 'final_decision',
} as const

export const TOKEN_STORAGE_KEY = 'ftp.token'

// Where we stash a room code from an invite link while the invitee signs in, so
// we can resume the join right after login.
export const PENDING_ROOM_KEY = 'ftp.pendingRoomCode'

// Build a shareable invite link for a room code.
export function inviteLink(roomCode: string) {
  return `${window.location.origin}/join/${roomCode}`
}
