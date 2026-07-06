import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/types'

interface AuthState {
  token: string | null
  user: User | null
  setAuth: (token: string, user: User) => void
  setUser: (user: User) => void
  logout: () => void
  isAuthenticated: () => boolean
}

// Auth is persisted so a refresh keeps the session. The token only identifies
// a user (email-as-identity); it is not a password-grade secret.
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      setAuth: (token, user) => set({ token, user }),
      setUser: (user) => set({ user }),
      logout: () => set({ token: null, user: null }),
      isAuthenticated: () => Boolean(get().token),
    }),
    { name: 'ftp.auth' },
  ),
)

// Non-hook accessor for use inside the axios interceptor.
export function getToken(): string | null {
  return useAuthStore.getState().token
}
