import { Link } from 'react-router-dom'
import { Wordmark } from '@/components/icons/Logo'
import { BRAND } from '@/constants'

// Footer links the legal/trust pages (real routes) so they are crawlable and so
// the site meets the usual quality bar for ad networks later on.
export function Footer() {
  const year = new Date().getFullYear()
  return (
    <footer className="border-t border-border bg-muted/40">
      <div className="container-app py-12">
        <div className="flex flex-col gap-8 sm:flex-row sm:items-start sm:justify-between">
          <div className="max-w-sm">
            <Wordmark />
            <p className="mt-3 text-sm text-muted-foreground">
              Planning Poker for teams that build with AI. Estimate tokens, cost, model
              choice, and effort together, before the first prompt is written.
            </p>
          </div>
          <nav className="grid grid-cols-2 gap-x-12 gap-y-2 text-sm">
            <span className="col-span-2 font-semibold text-foreground">Company</span>
            <Link to="/about" className="text-muted-foreground hover:text-foreground">
              About
            </Link>
            <Link to="/contact" className="text-muted-foreground hover:text-foreground">
              Contact
            </Link>
            <Link to="/privacy" className="text-muted-foreground hover:text-foreground">
              Privacy
            </Link>
            <Link to="/terms" className="text-muted-foreground hover:text-foreground">
              Terms
            </Link>
          </nav>
        </div>
        <div className="mt-8 border-t border-border pt-6 text-sm text-muted-foreground">
          <p>
            {year} {BRAND.name}. Built for engineering teams entering the AI era.
          </p>
        </div>
      </div>
    </footer>
  )
}
