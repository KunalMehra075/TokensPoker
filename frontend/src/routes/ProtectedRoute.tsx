import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

// Gate app routes behind an identity. Unauthenticated users are sent to login,
// preserving where they were headed so we can bounce them back.
export function ProtectedRoute() {
  const token = useAuthStore((s) => s.token)
  const location = useLocation()
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }
  return <Outlet />
}
