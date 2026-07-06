import { Link, useNavigate } from 'react-router-dom'
import { History, LogOut } from 'lucide-react'
import { Wordmark } from '@/components/icons/Logo'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'
import { realtime } from '@/services/socket'

// Persistent header for authenticated app pages.
export function AppHeader() {
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const navigate = useNavigate()

  const onLogout = () => {
    realtime.disconnect()
    logout()
    navigate('/')
  }

  return (
    <header className="sticky top-0 z-30 border-b border-border bg-white/90 backdrop-blur">
      <div className="container-app flex h-14 items-center justify-between">
        <Link to="/app" aria-label="Go to dashboard">
          <Wordmark />
        </Link>
        <nav className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={() => navigate('/history')}>
            <History className="h-4 w-4" />
            <span className="hidden sm:inline">History</span>
          </Button>
          {user && (
            <span className="hidden items-center gap-2 sm:flex">
              <Avatar id={user.id} name={user.name} size="sm" />
              <span className="text-sm font-medium text-foreground">{user.name}</span>
            </span>
          )}
          <Button variant="ghost" size="icon" onClick={onLogout} aria-label="Log out">
            <LogOut className="h-4 w-4" />
          </Button>
        </nav>
      </div>
    </header>
  )
}
