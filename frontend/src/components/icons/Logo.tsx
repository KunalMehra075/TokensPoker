import { cn } from '@/lib/utils'

// Brand mark: a stack of estimation cards. Hand-authored SVG, kept in the shared
// icons folder per the house design rules (no emoji, icons live in one place).
export function Logo({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      className={cn('h-6 w-6', className)}
      aria-hidden="true"
    >
      <rect x="3" y="5" width="11" height="15" rx="2" fill="hsl(var(--primary))" opacity="0.25" />
      <rect
        x="7"
        y="3"
        width="11"
        height="15"
        rx="2"
        fill="hsl(var(--primary))"
        stroke="hsl(var(--primary))"
        strokeWidth="1.2"
      />
      <path
        d="M12.5 7.2v6.6M9.6 10.5h5.8"
        stroke="white"
        strokeWidth="1.4"
        strokeLinecap="round"
      />
    </svg>
  )
}

export function Wordmark({ className }: { className?: string }) {
  return (
    <span className={cn('inline-flex items-center gap-2', className)}>
      <Logo />
      <span className="text-[15px] font-bold tracking-tight text-foreground">
        FreeTokens<span className="text-primary">Poker</span>
      </span>
    </span>
  )
}
