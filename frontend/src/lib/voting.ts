import type { EstimationMode, MemberVoteView } from '@/types'

// parseCardValue turns a card label into a comparable number, or null when the
// card is non-numeric ("?", model names). Used for spread + average insights.
export function parseCardValue(mode: EstimationMode, card: string): number | null {
  if (mode === 'MODEL') return null
  const c = card.trim()
  if (c === '?' || c === '') return null
  if (mode === 'TOKENS') {
    const m = c.match(/^([\d.]+)\s*([KMB])?$/i)
    if (!m) return null
    const n = parseFloat(m[1])
    const unit = (m[2] || '').toUpperCase()
    const mult = unit === 'B' ? 1e9 : unit === 'M' ? 1e6 : unit === 'K' ? 1e3 : 1
    return n * mult
  }
  if (mode === 'COST') {
    const n = parseFloat(c.replace(/[$,]/g, ''))
    return Number.isNaN(n) ? null : n
  }
  // DAYS
  const n = parseFloat(c)
  return Number.isNaN(n) ? null : n
}

export interface VoteDistribution {
  card: string
  count: number
}

export function buildDistribution(votes: MemberVoteView[]): VoteDistribution[] {
  const counts = new Map<string, number>()
  for (const v of votes) {
    if (v.card) counts.set(v.card, (counts.get(v.card) ?? 0) + 1)
  }
  return [...counts.entries()]
    .map(([card, count]) => ({ card, count }))
    .sort((a, b) => b.count - a.count)
}

export interface Spread {
  low: { name: string; card: string }
  high: { name: string; card: string }
}

// outliers finds the lowest and highest numeric votes, the pair worth discussing.
export function findSpread(mode: EstimationMode, votes: MemberVoteView[]): Spread | null {
  const numeric = votes
    .filter((v) => v.card)
    .map((v) => ({ name: v.name, card: v.card as string, value: parseCardValue(mode, v.card as string) }))
    .filter((v): v is { name: string; card: string; value: number } => v.value !== null)
  if (numeric.length < 2) return null
  let low = numeric[0]
  let high = numeric[0]
  for (const v of numeric) {
    if (v.value < low.value) low = v
    if (v.value > high.value) high = v
  }
  if (low.card === high.card) return null
  return { low: { name: low.name, card: low.card }, high: { name: high.name, card: high.card } }
}

export function isUnanimous(votes: MemberVoteView[]): boolean {
  const cards = votes.filter((v) => v.card).map((v) => v.card)
  return cards.length > 1 && cards.every((c) => c === cards[0])
}
