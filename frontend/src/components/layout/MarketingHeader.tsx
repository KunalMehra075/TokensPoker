import { Link } from 'react-router-dom'
import { Wordmark } from '@/components/icons/Logo'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'

// Header for public marketing and legal pages.
export function MarketingHeader() {
  const token = useAuthStore((s) => s.token)
  return (
    <header className="sticky top-0 z-30 border-b border-border bg-white/90 backdrop-blur">
      <div className="container-app flex h-16 items-center justify-between">
        <Link to="/" aria-label="FreeTokensPoker home">
          <Wordmark />
        </Link>
        <nav className="flex items-center gap-1 sm:gap-2">
          <a
            href="/#how-it-works"
            className="hidden rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:text-foreground sm:block"
          >
            How it works
          </a>
          <a
            href="/#modes"
            className="hidden rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:text-foreground sm:block"
          >
            Estimation modes
          </a>
          {token ? (
            <Link to="/app">
              <Button size="sm">Open app</Button>
            </Link>
          ) : (
            <>
              <Link to="/login">
                <Button variant="ghost" size="sm">
                  Log in
                </Button>
              </Link>
              <Link to="/login">
                <Button size="sm">Start a room</Button>
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  )
}
