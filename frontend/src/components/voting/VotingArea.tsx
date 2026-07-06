import { useEffect, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Check, Eye, Loader2, Trophy } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Avatar } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'
import { buildDistribution, findSpread, isUnanimous } from '@/lib/voting'
import { finalizeTask, revealTask, submitVote } from '@/services/api'
import type { ApiErrorShape } from '@/services/api'
import type { ModeDefinition, TaskDetail } from '@/types'

interface Props {
  detail: TaskDetail
  modeDef?: ModeDefinition
  isOwner: boolean
  currentUserId: string
  roomId: string
}

export function VotingArea({ detail, modeDef, isOwner, currentUserId, roomId }: Props) {
  const qc = useQueryClient()
  const { task, votes } = detail
  const revealed = task.revealed
  const [myCard, setMyCard] = useState<string | null>(null)
  const [finalChoice, setFinalChoice] = useState('')

  // The server hides cards before reveal, even our own, so we remember our pick
  // locally. Reset when the task changes.
  useEffect(() => {
    setMyCard(null)
    setFinalChoice('')
  }, [task.id])

  const invalidate = () => qc.invalidateQueries({ queryKey: ['roomTasks', roomId] })

  const vote = useMutation({
    mutationFn: (card: string) => submitVote(task.id, card),
    onMutate: (card) => setMyCard(card),
    onError: (e: ApiErrorShape) => toast.error(e.message),
    onSuccess: () => invalidate(),
  })

  const reveal = useMutation({
    mutationFn: () => revealTask(task.id),
    onError: (e: ApiErrorShape) => toast.error(e.message),
    onSuccess: () => invalidate(),
  })

  const finalize = useMutation({
    mutationFn: (value: string) => finalizeTask(task.id, value),
    onError: (e: ApiErrorShape) => toast.error(e.message),
    onSuccess: () => {
      invalidate()
      qc.invalidateQueries({ queryKey: ['history'] })
      toast.success('Final decision saved')
    },
  })

  const cards = modeDef?.cards ?? []
  const distribution = buildDistribution(votes)
  const spread = findSpread(task.mode, votes)
  const unanimous = isUnanimous(votes)
  const maxCount = distribution[0]?.count ?? 0

  return (
    <Card>
      <CardHeader className="border-b border-border">
        <div className="flex flex-wrap items-start justify-between gap-2">
          <div>
            <div className="flex items-center gap-2">
              <Badge variant="primary">{modeDef?.name ?? task.mode}</Badge>
              {revealed ? (
                <Badge variant="success">Revealed</Badge>
              ) : (
                <Badge variant="warning">Voting open</Badge>
              )}
            </div>
            <h2 className="mt-2 text-lg font-semibold text-foreground">{task.title}</h2>
            {task.description && (
              <p className="mt-1 text-sm text-muted-foreground">{task.description}</p>
            )}
          </div>
          <div className="text-right">
            <p className="text-2xl font-bold text-foreground">
              {detail.voteCount}
              <span className="text-base font-normal text-muted-foreground">
                /{detail.memberCount}
              </span>
            </p>
            <p className="text-xs text-muted-foreground">voted</p>
          </div>
        </div>
      </CardHeader>

      <CardContent className="pt-5">
        {/* Voting deck (pre-reveal) */}
        {!revealed && (
          <>
            <p className="mb-3 text-sm font-medium text-foreground">
              {myCard ? 'Your pick (change it any time before reveal):' : 'Pick your card:'}
            </p>
            <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 lg:grid-cols-8">
              {cards.map((card) => {
                const active = myCard === card
                return (
                  <button
                    key={card}
                    onClick={() => vote.mutate(card)}
                    disabled={vote.isPending}
                    className={cn(
                      'flex h-20 items-center justify-center rounded-md border text-base font-bold transition-all',
                      active
                        ? '-translate-y-1 border-primary bg-primary text-primary-foreground shadow-card-hover'
                        : 'border-border bg-white text-foreground hover:-translate-y-0.5 hover:border-primary hover:shadow-card',
                    )}
                    aria-pressed={active}
                  >
                    {card}
                  </button>
                )
              })}
            </div>
          </>
        )}

        {/* Reveal results */}
        {revealed && (
          <div className="space-y-5">
            {unanimous && (
              <div className="flex items-center gap-2 rounded-md bg-emerald-50 px-4 py-3 text-sm font-medium text-success">
                <Trophy className="h-4 w-4" />
                Unanimous. That basically never happens. Bank it.
              </div>
            )}
            {spread && (
              <div className="rounded-md border border-border bg-muted/40 px-4 py-3 text-sm">
                <span className="font-medium text-foreground">Biggest gap to discuss: </span>
                <span className="text-muted-foreground">
                  {spread.low.name} said {spread.low.card}, {spread.high.name} said{' '}
                  {spread.high.card}.
                </span>
              </div>
            )}

            {/* Distribution */}
            <div className="space-y-2">
              {distribution.map((d) => (
                <div key={d.card} className="flex items-center gap-3">
                  <span className="w-16 shrink-0 text-right font-mono text-sm font-semibold text-foreground">
                    {d.card}
                  </span>
                  <div className="h-7 flex-1 overflow-hidden rounded bg-muted">
                    <div
                      className="flex h-full items-center bg-primary/80 px-2 text-xs font-medium text-white"
                      style={{ width: `${Math.max((d.count / maxCount) * 100, 12)}%` }}
                    >
                      {d.count}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Per-member status / cards */}
        <div className="mt-6">
          <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
            {revealed ? 'Everyone showed' : 'Voting status'}
          </p>
          <div className="grid gap-2 sm:grid-cols-2">
            {votes.map((v) => (
              <div
                key={v.userId}
                className="flex items-center justify-between rounded-md border border-border bg-white px-3 py-2"
              >
                <span className="flex items-center gap-2">
                  <Avatar id={v.userId} name={v.name} size="sm" />
                  <span className="text-sm font-medium text-foreground">
                    {v.name}
                    {v.userId === currentUserId && (
                      <span className="text-muted-foreground"> (you)</span>
                    )}
                  </span>
                </span>
                {revealed ? (
                  <span className="font-mono text-sm font-bold text-foreground">
                    {v.card ?? 'no vote'}
                  </span>
                ) : v.hasVoted ? (
                  <Badge variant="success">
                    <Check className="h-3 w-3" />
                    in
                  </Badge>
                ) : (
                  <Badge variant="neutral">
                    <Loader2 className="h-3 w-3 animate-spin" />
                    thinking
                  </Badge>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Owner controls */}
        {isOwner && !revealed && (
          <div className="mt-6 flex justify-end">
            <Button onClick={() => reveal.mutate()} loading={reveal.isPending} disabled={detail.voteCount === 0}>
              <Eye className="h-4 w-4" />
              Reveal votes
            </Button>
          </div>
        )}

        {isOwner && revealed && !detail.final && (
          <div className="mt-6 rounded-md border border-border bg-muted/40 p-4">
            <p className="mb-2 text-sm font-medium text-foreground">Lock in the final decision</p>
            <div className="mb-3 flex flex-wrap gap-2">
              {cards.map((card) => (
                <button
                  key={card}
                  onClick={() => setFinalChoice(card)}
                  className={cn(
                    'rounded-md border px-3 py-1.5 text-sm font-semibold',
                    finalChoice === card
                      ? 'border-primary bg-primary text-primary-foreground'
                      : 'border-border bg-white text-foreground hover:bg-muted',
                  )}
                >
                  {card}
                </button>
              ))}
            </div>
            <div className="flex gap-2">
              <Input
                placeholder="or type a custom value"
                value={finalChoice}
                onChange={(e) => setFinalChoice(e.target.value)}
                maxLength={80}
              />
              <Button
                onClick={() => finalChoice.trim() && finalize.mutate(finalChoice.trim())}
                loading={finalize.isPending}
                disabled={!finalChoice.trim()}
              >
                Save
              </Button>
            </div>
          </div>
        )}

        {detail.final && (
          <div className="mt-6 flex items-center gap-2 rounded-md bg-blue-50 px-4 py-3">
            <Trophy className="h-4 w-4 text-primary" />
            <span className="text-sm text-foreground">
              Final decision: <span className="font-bold">{detail.final.finalValue}</span>
            </span>
          </div>
        )}

        {!isOwner && !revealed && (
          <p className="mt-6 text-center text-sm text-muted-foreground">
            Waiting for the room owner to reveal once everyone has voted.
          </p>
        )}
      </CardContent>
    </Card>
  )
}
