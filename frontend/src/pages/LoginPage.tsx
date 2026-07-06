import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { ArrowLeft } from 'lucide-react'
import { Wordmark } from '@/components/icons/Logo'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { login } from '@/services/api'
import type { ApiErrorShape } from '@/services/api'
import { useAuthStore } from '@/store/auth'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'
import { PENDING_ROOM_KEY } from '@/constants'

const schema = z.object({
  name: z.string().trim().min(1, 'Tell us what to call you').max(80),
  email: z.string().trim().email('That email does not look right').max(160),
})
type FormValues = z.infer<typeof schema>

export function LoginPage() {
  useDocumentTitle('Sign in', 'Enter your name and email to start estimating with your team.')
  const setAuth = useAuthStore((s) => s.setAuth)
  const navigate = useNavigate()
  const location = useLocation()
  const from = (location.state as { from?: { pathname: string } } | null)?.from?.pathname ?? '/app'

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema) })

  const onSubmit = async (values: FormValues) => {
    try {
      const { token, user } = await login(values.name, values.email)
      setAuth(token, user)
      toast.success(`Welcome, ${user.name.split(' ')[0]}`)
      // If the user arrived from an invite link, resume the join (which then
      // asks them to confirm); otherwise go where they were headed.
      const pendingRoom = sessionStorage.getItem(PENDING_ROOM_KEY)
      navigate(pendingRoom ? `/join/${pendingRoom}` : from, { replace: true })
    } catch (err) {
      toast.error((err as ApiErrorShape).message ?? 'Could not sign you in')
    }
  }

  return (
    <div className="flex min-h-screen flex-col bg-muted/40">
      <div className="container-app py-6">
        <Link to="/" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground">
          <ArrowLeft className="h-4 w-4" />
          Back home
        </Link>
      </div>

      <div className="flex flex-1 items-center justify-center px-4 pb-16">
        <div className="w-full max-w-sm">
          <div className="mb-8 flex flex-col items-center text-center">
            <Wordmark className="mb-6" />
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Join the table
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              No password to forget. We just need a name for your cards and an email so your
              history sticks around.
            </p>
          </div>

          <form
            onSubmit={handleSubmit(onSubmit)}
            className="rounded-lg border border-border bg-white p-6 shadow-card"
          >
            <div className="mb-4">
              <Label htmlFor="name">Your name</Label>
              <Input id="name" placeholder="Ada Lovelace" autoFocus {...register('name')} />
              {errors.name && <p className="mt-1 text-xs text-danger">{errors.name.message}</p>}
            </div>
            <div className="mb-5">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="ada@yourteam.com"
                {...register('email')}
              />
              {errors.email && <p className="mt-1 text-xs text-danger">{errors.email.message}</p>}
            </div>
            <Button type="submit" className="w-full" loading={isSubmitting}>
              Continue
            </Button>
          </form>

          <p className="mt-4 text-center text-xs text-muted-foreground">
            By continuing you agree to our{' '}
            <Link to="/terms" className="underline hover:text-foreground">
              Terms
            </Link>{' '}
            and{' '}
            <Link to="/privacy" className="underline hover:text-foreground">
              Privacy Policy
            </Link>
            .
          </p>
        </div>
      </div>
    </div>
  )
}
