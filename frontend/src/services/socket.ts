import { WS_URL } from '@/constants'
import type { SocketEnvelope } from '@/types'

type Listener = (env: SocketEnvelope) => void

// RealtimeClient is a thin, reconnecting WebSocket wrapper. It speaks the same
// typed-envelope protocol as the Go hub: inbound control messages
// ({type:'join'|'leave'|'ping'}) and outbound events ({event,roomId,payload}).
class RealtimeClient {
  private ws: WebSocket | null = null
  private token: string | null = null
  private listeners = new Set<Listener>()
  private joinedRooms = new Set<string>()
  private reconnectAttempts = 0
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private heartbeat: ReturnType<typeof setInterval> | null = null
  private manuallyClosed = false

  connect(token: string) {
    // Reconnect only when the token changes or the socket is closed.
    if (this.ws && this.token === token && this.ws.readyState <= WebSocket.OPEN) {
      return
    }
    this.token = token
    this.manuallyClosed = false
    this.open()
  }

  private open() {
    if (!this.token) return
    const ws = new WebSocket(`${WS_URL}?token=${encodeURIComponent(this.token)}`)
    this.ws = ws

    ws.onopen = () => {
      this.reconnectAttempts = 0
      // Re-join any rooms we were in before a reconnect.
      for (const roomId of this.joinedRooms) {
        this.send({ type: 'join', roomId })
      }
      this.heartbeat = setInterval(() => this.send({ type: 'ping' }), 25000)
    }

    ws.onmessage = (e) => {
      try {
        const env = JSON.parse(e.data) as SocketEnvelope
        this.listeners.forEach((l) => l(env))
      } catch {
        // Ignore malformed frames.
      }
    }

    ws.onclose = () => {
      if (this.heartbeat) clearInterval(this.heartbeat)
      if (!this.manuallyClosed) this.scheduleReconnect()
    }

    ws.onerror = () => {
      ws.close()
    }
  }

  private scheduleReconnect() {
    if (this.reconnectTimer) return
    const delay = Math.min(1000 * 2 ** this.reconnectAttempts, 15000)
    this.reconnectAttempts += 1
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      this.open()
    }, delay)
  }

  private send(msg: Record<string, unknown>) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg))
    }
  }

  joinRoom(roomId: string) {
    this.joinedRooms.add(roomId)
    this.send({ type: 'join', roomId })
  }

  leaveRoom(roomId: string) {
    this.joinedRooms.delete(roomId)
    this.send({ type: 'leave', roomId })
  }

  subscribe(listener: Listener) {
    this.listeners.add(listener)
    return () => this.listeners.delete(listener)
  }

  disconnect() {
    this.manuallyClosed = true
    this.joinedRooms.clear()
    if (this.heartbeat) clearInterval(this.heartbeat)
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer)
    this.ws?.close()
    this.ws = null
  }
}

export const realtime = new RealtimeClient()
