import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Wordmark } from '@/components/icons/Logo'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'

export function NotFoundPage() {
  useDocumentTitle('Page not found')
  return (
    <div className="flex min-h-screen flex-col items-center justify-center px-4 text-center">
      <Wordmark className="mb-8" />
      <p className="text-5xl font-bold text-foreground">404</p>
      <h1 className="mt-3 text-xl font-semibold text-foreground">This page folded early</h1>
      <p className="mt-2 max-w-sm text-muted-foreground">
        The page you are after is not here. It either moved, never existed, or got estimated at
        zero story points and quietly dropped from the sprint.
      </p>
      <Link to="/" className="mt-8">
        <Button>Back to home</Button>
      </Link>
    </div>
  )
}
