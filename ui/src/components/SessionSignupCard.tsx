import { Mail, User } from 'lucide-react'
import { DiscordOAuthButton } from './DiscordOAuthButton'

interface SessionSignupCardProps {
  onContinueWithDiscord: () => void
}

export function SessionSignupCard({ onContinueWithDiscord }: SessionSignupCardProps) {
  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-2xl shadow-indigo-100/70 sm:p-8">
      <div className="mb-6">
        <h2 className="text-2xl font-bold tracking-tight text-slate-950">Sign Up for a Session</h2>
        <p className="mt-2 text-sm leading-6 text-slate-500">
          Reserve your spot and get your bot credentials.
        </p>
      </div>

      <div className="space-y-4" aria-hidden="true">
        <label className="flex items-center gap-3 rounded-xl border border-slate-200 bg-white px-4 py-3 text-slate-400">
          <User className="h-5 w-5 text-slate-400" />
          <span>Your Name</span>
        </label>
        <label className="flex items-center gap-3 rounded-xl border border-slate-200 bg-white px-4 py-3 text-slate-400">
          <Mail className="h-5 w-5 text-slate-400" />
          <span>Email Address</span>
        </label>
      </div>

      <div className="mt-5">
        <DiscordOAuthButton onClick={onContinueWithDiscord} fullWidth />
      </div>

      <p className="mt-4 text-center text-xs leading-5 text-slate-500">
        We’ll email you session details and your bot credentials. No password needed — Discord handles sign-in.
      </p>
    </section>
  )
}
