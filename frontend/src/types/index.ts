// Domain types mirroring the Go backend's JSON shapes.

export type EstimationMode = 'TOKENS' | 'COST' | 'DAYS' | 'MODEL'
export type TaskStatus = 'ACTIVE' | 'REVEALED' | 'CLOSED'

export interface User {
  id: string
  email: string
  name: string
  createdAt: string
  lastLogin: string
}

export interface Member {
  userId: string
  email: string
  name: string
  joinedAt: string
}

export interface Room {
  id: string
  roomCode: string
  name: string
  ownerId: string
  members: Member[]
  createdAt: string
}

export interface Task {
  id: string
  roomId: string
  title: string
  description: string
  mode: EstimationMode
  status: TaskStatus
  revealed: boolean
  createdBy: string
  createdAt: string
  revealedAt?: string
  closedAt?: string
}

export interface MemberVoteView {
  userId: string
  name: string
  hasVoted: boolean
  card?: string
}

export interface FinalDecision {
  id: string
  taskId: string
  roomId: string
  finalValue: string
  selectedBy: string
  createdAt: string
}

export interface TaskDetail {
  task: Task
  votes: MemberVoteView[]
  final?: FinalDecision
  voteCount: number
  memberCount: number
}

export interface Vote {
  id: string
  taskId: string
  roomId: string
  userId: string
  userName: string
  selectedCard: string
  createdAt: string
  updatedAt: string
}

export interface HistoryItem {
  task: Task
  roomName: string
  roomCode: string
  final?: FinalDecision
  votes: Vote[]
}

export interface ModeDefinition {
  mode: EstimationMode
  name: string
  description: string
  numeric: boolean
  cards: string[]
}

export interface RoomPreview {
  name: string
  roomCode: string
  memberCount: number
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Presence {
  userId: string
  name: string
  email: string
}

// Realtime envelope from the WebSocket hub.
export interface SocketEnvelope<T = unknown> {
  event: string
  roomId: string
  payload: T
}
