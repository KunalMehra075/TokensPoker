import { lazy, Suspense } from 'react'
import { Routes, Route } from 'react-router-dom'
import { Loader2 } from 'lucide-react'
import { ProtectedRoute } from '@/routes/ProtectedRoute'

// Eager: the public entry points users land on first.
import { LandingPage } from '@/pages/LandingPage'

// Lazy: everything behind a click or auth, so the landing bundle stays small.
const LoginPage = lazy(() => import('@/pages/LoginPage').then((m) => ({ default: m.LoginPage })))
const InvitePage = lazy(() => import('@/pages/InvitePage').then((m) => ({ default: m.InvitePage })))
const DashboardPage = lazy(() =>
  import('@/pages/DashboardPage').then((m) => ({ default: m.DashboardPage })),
)
const RoomPage = lazy(() => import('@/pages/RoomPage').then((m) => ({ default: m.RoomPage })))
const HistoryPage = lazy(() =>
  import('@/pages/HistoryPage').then((m) => ({ default: m.HistoryPage })),
)
const NotFoundPage = lazy(() =>
  import('@/pages/NotFoundPage').then((m) => ({ default: m.NotFoundPage })),
)
const PrivacyPage = lazy(() => import('@/pages/legal').then((m) => ({ default: m.PrivacyPage })))
const TermsPage = lazy(() => import('@/pages/legal').then((m) => ({ default: m.TermsPage })))
const AboutPage = lazy(() => import('@/pages/legal').then((m) => ({ default: m.AboutPage })))
const ContactPage = lazy(() => import('@/pages/legal').then((m) => ({ default: m.ContactPage })))

function RouteFallback() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
    </div>
  )
}

export function App() {
  return (
    <Suspense fallback={<RouteFallback />}>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/join/:code" element={<InvitePage />} />

        <Route path="/privacy" element={<PrivacyPage />} />
        <Route path="/terms" element={<TermsPage />} />
        <Route path="/about" element={<AboutPage />} />
        <Route path="/contact" element={<ContactPage />} />

        <Route element={<ProtectedRoute />}>
          <Route path="/app" element={<DashboardPage />} />
          <Route path="/room/:id" element={<RoomPage />} />
          <Route path="/history" element={<HistoryPage />} />
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Suspense>
  )
}
