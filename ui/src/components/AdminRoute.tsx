import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { ProtectedRoute } from './ProtectedRoute'

export function AdminRoute({ children }: { children: ReactNode }) {
  return (
    <ProtectedRoute>
      <AdminGate>{children}</AdminGate>
    </ProtectedRoute>
  )
}

function AdminGate({ children }: { children: ReactNode }) {
  const { user } = useAuth()

  if (user?.role !== 'admin') {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
        <section className="max-w-md rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-xl shadow-indigo-100">
          <h1 className="text-2xl font-black text-slate-950">Admin access required</h1>
          <p className="mt-3 text-slate-500">Your Discord account is signed in, but it is not an admin account.</p>
          <Link className="mt-6 inline-flex rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white hover:bg-[#4752C4]" to="/">
            Return home
          </Link>
        </section>
      </main>
    )
  }

  return <>{children}</>
}
