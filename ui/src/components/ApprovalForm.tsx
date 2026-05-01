import { useState } from 'react'

export interface ApprovalFormValues {
  application_id: string
  bot_token: string
  guild_id: string
  public_key: string
}

interface ApprovalFormProps {
  onSubmit: (values: ApprovalFormValues) => Promise<void> | void
  isSubmitting?: boolean
  submitLabel?: string
  submittingLabel?: string
}

const defaultGuildID = import.meta.env.VITE_DEFAULT_GUILD_ID || '822583790773862470'

const emptyValues: ApprovalFormValues = {
  application_id: '',
  bot_token: '',
  guild_id: defaultGuildID,
  public_key: '',
}

function parseCredentialPaste(text: string): Partial<ApprovalFormValues> {
  const trimmed = text.trim()
  if (!trimmed) return {}

  try {
    const parsed = JSON.parse(trimmed) as Record<string, unknown>
    const values: Partial<ApprovalFormValues> = {}
    assignIfPresent(values, 'application_id', readString(parsed, 'application_id', 'applicationId', 'client_id', 'clientId'))
    assignIfPresent(values, 'bot_token', readString(parsed, 'bot_token', 'botToken', 'token'))
    assignIfPresent(values, 'guild_id', readString(parsed, 'guild_id', 'guildId'))
    assignIfPresent(values, 'public_key', readString(parsed, 'public_key', 'publicKey'))
    return values
  } catch {
    // Fall back to parsing key/value text copied from the Developer Portal or a notes file.
  }

  const values: Partial<ApprovalFormValues> = {}
  const patterns: Array<[keyof ApprovalFormValues, RegExp]> = [
    ['application_id', /(?:application|client)[ _-]*id\s*[:=]\s*([0-9]{15,25})/i],
    ['bot_token', /(?:bot[ _-]*)?token\s*[:=]\s*([^\s,;]+)/i],
    ['guild_id', /guild[ _-]*id\s*[:=]\s*([0-9]{15,25})/i],
    ['public_key', /public[ _-]*key\s*[:=]\s*([a-f0-9]{32,})/i],
  ]

  for (const [field, pattern] of patterns) {
    const match = trimmed.match(pattern)
    if (match?.[1]) values[field] = match[1].trim().replace(/^['\"]|['\"]$/g, '')
  }

  return values
}

function readString(source: Record<string, unknown>, ...keys: string[]) {
  for (const key of keys) {
    const value = source[key]
    if (typeof value === 'string' && value.trim() !== '') return value.trim()
  }
  return undefined
}

function assignIfPresent(values: Partial<ApprovalFormValues>, field: keyof ApprovalFormValues, value: string | undefined) {
  if (value !== undefined) values[field] = value
}

export function ApprovalForm({ onSubmit, isSubmitting = false, submitLabel = 'Approve User', submittingLabel = 'Approving…' }: ApprovalFormProps) {
  const [values, setValues] = useState(emptyValues)
  const [error, setError] = useState('')

  const update = (field: keyof ApprovalFormValues, value: string) => {
    setError('')
    setValues((current) => ({ ...current, [field]: value }))
  }

  const fillFromPaste = (text: string) => {
    const parsed = parseCredentialPaste(text)
    const filledValues = Object.fromEntries(Object.entries(parsed).filter(([, value]) => value !== undefined && value.trim() !== '')) as Partial<ApprovalFormValues>
    if (Object.keys(filledValues).length === 0) return false

    setError('')
    setValues((current) => ({ ...current, ...filledValues }))
    return true
  }

  const handlePaste = (event: React.ClipboardEvent) => {
    if (fillFromPaste(event.clipboardData.getData('text'))) {
      event.preventDefault()
    }
  }

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    if (Object.values(values).some((value) => value.trim() === '')) {
      setError('All credential fields are required.')
      return
    }
    try {
      await onSubmit(values)
    } catch (err) {
      setError(errorMessage(err))
    }
  }

  return (
    <form onSubmit={submit} className="space-y-5 rounded-3xl border border-slate-200 bg-white p-6 shadow-xl shadow-indigo-100/60">
      <div>
        <h2 className="text-2xl font-black text-slate-950">Discord Bot Credentials</h2>
        <p className="mt-2 text-sm leading-6 text-slate-500">
          Paste the application values created in the Discord Developer Portal. These values unlock the user’s bot profile.
        </p>
      </div>
      {error && <div className="rounded-xl bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">{error}</div>}
      <div className="rounded-2xl border border-dashed border-indigo-200 bg-indigo-50/60 p-4">
        <label className="block">
          <span className="text-sm font-bold text-slate-800">Paste all credentials</span>
          <textarea
            rows={4}
            onPaste={handlePaste}
            placeholder={'Paste JSON or key/value text here to fill Application ID, Bot Token, Guild ID, and Public Key.'}
            className="mt-2 w-full rounded-xl border border-indigo-100 bg-white px-4 py-3 font-mono text-sm text-slate-800 outline-none transition focus:border-[#5865F2] focus:ring-4 focus:ring-indigo-100"
          />
          <span className="mt-1 block text-xs text-slate-500">Supported keys: application_id/applicationId, bot_token/botToken/token, guild_id/guildId, public_key/publicKey.</span>
        </label>
      </div>
      <CredentialInput label="Application ID" help="Developer Portal → General Information → Application ID" value={values.application_id} onChange={(v) => update('application_id', v)} onPaste={handlePaste} />
      <CredentialInput label="Bot Token" help="Developer Portal → Bot → Token" value={values.bot_token} onChange={(v) => update('bot_token', v)} onPaste={handlePaste} secret />
      <CredentialInput label="Guild ID" help="Discord server ID where the bot will run" value={values.guild_id} onChange={(v) => update('guild_id', v)} onPaste={handlePaste} />
      <CredentialInput label="Public Key" help="Developer Portal → General Information → Public Key" value={values.public_key} onChange={(v) => update('public_key', v)} onPaste={handlePaste} />
      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white shadow-lg shadow-indigo-200 transition hover:bg-[#4752C4] disabled:cursor-not-allowed disabled:opacity-60"
      >
        {isSubmitting ? submittingLabel : submitLabel}
      </button>
    </form>
  )
}

function errorMessage(err: unknown) {
  if (typeof err === 'object' && err !== null && 'data' in err) {
    const data = (err as { data?: unknown }).data
    if (typeof data === 'object' && data !== null && 'error' in data) {
      const message = (data as { error?: unknown }).error
      if (typeof message === 'string') return message
    }
    if (typeof data === 'object' && data !== null && 'message' in data) {
      const message = (data as { message?: unknown }).message
      if (typeof message === 'string') return message
    }
  }
  return 'Could not save credentials. Please refresh and try again.'
}

function CredentialInput({
  label,
  help,
  value,
  onChange,
  onPaste,
  secret = false,
}: {
  label: string
  help: string
  value: string
  onChange: (value: string) => void
  onPaste?: (event: React.ClipboardEvent) => void
  secret?: boolean
}) {
  return (
    <label className="block">
      <span className="text-sm font-bold text-slate-800">{label}</span>
      <input
        value={value}
        type={secret ? 'password' : 'text'}
        onChange={(event) => onChange(event.target.value)}
        onPaste={onPaste}
        className="mt-2 w-full rounded-xl border border-slate-200 px-4 py-3 font-mono text-sm text-slate-800 outline-none transition focus:border-[#5865F2] focus:ring-4 focus:ring-indigo-100"
      />
      <span className="mt-1 block text-xs text-slate-500">{help}</span>
    </label>
  )
}
