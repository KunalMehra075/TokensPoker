import type { HTMLAttributes } from 'react'
import { cva } from 'class-variance-authority'
import type { VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const badgeVariants = cva(
  'inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium',
  {
    variants: {
      variant: {
        neutral: 'bg-muted text-muted-foreground',
        primary: 'bg-blue-50 text-primary',
        success: 'bg-emerald-50 text-success',
        warning: 'bg-amber-50 text-amber-700',
        danger: 'bg-red-50 text-danger',
        outline: 'border border-border text-foreground',
      },
    },
    defaultVariants: { variant: 'neutral' },
  },
)

export interface BadgeProps
  extends HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />
}
