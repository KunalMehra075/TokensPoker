import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { Inbox } from 'lucide-react'
import { AppHeader } from '@/components/layout/AppHeader'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { fetchHistory } from '@/services/api'
import { useModes } from '@/hooks/useModes'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'

export function HistoryPage() {
  useDocumentTitle('History')
  const { data, isLoading } = useQuery({ queryKey: ['history'], queryFn: fetchHistory })
  const { data: modes } = useModes()
  const modeName = (m: string) => modes?.find((d) => d.mode === m)?.name ?? m

  return (
    <div className="min-h-screen bg-muted/40">
      <AppHeader />
      <main className="container-app py-10">
        <div className="mx-auto max-w-3xl">
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Estimation history</h1>
          <p className="mt-2 text-muted-foreground">
            Every round you have closed, across all your rooms. Proof of what the team thought,
            before reality had opinions.
          </p>

          <div className="mt-8 space-y-3">
            {isLoading && (
              <>
                <Skeleton className="h-24 w-full" />
                <Skeleton className="h-24 w-full" />
              </>
            )}

            {!isLoading && data && data.length === 0 && (
              <Card>
                <CardContent className="flex flex-col items-center py-16 text-center">
                  <span className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-muted text-muted-foreground">
                    <Inbox className="h-6 w-6" />
                  </span>
                  <h2 className="mt-4 text-lg font-semibold text-foreground">Nothing here yet</h2>
                  <p className="mt-2 max-w-sm text-sm text-muted-foreground">
                    Finish your first estimation round and it will show up here for posterity.
                  </p>
                  <Link to="/app" className="mt-6">
                    <Button>Go to dashboard</Button>
                  </Link>
                </CardContent>
              </Card>
            )}

            {!isLoading &&
              data?.map((item) => (
                <Card key={item.task.id}>
                  <CardContent className="p-5">
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div className="min-w-0">
                        <h3 className="font-semibold text-foreground">{item.task.title}</h3>
                        <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                          <Badge variant="neutral">{modeName(item.task.mode)}</Badge>
                          <span>Room: {item.roomName}</span>
                          <span>{new Date(item.task.createdAt).toLocaleDateString()}</span>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-xs text-muted-foreground">Final</p>
                        <p className="text-lg font-bold text-primary">
                          {item.final?.finalValue ?? '-'}
                        </p>
                      </div>
                    </div>
                    {item.votes.length > 0 && (
                      <div className="mt-4 flex flex-wrap gap-1.5">
                        {item.votes.map((v) => (
                          <span
                            key={v.id}
                            className="inline-flex items-center gap-1 rounded-md border border-border bg-muted/40 px-2 py-1 text-xs"
                          >
                            <span className="text-muted-foreground">{v.userName}</span>
                            <span className="font-mono font-semibold text-foreground">
                              {v.selectedCard}
                            </span>
                          </span>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
          </div>
        </div>
      </main>
    </div>
  )
}
