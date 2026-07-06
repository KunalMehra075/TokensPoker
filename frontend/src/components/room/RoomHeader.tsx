import { useState } from 'react'
import { Check, Copy, Link2, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { inviteLink } from '@/constants'
import type { Room } from '@/types'

interface Props {
  room: Room
  isOwner: boolean
  onDelete: () => void
}

export function RoomHeader({ room, isOwner, onDelete }: Props) {
  const [copied, setCopied] = useState(false)

  const copy = async () => {
    await navigator.clipboard.writeText(room.roomCode)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  const copyLink = async () => {
    await navigator.clipboard.writeText(inviteLink(room.roomCode))
    toast.success('Invite link copied. Send it to your team.')
  }

  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 className="text-xl font-bold tracking-tight text-foreground">{room.name}</h1>
        <div className="mt-1 flex flex-wrap items-center gap-2">
          <span className="text-sm text-muted-foreground">Invite your team:</span>
          <button
            onClick={copy}
            className="inline-flex items-center gap-1.5 rounded-md border border-border bg-white px-2.5 py-1 font-mono text-sm font-semibold tracking-widest text-foreground hover:bg-muted"
            aria-label="Copy room code"
          >
            {room.roomCode}
            {copied ? (
              <Check className="h-3.5 w-3.5 text-success" />
            ) : (
              <Copy className="h-3.5 w-3.5 text-muted-foreground" />
            )}
          </button>
          <Button variant="secondary" size="sm" onClick={copyLink}>
            <Link2 className="h-3.5 w-3.5" />
            Copy invite link
          </Button>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Badge variant={isOwner ? 'primary' : 'neutral'}>{isOwner ? 'Owner' : 'Member'}</Badge>
        {isOwner && (
          <Button variant="ghost" size="sm" onClick={onDelete}>
            <Trash2 className="h-4 w-4" />
            <span className="hidden sm:inline">Delete room</span>
          </Button>
        )}
      </div>
    </div>
  )
}
