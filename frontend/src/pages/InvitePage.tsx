import { useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ArrowRight, Loader2, Users } from 'lucide-react'
import { Wordmark } from '@/components/icons/Logo'
import { Button } from '@/components/ui/button'
import { fetchRoomPreview, joinRoom } from '@/services/api'
import type { ApiErrorShape } from '@/services/api'
import { useAuthStore } from '@/store/auth'
import { PENDING_ROOM_KEY } from '@/constants'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'

// InvitePage backs the shareable /join/:code link. It shows what room you are
// joining, then either routes you through sign-in (stashing the code) or, once
// authenticated, asks you to confirm the join.
export function InvitePage() {
  const { code } = useParams<{ code: string }>()
  const navigate = useNavigate()
  const token = useAuthStore((s) => s.token)
  const user = useAuthStore((s) => s.user)

  useDocumentTitle('Join a room')

  const preview = useQuery({
    queryKey: ['invite', code],
    queryFn: () => fetchRoomPreview(code!),
    enabled: !!code,
    retry: false,
  })

  // Remember the code while the invitee signs in; clear it once authenticated so
  // it cannot leak into a later, unrelated login.
  useEffect(() => {
    if (!code) return
    if (token) {
      sessionStorage.removeItem(PENDING_ROOM_KEY)
    } else {
      sessionStorage.setItem(PENDING_ROOM_KEY, code)
    }
  }, [code, token])

  const join = useMutation({
    mutationFn: () => joinRoom(code!),
    onSuccess: (room) => {
      sessionStorage.removeItem(PENDING_ROOM_KEY)
      toast.success(`Joined ${room.name}`)
      navigate(`/room/${room.id}`, { replace: true })
    },
    onError: (e: ApiErrorShape) => toast.error(e.message),
  })

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-muted/40 px-4">
      <div className="w-full max-w-md rounded-lg border border-border bg-white p-8 shadow-card">
        <div className="flex justify-center">
          <Wordmark />
        </div>

        {preview.isLoading && (
          <div className="mt-8 flex flex-col items-center text-center">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            <p className="mt-3 text-sm text-muted-foreground">Looking up the room...</p>
          </div>
        )}

        {preview.isError && (
          <div className="mt-8 text-center">
            <h1 className="text-lg font-semibold text-foreground">This invite is not valid</h1>
            <p className="mt-2 text-sm text-muted-foreground">
              The link may be mistyped, or the room was deleted. You can still create your own or
              join with a different code.
            </p>
            <div className="mt-6 flex justify-center gap-2">
              <Button onClick={() => navigate(token ? '/app' : '/login')}>
                {token ? 'Go to dashboard' : 'Sign in'}
              </Button>
            </div>
          </div>
        )}

        {preview.data && (
          <div className="mt-8 text-center">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-blue-50 px-3 py-1 text-xs font-medium text-primary">
              <Users className="h-3.5 w-3.5" />
              {preview.data.memberCount}{' '}
              {preview.data.memberCount === 1 ? 'member' : 'members'}
            </span>
            <h1 className="mt-4 text-xl font-bold tracking-tight text-foreground">
              You are invited to
            </h1>
            <p className="mt-1 text-2xl font-bold text-primary">{preview.data.name}</p>

            {token ? (
              <>
                <p className="mt-4 text-sm text-muted-foreground">
                  Signed in as {user?.name}. Want to join this estimation room?
                </p>
                <div className="mt-6 flex flex-col gap-2">
                  <Button onClick={() => join.mutate()} loading={join.isPending}>
                    Join {preview.data.name}
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" onClick={() => navigate('/app')} disabled={join.isPending}>
                    Maybe later
                  </Button>
                </div>
              </>
            ) : (
              <>
                <p className="mt-4 text-sm text-muted-foreground">
                  Sign in with your name and email to join. It takes about ten seconds, no
                  password required.
                </p>
                <div className="mt-6">
                  <Button className="w-full" onClick={() => navigate('/login')}>
                    Sign in to join
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
