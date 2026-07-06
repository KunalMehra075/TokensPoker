import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import type { ModeDefinition, TaskDetail } from '@/types'

interface Props {
  tasks: TaskDetail[]
  modes: ModeDefinition[]
}

export function PreviousTasks({ tasks, modes }: Props) {
  if (tasks.length === 0) return null
  const modeName = (m: string) => modes.find((d) => d.mode === m)?.name ?? m

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle>Previous estimations</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {tasks.map((t) => (
          <details key={t.task.id} className="group rounded-md border border-border">
            <summary className="flex cursor-pointer list-none items-center justify-between gap-2 p-3">
              <span className="min-w-0">
                <span className="block truncate text-sm font-medium text-foreground">
                  {t.task.title}
                </span>
                <span className="text-xs text-muted-foreground">{modeName(t.task.mode)}</span>
              </span>
              <Badge variant="primary" className="shrink-0">
                {t.final?.finalValue ?? 'closed'}
              </Badge>
            </summary>
            <div className="border-t border-border p-3">
              <div className="grid gap-1.5 sm:grid-cols-2">
                {t.votes.map((v) => (
                  <div key={v.userId} className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">{v.name}</span>
                    <span className="font-mono font-semibold text-foreground">
                      {v.card ?? '-'}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </details>
        ))}
      </CardContent>
    </Card>
  )
}
