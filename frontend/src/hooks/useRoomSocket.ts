import { useEffect } from 'react'
import { realtime } from '@/services/socket'
import { useAuthStore } from '@/store/auth'
import type { SocketEnvelope } from '@/types'

// useRoomSocket connects the realtime client, joins a room, and routes inbound
// events to a handler for the duration of the component's life.
export function useRoomSocket(
  roomId: string | undefined,
  onEvent: (env: SocketEnvelope) => void,
) {
  const token = useAuthStore((s) => s.token)

  useEffect(() => {
    if (!token || !roomId) return
    realtime.connect(token)
    const unsub = realtime.subscribe((env) => {
      if (env.roomId === roomId) onEvent(env)
    })
    realtime.joinRoom(roomId)
    return () => {
      realtime.leaveRoom(roomId)
      unsub()
    }
    // onEvent is kept stable by the caller (useCallback).
  }, [token, roomId, onEvent])
}
