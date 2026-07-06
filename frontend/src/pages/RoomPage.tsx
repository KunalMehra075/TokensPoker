import { useCallback, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Plus } from 'lucide-react'
import { AppHeader } from '@/components/layout/AppHeader'
import { RoomHeader } from '@/components/room/RoomHeader'
import { MembersPanel } from '@/components/room/MembersPanel'
import { VotingArea } from '@/components/voting/VotingArea'
import { CreateTaskDialog } from '@/components/task/CreateTaskDialog'
import { PreviousTasks } from '@/components/task/PreviousTasks'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { deleteRoom, fetchRoom, fetchRoomTasks } from '@/services/api'
import { useModes, findMode } from '@/hooks/useModes'
import { useRoomSocket } from '@/hooks/useRoomSocket'
import { useAuthStore } from '@/store/auth'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'
import { SOCKET_EVENTS } from '@/constants'
import type { Presence, SocketEnvelope, TaskDetail } from '@/types'

export function RoomPage() {
  const { id } = useParams<{ id: string }>()
  const qc = useQueryClient()
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const [presence, setPresence] = useState<Presence[]>([])
  const [createOpen, setCreateOpen] = useState(false)

  const roomQuery = useQuery({ queryKey: ['room', id], queryFn: () => fetchRoom(id!), enabled: !!id })
  const tasksQuery = useQuery({
    queryKey: ['roomTasks', id],
    queryFn: () => fetchRoomTasks(id!),
    enabled: !!id,
  })
  const modesQuery = useModes()

  useDocumentTitle(roomQuery.data ? roomQuery.data.name : 'Room')

  // Route realtime events to cache invalidation + presence state.
  const onEvent = useCallback(
    (env: SocketEnvelope) => {
      switch (env.event) {
        case SOCKET_EVENTS.presence:
          setPresence(env.payload as Presence[])
          break
        case SOCKET_EVENTS.memberJoined:
          qc.invalidateQueries({ queryKey: ['room', id] })
          break
        case SOCKET_EVENTS.finalDecision:
        case SOCKET_EVENTS.taskClosed:
          qc.invalidateQueries({ queryKey: ['roomTasks', id] })
          qc.invalidateQueries({ queryKey: ['history'] })
          break
        default:
          // task_created, vote_received, votes_revealed
          qc.invalidateQueries({ queryKey: ['roomTasks', id] })
      }
    },
    [qc, id],
  )

  useRoomSocket(id, onEvent)

  const remove = useMutation({
    mutationFn: () => deleteRoom(id!),
    onSuccess: () => {
      toast.success('Room deleted')
      navigate('/app')
    },
    onError: () => toast.error('Could not delete room'),
  })

  if (roomQuery.isLoading) {
    return (
      <div className="min-h-screen bg-muted/40">
        <AppHeader />
        <main className="container-app py-8">
          <Skeleton className="h-8 w-64" />
          <Skeleton className="mt-6 h-64 w-full" />
        </main>
      </div>
    )
  }

  if (roomQuery.isError || !roomQuery.data || !user) {
    return (
      <div className="min-h-screen bg-muted/40">
        <AppHeader />
        <main className="container-app py-16 text-center">
          <h1 className="text-xl font-semibold text-foreground">Room not found</h1>
          <p className="mt-2 text-muted-foreground">
            It may have been deleted, or you are not a member.
          </p>
          <Button className="mt-6" onClick={() => navigate('/app')}>
            Back to dashboard
          </Button>
        </main>
      </div>
    )
  }

  const room = roomQuery.data
  const tasks = tasksQuery.data ?? []
  const modes = modesQuery.data ?? []
  const isOwner = room.ownerId === user.id

  const activeTask: TaskDetail | undefined = tasks.find((t) => t.task.status !== 'CLOSED')
  const previous = tasks.filter((t) => t.task.status === 'CLOSED')

  return (
    <div className="min-h-screen bg-muted/40">
      <AppHeader />
      <main className="container-app py-8">
        <RoomHeader room={room} isOwner={isOwner} onDelete={() => remove.mutate()} />

        <div className="mt-6 grid gap-6 lg:grid-cols-[1fr_300px]">
          <div className="space-y-6">
            {activeTask ? (
              <VotingArea
                detail={activeTask}
                modeDef={findMode(modes, activeTask.task.mode)}
                isOwner={isOwner}
                currentUserId={user.id}
                roomId={room.id}
              />
            ) : (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-16 text-center">
                  <h2 className="text-lg font-semibold text-foreground">No active estimation</h2>
                  <p className="mt-2 max-w-sm text-sm text-muted-foreground">
                    {isOwner
                      ? 'Start a round and the whole room gets cards to vote with.'
                      : 'Hang tight. The room owner will start the next estimation soon.'}
                  </p>
                  {isOwner && (
                    <Button className="mt-6" onClick={() => setCreateOpen(true)}>
                      <Plus className="h-4 w-4" />
                      Start an estimation
                    </Button>
                  )}
                </CardContent>
              </Card>
            )}

            {isOwner && activeTask && (
              <p className="text-center text-xs text-muted-foreground">
                Finish this round (reveal, then pick a final value) before starting the next one.
              </p>
            )}

            <PreviousTasks tasks={previous} modes={modes} />
          </div>

          <aside className="space-y-6">
            <MembersPanel members={room.members} ownerId={room.ownerId} online={presence} />
            {isOwner && !activeTask && (
              <Button className="w-full" onClick={() => setCreateOpen(true)}>
                <Plus className="h-4 w-4" />
                New estimation
              </Button>
            )}
          </aside>
        </div>
      </main>

      <CreateTaskDialog
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        roomId={room.id}
        modes={modes}
      />
    </div>
  )
}
