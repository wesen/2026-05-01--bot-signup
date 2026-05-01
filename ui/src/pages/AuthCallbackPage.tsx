import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'

export function AuthCallbackPage() {
  const navigate = useNavigate()
  const { user, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading) {
      navigate(user?.status === 'approved' ? '/profile' : '/waiting-list', { replace: true })
    }
  }, [isLoading, navigate, user])

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
      <section className="rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-xl">
        <h1 className="text-2xl font-black text-slate-950">Signing you in with Discord...</h1>
        <p className="mt-3 text-slate-500">This should only take a moment.</p>
      </section>
    </main>
  )
}
