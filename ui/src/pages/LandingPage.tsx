import { Bot, CalendarDays, Code2, KeyRound, Rocket, UsersRound } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { FeatureCard } from '../components/FeatureCard'
import { SessionSignupCard } from '../components/SessionSignupCard'
import { StatusBadge } from '../components/StatusBadge'
import { useGetStatsQuery } from '../store/api'

export function LandingPage() {
  const { data: stats } = useGetStatsQuery()
  const { user } = useAuth()

  const startDiscord = () => {
    window.location.href = '/auth/discord/login?return_to=/waiting-list'
  }

  return (
    <main className="min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(88,101,242,0.14),_transparent_34%),radial-gradient(circle_at_top_right,_rgba(168,85,247,0.12),_transparent_30%),#fbfbff] text-slate-950">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
        <a href="/" className="flex items-center gap-2 text-lg font-bold tracking-tight">
          <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-[#5865F2] text-white">
            <Bot className="h-5 w-5" aria-hidden="true" />
          </span>
          VibeBot Sessions
        </a>
        <div className="flex items-center gap-6 text-sm font-medium text-slate-600">
          <Link className="hidden hover:text-[#5865F2] sm:inline" to="/tutorial">
            Docs
          </Link>
        </div>
      </nav>

      <section id="about" className="mx-auto grid max-w-6xl items-center gap-10 px-6 py-12 lg:grid-cols-[1.08fr_0.92fr] lg:py-20">
        <div>
          <div className="inline-flex rounded-full bg-indigo-50 px-4 py-2 text-xs font-bold tracking-[0.22em] text-[#5865F2]">
            VIBE + CODE + DISCORD
          </div>
          <h1 className="mt-6 max-w-3xl text-5xl font-black leading-[1.02] tracking-[-0.05em] text-slate-950 sm:text-6xl lg:text-7xl">
            Build a Discord Bot. <span className="text-[#5865F2]">Vibe.</span> Code. Deploy.
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-600 sm:text-xl">
            Join a live vibe coding session and build your own Discord bot. Leave with working code and your own bot credentials.
          </p>

          <div className="mt-8 flex flex-col gap-4 text-sm font-semibold text-slate-700 sm:flex-row sm:items-center">
            <div className="flex items-center gap-3">
              <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-indigo-50 text-[#5865F2]">
                <CalendarDays className="h-5 w-5" />
              </span>
              Live Sessions
            </div>
            <div className="flex items-center gap-3">
              <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-indigo-50 text-[#5865F2]">
                <UsersRound className="h-5 w-5" />
              </span>
              Limited Spots
            </div>
          </div>

          {stats && (
            <div className="mt-8 flex flex-wrap gap-3 text-xs font-semibold text-slate-500">
              <span className="rounded-full border border-slate-200 bg-white/80 px-3 py-1">{stats.total_users} users</span>
              <span className="rounded-full border border-slate-200 bg-white/80 px-3 py-1">{stats.approved_users} approved</span>
              <span className="rounded-full border border-slate-200 bg-white/80 px-3 py-1">{stats.waiting_users} waiting</span>
            </div>
          )}
        </div>

        {user ? <LoggedInStatusCard user={user} /> : <SessionSignupCard onContinueWithDiscord={startDiscord} />}
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-16 pt-4 lg:pb-24">
        <h2 className="text-center text-3xl font-black tracking-tight text-slate-950">What you get</h2>
        <div className="mt-8 grid gap-5 md:grid-cols-3">
          <FeatureCard
            icon={KeyRound}
            title="Your Bot Credentials"
            description="Get your own Discord bot token and client ID to keep."
          />
          <FeatureCard
            icon={Code2}
            title="Guided Vibe Coding"
            description="Follow along in a live session and build something awesome."
          />
          <FeatureCard
            icon={Rocket}
            title="Deploy & Use"
            description="Invite your bot, test it, and start building in your server."
          />
        </div>
      </section>

      <footer className="pb-10 text-center text-sm text-slate-400">
        Built for creators. Powered by Discord. <Link className="font-semibold text-[#5865F2] hover:text-[#4752C4]" to="/tutorial">Read the docs.</Link>
      </footer>
    </main>
  )
}

function LoggedInStatusCard({ user }: { user: NonNullable<ReturnType<typeof useAuth>['user']> }) {
  const statusCopy = {
    waiting: 'You are on the waiting list. We will review your signup and approve your bot credentials soon.',
    approved: 'You are approved. Your bot credentials and tutorial are ready.',
    rejected: 'Your signup was not approved. Contact the organizers if this looks wrong.',
    suspended: 'Your access is currently suspended. Contact the organizers for help.',
  }[user.status]

  const primaryLink = user.role === 'admin' ? '/admin' : user.status === 'approved' ? '/profile' : '/waiting-list'
  const primaryLabel = user.role === 'admin' ? 'Open admin dashboard' : user.status === 'approved' ? 'View your profile' : 'View waiting list status'

  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-2xl shadow-indigo-100/70 sm:p-8">
      <div className="mb-6 flex items-start justify-between gap-4">
        <div>
          <p className="text-sm font-bold uppercase tracking-[0.2em] text-[#5865F2]">Signed in</p>
          <h2 className="mt-2 text-2xl font-bold tracking-tight text-slate-950">Welcome, {user.display_name}</h2>
        </div>
        <StatusBadge status={user.status} />
      </div>
      <p className="text-sm leading-6 text-slate-500">{statusCopy}</p>
      <div className="mt-6 grid gap-3">
        <Link
          className="inline-flex justify-center rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white shadow-sm transition hover:bg-[#4752C4]"
          to={primaryLink}
        >
          {primaryLabel}
        </Link>
        <Link
          className="inline-flex justify-center rounded-xl border border-slate-200 px-5 py-3 text-sm font-bold text-slate-700 transition hover:border-[#5865F2] hover:text-[#5865F2]"
          to="/tutorial"
        >
          Read the docs
        </Link>
      </div>
    </section>
  )
}
