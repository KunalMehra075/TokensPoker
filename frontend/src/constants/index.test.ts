import { describe, it, expect } from 'vitest'
import { BRAND, SOCKET_EVENTS, inviteLink } from './index'

describe('inviteLink', () => {
  it('builds an absolute /join/<code> link from the current origin', () => {
    const link = inviteLink('ABC123')
    expect(link).toBe(`${window.location.origin}/join/ABC123`)
    expect(link).toMatch(/\/join\/ABC123$/)
  })
})

describe('BRAND', () => {
  it('points at the freetokenspoker.com domain and contact inbox', () => {
    expect(BRAND.domain).toBe('freetokenspoker.com')
    expect(BRAND.url).toBe('https://freetokenspoker.com')
    expect(BRAND.contactEmail).toBe('hello@freetokenspoker.com')
  })
})

describe('SOCKET_EVENTS', () => {
  it('matches the backend wire contract in internal/realtime/events.go', () => {
    // If these drift from the Go constants, realtime updates silently break.
    expect(SOCKET_EVENTS).toEqual({
      presence: 'presence',
      memberJoined: 'member_joined',
      memberLeft: 'member_left',
      taskCreated: 'task_created',
      voteReceived: 'vote_received',
      votesRevealed: 'votes_revealed',
      taskClosed: 'task_closed',
      finalDecision: 'final_decision',
    })
  })
})
