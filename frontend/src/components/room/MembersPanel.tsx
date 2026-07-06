import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Avatar } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import type { Member, Presence } from '@/types'

interface Props {
  members: Member[]
  ownerId: string
  online: Presence[]
}

export function MembersPanel({ members, ownerId, online }: Props) {
  const onlineIds = new Set(online.map((p) => p.userId))
  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center justify-between">
          <span>Members</span>
          <Badge variant="neutral">{members.length}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {members.map((m) => (
          <div key={m.userId} className="flex items-center gap-3">
            <Avatar id={m.userId} name={m.name} size="sm" online={onlineIds.has(m.userId)} />
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium text-foreground">{m.name}</p>
              <p className="truncate text-xs text-muted-foreground">{m.email}</p>
            </div>
            {m.userId === ownerId && (
              <Badge variant="primary" className="shrink-0">
                Owner
              </Badge>
            )}
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
