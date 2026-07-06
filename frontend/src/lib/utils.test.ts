import { describe, it, expect } from 'vitest'
import { cn, colorForId, initials } from './utils'

describe('cn', () => {
  it('merges class names and resolves Tailwind conflicts (last wins)', () => {
    const isHidden: boolean = false
    expect(cn('px-2', 'px-4')).toBe('px-4')
    expect(cn('text-sm', isHidden && 'hidden', 'font-medium')).toBe('text-sm font-medium')
  })
})

describe('colorForId', () => {
  it('is deterministic for the same id', () => {
    expect(colorForId('user-123')).toBe(colorForId('user-123'))
  })

  it('always returns one of the known palette classes', () => {
    const palette = new Set([
      'bg-blue-100 text-blue-700',
      'bg-emerald-100 text-emerald-700',
      'bg-amber-100 text-amber-700',
      'bg-violet-100 text-violet-700',
      'bg-rose-100 text-rose-700',
      'bg-cyan-100 text-cyan-700',
    ])
    for (const id of ['a', 'bbb', '507f1f77bcf86cd799439011', '', 'Zzz']) {
      expect(palette.has(colorForId(id))).toBe(true)
    }
  })
})

describe('initials', () => {
  it('takes first + last initial for multi-word names', () => {
    expect(initials('Ada Lovelace')).toBe('AL')
    expect(initials('  grace   brewster  hopper ')).toBe('GH')
  })

  it('takes the first two letters of a single-word name', () => {
    expect(initials('Cher')).toBe('CH')
    expect(initials('x')).toBe('X')
  })

  it('falls back to ? for an empty name', () => {
    expect(initials('')).toBe('?')
    expect(initials('   ')).toBe('?')
  })
})
