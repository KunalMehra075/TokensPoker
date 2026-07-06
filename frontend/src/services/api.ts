import axios from 'axios'
import type { AxiosError } from 'axios'
import { API_BASE_URL } from '@/constants'
import { getToken, useAuthStore } from '@/store/auth'
import type {
  AuthResponse,
  HistoryItem,
  ModeDefinition,
  Room,
  RoomPreview,
  TaskDetail,
  User,
  Vote,
  EstimationMode,
} from '@/types'

// The backend wraps success payloads as { success, data }. unwrap pulls data.
interface Envelope<T> {
  success: boolean
  data: T
}

export interface ApiErrorShape {
  message: string
  errorCode: string
  status: number
}

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res,
  (error: AxiosError<{ message?: string; errorCode?: string }>) => {
    // An expired or invalid identity token logs the user out cleanly.
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
    }
    const shaped: ApiErrorShape = {
      message: error.response?.data?.message ?? 'Something went wrong. Please try again.',
      errorCode: error.response?.data?.errorCode ?? 'UNKNOWN',
      status: error.response?.status ?? 0,
    }
    return Promise.reject(shaped)
  },
)

async function unwrap<T>(p: Promise<{ data: Envelope<T> }>): Promise<T> {
  const res = await p
  return res.data.data
}

// --- Auth ---
export const login = (name: string, email: string) =>
  unwrap<AuthResponse>(api.post('/api/auth/login', { name, email }))

export const fetchMe = () => unwrap<User>(api.get('/api/me'))

// --- Meta ---
export const fetchModes = () => unwrap<ModeDefinition[]>(api.get('/api/modes'))

// --- Rooms ---
export const createRoom = (name: string) =>
  unwrap<Room>(api.post('/api/rooms', { name }))

export const joinRoom = (roomCode: string) =>
  unwrap<Room>(api.post('/api/rooms/join', { roomCode }))

// Public: preview a room behind an invite link before signing in.
export const fetchRoomPreview = (code: string) =>
  unwrap<RoomPreview>(api.get(`/api/invite/${encodeURIComponent(code)}`))

export const fetchRoom = (id: string) => unwrap<Room>(api.get(`/api/rooms/${id}`))

export const deleteRoom = (id: string) =>
  unwrap<{ deleted: boolean }>(api.delete(`/api/rooms/${id}`))

// --- Tasks ---
export const fetchRoomTasks = (roomId: string) =>
  unwrap<TaskDetail[]>(api.get(`/api/rooms/${roomId}/tasks`))

export const createTask = (input: {
  roomId: string
  title: string
  description?: string
  mode: EstimationMode
}) => unwrap<TaskDetail>(api.post('/api/tasks', input))

export const fetchTask = (id: string) =>
  unwrap<TaskDetail>(api.get(`/api/tasks/${id}`))

export const revealTask = (id: string) =>
  unwrap<TaskDetail>(api.patch(`/api/tasks/${id}/reveal`))

export const finalizeTask = (id: string, finalValue: string) =>
  unwrap<TaskDetail>(api.patch(`/api/tasks/${id}/final`, { finalValue }))

// --- Votes ---
export const submitVote = (taskId: string, selectedCard: string) =>
  unwrap<Vote>(api.post('/api/votes', { taskId, selectedCard }))

// --- History ---
export const fetchHistory = () => unwrap<HistoryItem[]>(api.get('/api/history'))
