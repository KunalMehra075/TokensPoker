import { describe, it, expect } from 'vitest'
import type { MemberVoteView } from '@/types'
import {
  parseCardValue,
  buildDistribution,
  findSpread,
  isUnanimous,
} from './voting'

// Small helper to build vote views tersely.
function vote(name: string, card?: string): MemberVoteView {
  return { userId: name, name, hasVoted: card !== undefined, card }
}

describe('parseCardValue', () => {
  it('parses TOKENS suffixes K/M/B into raw numbers', () => {
    expect(parseCardValue('TOKENS', '500K')).toBe(500_000)
    expect(parseCardValue('TOKENS', '1M')).toBe(1_000_000)
    expect(parseCardValue('TOKENS', '2.5M')).toBe(2_500_000)
    expect(parseCardValue('TOKENS', '1B')).toBe(1_000_000_000)
    expect(parseCardValue('TOKENS', '42')).toBe(42)
  })

  it('is case-insensitive on the unit and tolerates spacing', () => {
    expect(parseCardValue('TOKENS', '5m')).toBe(5_000_000)
    expect(parseCardValue('TOKENS', '10 M')).toBe(10_000_000)
  })

  it('parses COST by stripping $ and commas', () => {
    expect(parseCardValue('COST', '$25')).toBe(25)
    expect(parseCardValue('COST', '$1,250')).toBe(1250)
  })

  it('parses DAYS as plain numbers', () => {
    expect(parseCardValue('DAYS', '8')).toBe(8)
  })

  it('returns null for non-numeric and unsure cards', () => {
    expect(parseCardValue('TOKENS', '?')).toBeNull()
    expect(parseCardValue('TOKENS', '')).toBeNull()
    expect(parseCardValue('MODEL', 'Claude')).toBeNull()
    expect(parseCardValue('TOKENS', 'abc')).toBeNull()
    expect(parseCardValue('COST', '$')).toBeNull()
  })
})

describe('buildDistribution', () => {
  it('counts cards and sorts by frequency descending', () => {
    const dist = buildDistribution([
      vote('a', '5M'),
      vote('b', '1M'),
      vote('c', '5M'),
      vote('d', '5M'),
      vote('e'), // no card, ignored
    ])
    expect(dist).toEqual([
      { card: '5M', count: 3 },
      { card: '1M', count: 1 },
    ])
  })

  it('returns an empty distribution when nobody has a card', () => {
    expect(buildDistribution([vote('a'), vote('b')])).toEqual([])
  })
})

describe('findSpread', () => {
  it('returns the lowest and highest numeric votes', () => {
    const spread = findSpread('TOKENS', [
      vote('low', '500K'),
      vote('mid', '5M'),
      vote('high', '50M'),
    ])
    expect(spread).not.toBeNull()
    expect(spread!.low).toEqual({ name: 'low', card: '500K' })
    expect(spread!.high).toEqual({ name: 'high', card: '50M' })
  })

  it('returns null when fewer than two numeric votes exist', () => {
    expect(findSpread('TOKENS', [vote('a', '5M')])).toBeNull()
    expect(findSpread('TOKENS', [vote('a', '5M'), vote('b', '?')])).toBeNull()
  })

  it('returns null when everyone agrees (no spread to discuss)', () => {
    expect(
      findSpread('TOKENS', [vote('a', '5M'), vote('b', '5M'), vote('c', '5M')]),
    ).toBeNull()
  })

  it('ignores non-numeric outliers when computing the spread', () => {
    const spread = findSpread('COST', [
      vote('a', '$1'),
      vote('b', '$100'),
      vote('c', '?'),
    ])
    expect(spread!.low.card).toBe('$1')
    expect(spread!.high.card).toBe('$100')
  })
})

describe('isUnanimous', () => {
  it('is true when every cast vote matches and there is more than one', () => {
    expect(isUnanimous([vote('a', 'Claude'), vote('b', 'Claude')])).toBe(true)
  })

  it('is false for a single vote or any disagreement', () => {
    expect(isUnanimous([vote('a', 'Claude')])).toBe(false)
    expect(isUnanimous([vote('a', 'Claude'), vote('b', 'GPT')])).toBe(false)
  })

  it('ignores members who have not voted', () => {
    expect(isUnanimous([vote('a', '5M'), vote('b', '5M'), vote('c')])).toBe(true)
  })
})
