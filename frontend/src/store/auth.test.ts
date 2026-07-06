import { describe, it, expect, beforeEach } from 'vitest'
import type { User } from '@/types'
import { useAuthStore, getToken } from './auth'

const user: User = {
  id: 'u1',
  email: 'dev@example.com',
  name: 'Dev',
  createdAt: '2026-01-01T00:00:00Z',
  lastLogin: '2026-01-01T00:00:00Z',
}

beforeEach(() => {
  // Reset the singleton store between cases.
  useAuthStore.getState().logout()
})

describe('useAuthStore', () => {
  it('starts logged out', () => {
    const s = useAuthStore.getState()
    expect(s.token).toBeNull()
    expect(s.user).toBeNull()
    expect(s.isAuthenticated()).toBe(false)
  })

  it('setAuth stores the token and user and flips isAuthenticated', () => {
    useAuthStore.getState().setAuth('jwt-abc', user)
    const s = useAuthStore.getState()
    expect(s.token).toBe('jwt-abc')
    expect(s.user).toEqual(user)
    expect(s.isAuthenticated()).toBe(true)
  })

  it('setUser updates the profile without touching the token', () => {
    useAuthStore.getState().setAuth('jwt-abc', user)
    useAuthStore.getState().setUser({ ...user, name: 'Renamed' })
    const s = useAuthStore.getState()
    expect(s.token).toBe('jwt-abc')
    expect(s.user?.name).toBe('Renamed')
  })

  it('logout clears token and user', () => {
    useAuthStore.getState().setAuth('jwt-abc', user)
    useAuthStore.getState().logout()
    const s = useAuthStore.getState()
    expect(s.token).toBeNull()
    expect(s.user).toBeNull()
    expect(s.isAuthenticated()).toBe(false)
  })

  it('persists the token to localStorage so a refresh keeps the session', () => {
    useAuthStore.getState().setAuth('jwt-persist', user)
    const raw = localStorage.getItem('ftp.auth')
    expect(raw).toBeTruthy()
    expect(raw).toContain('jwt-persist')
  })
})

describe('getToken', () => {
  it('mirrors the store token for the axios interceptor', () => {
    expect(getToken()).toBeNull()
    useAuthStore.getState().setAuth('jwt-xyz', user)
    expect(getToken()).toBe('jwt-xyz')
  })
})
