import type { ReactNode } from 'react'
import { useAuth } from '../auth/useAuth'

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { user, isLoading, loginWithDiscord } = useAuth()

  if (isLoading) {
    return <FullPageMessage title="Checking your Discord session..." />
  }

  if (!user) {
    return (
      <FullPageMessage
        title="Sign in with Discord"
        message="Use Discord OAuth to view your session status and bot credentials."
        actionLabel="Continue with Discord"
        onAction={() => loginWithDiscord(window.location.pathname)}
      />
    )
  }

  return <>{children}</>
}

function FullPageMessage({
  title,
  message,
  actionLabel,
  onAction,
}: {
  title: string
  message?: string
  actionLabel?: string
  onAction?: () => void
}) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
      <section className="max-w-md rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-xl shadow-indigo-100">
        <h1 className="text-2xl font-black text-slate-950">{title}</h1>
        {message && <p className="mt-3 text-slate-500">{message}</p>}
        {actionLabel && onAction && (
          <button
            type="button"
            onClick={onAction}
            className="mt-6 rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white hover:bg-[#4752C4]"
          >
            {actionLabel}
          </button>
        )}
      </section>
    </main>
  )
}
