import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Plus, LogIn, History } from 'lucide-react'
import { AppHeader } from '@/components/layout/AppHeader'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { createRoom, joinRoom } from '@/services/api'
import type { ApiErrorShape } from '@/services/api'
import { useAuthStore } from '@/store/auth'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'

export function DashboardPage() {
  useDocumentTitle('Dashboard')
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const [roomName, setRoomName] = useState('')
  const [code, setCode] = useState('')

  const create = useMutation({
    mutationFn: () => createRoom(roomName.trim()),
    onSuccess: (room) => {
      toast.success('Room created')
      navigate(`/room/${room.id}`)
    },
    onError: (e: ApiErrorShape) => toast.error(e.message),
  })

  const join = useMutation({
    mutationFn: () => joinRoom(code.trim().toUpperCase()),
    onSuccess: (room) => {
      toast.success(`Joined ${room.name}`)
      navigate(`/room/${room.id}`)
    },
    onError: (e: ApiErrorShape) => toast.error(e.message),
  })

  return (
    <div className="min-h-screen bg-muted/40">
      <AppHeader />
      <main className="container-app py-12">
        <div className="mx-auto max-w-3xl">
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            {user ? `Hey ${user.name.split(' ')[0]}, ready to estimate?` : 'Ready to estimate?'}
          </h1>
          <p className="mt-2 text-muted-foreground">
            Spin up a fresh room for your sprint, or drop in with a code a teammate shared.
          </p>

          <div className="mt-8 grid gap-5 sm:grid-cols-2">
            <Card>
              <CardHeader>
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-md bg-blue-50 text-primary">
                  <Plus className="h-5 w-5" />
                </span>
                <CardTitle className="mt-2">Create a room</CardTitle>
                <CardDescription>You become the owner and run the reveals.</CardDescription>
              </CardHeader>
              <CardContent>
                <form
                  onSubmit={(e) => {
                    e.preventDefault()
                    if (roomName.trim()) create.mutate()
                  }}
                >
                  <Label htmlFor="roomName">Room name</Label>
                  <Input
                    id="roomName"
                    placeholder="Sprint 42 planning"
                    value={roomName}
                    onChange={(e) => setRoomName(e.target.value)}
                    maxLength={80}
                  />
                  <Button
                    type="submit"
                    className="mt-4 w-full"
                    loading={create.isPending}
                    disabled={!roomName.trim()}
                  >
                    Create room
                  </Button>
                </form>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-md bg-emerald-50 text-success">
                  <LogIn className="h-5 w-5" />
                </span>
                <CardTitle className="mt-2">Join a room</CardTitle>
                <CardDescription>Enter the six character code from your team.</CardDescription>
              </CardHeader>
              <CardContent>
                <form
                  onSubmit={(e) => {
                    e.preventDefault()
                    if (code.trim()) join.mutate()
                  }}
                >
                  <Label htmlFor="code">Room code</Label>
                  <Input
                    id="code"
                    placeholder="ABC123"
                    value={code}
                    onChange={(e) => setCode(e.target.value.toUpperCase())}
                    maxLength={12}
                    className="font-mono tracking-widest"
                  />
                  <Button
                    type="submit"
                    variant="secondary"
                    className="mt-4 w-full"
                    loading={join.isPending}
                    disabled={!code.trim()}
                  >
                    Join room
                  </Button>
                </form>
              </CardContent>
            </Card>
          </div>

          <button
            onClick={() => navigate('/history')}
            className="mt-8 inline-flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-foreground"
          >
            <History className="h-4 w-4" />
            View your estimation history
          </button>
        </div>
      </main>
    </div>
  )
}
