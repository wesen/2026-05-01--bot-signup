import { ArrowLeft, BookOpen, Clock, LogOut } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { StatusBadge } from '../components/StatusBadge'

export function WaitingListPage() {
  const { user, logout } = useAuth()

  if (!user) return null

  const messages = {
    waiting: {
      title: 'You are on the waiting list.',
      body: "Your request is being reviewed. We'll let you know when your bot credentials are ready.",
    },
    rejected: {
      title: 'Your request was not approved.',
      body: 'If you think this is a mistake, contact the session organizer.',
    },
    suspended: {
      title: 'Your account is suspended.',
      body: 'Your bot access is paused. Contact an admin for details.',
    },
    approved: {
      title: 'You are approved.',
      body: 'Your bot credentials are ready in your profile.',
    },
  }[user.status]

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(88,101,242,0.14),_transparent_35%),#fbfbff] px-6 py-10">
      <section className="mx-auto max-w-3xl rounded-3xl border border-slate-200 bg-white p-8 shadow-2xl shadow-indigo-100/70">
        <Link to="/" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
          <ArrowLeft className="h-4 w-4" /> Home
        </Link>
        <div className="mt-8 flex items-start justify-between gap-4">
          <div>
            <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-indigo-50 text-[#5865F2]">
              <Clock className="h-7 w-7" />
            </div>
            <h1 className="mt-5 text-4xl font-black tracking-tight text-slate-950">Your signup status</h1>
          </div>
          <StatusBadge status={user.status} />
        </div>
        <div className="mt-8 rounded-2xl bg-slate-50 p-6">
          <h2 className="text-xl font-bold text-slate-950">{messages.title}</h2>
          <p className="mt-3 leading-7 text-slate-600">{messages.body}</p>
          <dl className="mt-5 grid gap-3 text-sm text-slate-600 sm:grid-cols-2">
            <div>
              <dt className="font-bold text-slate-900">Discord ID</dt>
              <dd className="font-mono">{user.discord_id}</dd>
            </div>
            <div>
              <dt className="font-bold text-slate-900">Signed in as</dt>
              <dd>{user.display_name}</dd>
            </div>
          </dl>
        </div>
        <div className="mt-8 flex flex-col gap-3 sm:flex-row">
          {user.status === 'approved' ? (
            <Link className="rounded-xl bg-[#5865F2] px-5 py-3 text-center text-sm font-bold text-white hover:bg-[#4752C4]" to="/profile">
              View profile
            </Link>
          ) : (
            <Link className="inline-flex items-center justify-center gap-2 rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white hover:bg-[#4752C4]" to="/tutorial">
              <BookOpen className="h-4 w-4" /> Read the bot tutorial
            </Link>
          )}
          <button
            type="button"
            onClick={() => void logout()}
            className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-200 px-5 py-3 text-sm font-bold text-slate-600 hover:text-[#5865F2]"
          >
            <LogOut className="h-4 w-4" /> Log out
          </button>
        </div>
      </section>
    </main>
  )
}
