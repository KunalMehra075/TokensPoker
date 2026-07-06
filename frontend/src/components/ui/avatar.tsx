import { cn, colorForId, initials } from '@/lib/utils'

interface AvatarProps {
  id: string
  name: string
  size?: 'sm' | 'md'
  online?: boolean
  className?: string
}

// A deterministic, label-free initials avatar. No images, no emoji.
export function Avatar({ id, name, size = 'md', online, className }: AvatarProps) {
  const dims = size === 'sm' ? 'h-7 w-7 text-xs' : 'h-9 w-9 text-sm'
  return (
    <span className={cn('relative inline-flex shrink-0', className)}>
      <span
        className={cn(
          'inline-flex items-center justify-center rounded-full font-semibold',
          dims,
          colorForId(id),
        )}
        aria-hidden="true"
      >
        {initials(name)}
      </span>
      {online !== undefined && (
        <span
          className={cn(
            'absolute -bottom-0 -right-0 h-2.5 w-2.5 rounded-full border-2 border-white',
            online ? 'bg-success' : 'bg-gray-300',
          )}
          title={online ? 'online' : 'offline'}
        />
      )}
    </span>
  )
}
