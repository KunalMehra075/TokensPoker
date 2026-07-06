import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Dialog } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { createTask } from '@/services/api'
import type { ApiErrorShape } from '@/services/api'
import type { EstimationMode, ModeDefinition } from '@/types'

interface Props {
  open: boolean
  onClose: () => void
  roomId: string
  modes: ModeDefinition[]
}

export function CreateTaskDialog({ open, onClose, roomId, modes }: Props) {
  const qc = useQueryClient()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [mode, setMode] = useState<EstimationMode>('TOKENS')

  const reset = () => {
    setTitle('')
    setDescription('')
    setMode('TOKENS')
  }

  const create = useMutation({
    mutationFn: () => createTask({ roomId, title: title.trim(), description: description.trim(), mode }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['roomTasks', roomId] })
      toast.success('Estimation started')
      reset()
      onClose()
    },
    onError: (e: ApiErrorShape) => toast.error(e.message),
  })

  const selected = modes.find((m) => m.mode === mode)

  return (
    <Dialog open={open} onClose={onClose} title="Start an estimation" description="Everyone in the room votes on this.">
      <form
        onSubmit={(e) => {
          e.preventDefault()
          if (title.trim()) create.mutate()
        }}
      >
        <div className="mb-4">
          <Label htmlFor="title">What are you estimating?</Label>
          <Input
            id="title"
            placeholder="Add semantic search to the docs site"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            autoFocus
            maxLength={160}
          />
        </div>

        <div className="mb-4">
          <Label htmlFor="desc">Context (optional)</Label>
          <Textarea
            id="desc"
            placeholder="Anything the team should know before voting"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            maxLength={2000}
          />
        </div>

        <div className="mb-5">
          <Label>Estimation mode</Label>
          <div className="grid grid-cols-2 gap-2">
            {modes.map((m) => (
              <button
                type="button"
                key={m.mode}
                onClick={() => setMode(m.mode)}
                className={`rounded-md border p-3 text-left transition-colors ${
                  mode === m.mode
                    ? 'border-primary bg-blue-50'
                    : 'border-border bg-white hover:bg-muted'
                }`}
              >
                <span className="block text-sm font-semibold text-foreground">{m.name}</span>
                <span className="mt-0.5 block text-xs text-muted-foreground">
                  {m.cards.slice(0, 4).join(', ')}...
                </span>
              </button>
            ))}
          </div>
          {selected && (
            <p className="mt-2 text-xs text-muted-foreground">{selected.description}</p>
          )}
        </div>

        <div className="flex justify-end gap-2">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" loading={create.isPending} disabled={!title.trim()}>
            Start voting
          </Button>
        </div>
      </form>
    </Dialog>
  )
}
