import { Check, Copy } from 'lucide-react'
import { useMemo, useState } from 'react'
import { discordBotInviteURL, type BotCredentials } from '../store/api'

interface EnvrcSnippetProps {
  credentials: BotCredentials
}

export function EnvrcSnippet({ credentials }: EnvrcSnippetProps) {
  const [copied, setCopied] = useState(false)
  const envrc = useMemo(() => buildEnvrc(credentials), [credentials])
  const inviteURL = useMemo(() => discordBotInviteURL(credentials.application_id), [credentials.application_id])

  const copy = async () => {
    await navigator.clipboard.writeText(envrc)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1200)
  }

  return (
    <section className="rounded-2xl border border-slate-200 bg-slate-950 p-4 text-white shadow-sm">
      <div className="mb-3 flex flex-col justify-between gap-3 sm:flex-row sm:items-center">
        <div>
          <h3 className="text-sm font-black uppercase tracking-[0.2em] text-indigo-200">Environment setup</h3>
          <p className="mt-1 text-sm text-slate-300">Put this in your shell environment or local <code>.envrc</code>.</p>
        </div>
        <button
          type="button"
          onClick={copy}
          className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-700 bg-slate-900 px-3 py-2 text-sm font-bold text-slate-100 transition hover:border-indigo-300 hover:text-white"
          aria-label="Copy Discord environment variables"
        >
          {copied ? <Check className="h-4 w-4 text-emerald-400" /> : <Copy className="h-4 w-4" />}
          {copied ? 'Copied' : 'Copy envrc'}
        </button>
      </div>
      <pre className="overflow-x-auto rounded-xl bg-slate-900 p-4 text-xs leading-6 text-slate-100">
        <code>{envrc}</code>
      </pre>
      <div className="mt-4 rounded-xl border border-slate-800 bg-slate-900 p-4">
        <p className="text-sm font-bold text-slate-100">Invite this bot to your Discord server</p>
        <a className="mt-2 block break-all font-mono text-xs leading-6 text-indigo-200 hover:text-indigo-100" href={inviteURL} target="_blank" rel="noreferrer">
          {inviteURL}
        </a>
      </div>
    </section>
  )
}

function buildEnvrc(credentials: BotCredentials) {
  return [
    `export DISCORD_BOT_TOKEN=${shellQuote(credentials.bot_token)}`,
    `export DISCORD_APPLICATION_ID=${shellQuote(credentials.application_id)}`,
    `export DISCORD_PUBLIC_KEY=${shellQuote(credentials.public_key)}`,
    `export DISCORD_GUILD_ID=${shellQuote(credentials.guild_id)}`,
  ].join('\n')
}

function shellQuote(value: string) {
  return `'${value.replaceAll("'", `'\\''`)}'`
}
