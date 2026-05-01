import { useMemo, useState } from 'react'
import { ArrowLeft, BookOpen, Check, Copy, LogOut, ShieldCheck } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../auth/useAuth'
import { CredentialTable } from '../components/CredentialTable'
import { EnvrcSnippet } from '../components/EnvrcSnippet'
import { SetupGuide } from '../components/SetupGuide'
import { StatusBadge } from '../components/StatusBadge'
import { discordBotInviteURL, useGetProfileQuery } from '../store/api'

export function ProfilePage() {
  const { logout } = useAuth()
  const { data, isLoading, error } = useGetProfileQuery()
  const credentials = data?.bot_credentials ?? null
  const inviteURL = useMemo(() => credentials ? discordBotInviteURL(credentials.application_id) : '', [credentials])

  if (isLoading) {
    return <ProfileShell title="Loading your bot profile..." />
  }
  if (error || !data) {
    return <ProfileShell title="Could not load profile" message="Refresh the page or sign in again." />
  }

  const { user } = data

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(88,101,242,0.14),_transparent_35%),#fbfbff] px-6 py-10">
      <section className="mx-auto max-w-5xl">
        <Link to="/" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
          <ArrowLeft className="h-4 w-4" /> Home
        </Link>
        <div className="mt-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-2xl shadow-indigo-100/70">
          <div className="flex flex-col justify-between gap-6 sm:flex-row sm:items-start">
            <div>
              <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-indigo-50 text-[#5865F2]">
                <ShieldCheck className="h-7 w-7" />
              </div>
              <h1 className="mt-5 text-4xl font-black tracking-tight text-slate-950">Your Bot Dashboard</h1>
              <p className="mt-3 text-slate-600">Signed in with Discord as {user.display_name}.</p>
            </div>
            <StatusBadge status={user.status} />
          </div>

          <div className="mt-8 grid gap-4 rounded-2xl bg-slate-50 p-6 text-sm sm:grid-cols-3">
            <div>
              <dt className="font-bold text-slate-900">Discord ID</dt>
              <dd className="mt-1 font-mono text-slate-600">{user.discord_id}</dd>
            </div>
            <div>
              <dt className="font-bold text-slate-900">Email</dt>
              <dd className="mt-1 text-slate-600">{user.email || 'Not provided by Discord'}</dd>
            </div>
            <div>
              <dt className="font-bold text-slate-900">Joined</dt>
              <dd className="mt-1 text-slate-600">{new Date(user.created_at).toLocaleDateString()}</dd>
            </div>
          </div>

          {credentials ? (
            <div className="mt-8">
              <h2 className="text-2xl font-black text-slate-950">Bot credentials</h2>
              <p className="mt-2 text-sm text-slate-500">Keep these secret. Never paste your bot token into public chat or commits.</p>
              <InviteLinkCard inviteURL={inviteURL} />
              <div className="mt-5">
                <CredentialTable credentials={credentials} />
              </div>
              <div className="mt-6">
                <EnvrcSnippet credentials={credentials} />
              </div>
              <div className="mt-8">
                <SetupGuide credentials={credentials} />
              </div>
            </div>
          ) : (
            <div className="mt-8 rounded-2xl border border-amber-200 bg-amber-50 p-5 text-amber-800">
              Your account is not approved yet, so bot credentials are not available.
            </div>
          )}

          <div className="mt-8 flex flex-col gap-3 sm:flex-row">
            <Link className="inline-flex items-center justify-center gap-2 rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white hover:bg-[#4752C4]" to="/tutorial">
              <BookOpen className="h-4 w-4" /> Read tutorial
            </Link>
            <button
              type="button"
              onClick={() => void logout()}
              className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-200 px-5 py-3 text-sm font-bold text-slate-600 hover:text-[#5865F2]"
            >
              <LogOut className="h-4 w-4" /> Log out
            </button>
          </div>
        </div>
      </section>
    </main>
  )
}

function InviteLinkCard({ inviteURL }: { inviteURL: string }) {
  const [copied, setCopied] = useState(false)

  const copyInviteURL = async () => {
    await navigator.clipboard.writeText(inviteURL)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1500)
  }

  return (
    <div className="mt-5 rounded-2xl border border-indigo-100 bg-indigo-50/60 p-4">
      <h3 className="text-sm font-black uppercase tracking-[0.18em] text-[#5865F2]">Server invite link</h3>
      <p className="mt-2 text-sm leading-6 text-slate-600">
        Send this link to a server admin so they can add your bot to the Discord server. Only a Discord server admin can complete the invite.
      </p>
      <div className="mt-3 flex flex-col gap-3 sm:flex-row">
        <a
          className="inline-flex items-center justify-center rounded-xl bg-[#5865F2] px-4 py-2 text-sm font-bold text-white shadow-sm transition hover:bg-[#4752C4]"
          href={inviteURL}
          target="_blank"
          rel="noreferrer"
        >
          Request server access / invite bot
        </a>
        <button
          type="button"
          onClick={() => void copyInviteURL()}
          className="inline-flex items-center justify-center gap-2 rounded-xl border border-indigo-200 bg-white px-4 py-2 text-sm font-bold text-[#5865F2] hover:bg-indigo-50"
        >
          {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
          {copied ? 'Copied invite link' : 'Copy invite link'}
        </button>
      </div>
      <code className="mt-3 block overflow-hidden text-ellipsis whitespace-nowrap rounded-xl bg-white px-3 py-2 font-mono text-xs text-slate-600">
        {inviteURL}
      </code>
    </div>
  )
}

function ProfileShell({ title, message }: { title: string; message?: string }) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
      <section className="rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-xl">
        <h1 className="text-2xl font-black text-slate-950">{title}</h1>
        {message && <p className="mt-3 text-slate-500">{message}</p>}
      </section>
    </main>
  )
}
