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

      <DiscordOAuthButton onClick={onContinueWithDiscord} fullWidth />

      <p className="mt-4 text-center text-xs leading-5 text-slate-500">
        We’ll email you session details and your bot credentials. No password needed — Discord handles sign-in.
      </p>
    </section>
  )
}
