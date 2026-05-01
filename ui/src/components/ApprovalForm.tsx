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
}

const defaultGuildID = import.meta.env.VITE_DEFAULT_GUILD_ID || '822583790773862470'

const emptyValues: ApprovalFormValues = {
  application_id: '',
  bot_token: '',
  guild_id: defaultGuildID,
  public_key: '',
}

export function ApprovalForm({ onSubmit, isSubmitting = false }: ApprovalFormProps) {
  const [values, setValues] = useState(emptyValues)
  const [error, setError] = useState('')

  const update = (field: keyof ApprovalFormValues, value: string) => {
    setError('')
    setValues((current) => ({ ...current, [field]: value }))
  }

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    if (Object.values(values).some((value) => value.trim() === '')) {
      setError('All credential fields are required.')
      return
    }
    await onSubmit(values)
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
      <CredentialInput label="Application ID" help="Developer Portal → General Information → Application ID" value={values.application_id} onChange={(v) => update('application_id', v)} />
      <CredentialInput label="Bot Token" help="Developer Portal → Bot → Token" value={values.bot_token} onChange={(v) => update('bot_token', v)} secret />
      <CredentialInput label="Guild ID" help="Discord server ID where the bot will run" value={values.guild_id} onChange={(v) => update('guild_id', v)} />
      <CredentialInput label="Public Key" help="Developer Portal → General Information → Public Key" value={values.public_key} onChange={(v) => update('public_key', v)} />
      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-bold text-white shadow-lg shadow-indigo-200 transition hover:bg-[#4752C4] disabled:cursor-not-allowed disabled:opacity-60"
      >
        {isSubmitting ? 'Approving…' : 'Approve User'}
      </button>
    </form>
  )
}

function CredentialInput({
  label,
  help,
  value,
  onChange,
  secret = false,
}: {
  label: string
  help: string
  value: string
  onChange: (value: string) => void
  secret?: boolean
}) {
  return (
    <label className="block">
      <span className="text-sm font-bold text-slate-800">{label}</span>
      <input
        value={value}
        type={secret ? 'password' : 'text'}
        onChange={(event) => onChange(event.target.value)}
        className="mt-2 w-full rounded-xl border border-slate-200 px-4 py-3 font-mono text-sm text-slate-800 outline-none transition focus:border-[#5865F2] focus:ring-4 focus:ring-indigo-100"
      />
      <span className="mt-1 block text-xs text-slate-500">{help}</span>
    </label>
  )
}
