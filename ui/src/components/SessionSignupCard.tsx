import { ShieldCheck, UserRoundCheck } from 'lucide-react'
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

      <div className="space-y-3">
        <div className="flex items-center gap-3 rounded-xl border border-indigo-100 bg-indigo-50/70 px-4 py-3 text-sm font-semibold text-slate-700">
          <UserRoundCheck className="h-5 w-5 text-[#5865F2]" />
          <span>Your Discord profile provides your name and email.</span>
        </div>
        <div className="flex items-center gap-3 rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm font-semibold text-slate-600">
          <ShieldCheck className="h-5 w-5 text-emerald-500" />
          <span>No password or extra signup form required.</span>
        </div>
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
