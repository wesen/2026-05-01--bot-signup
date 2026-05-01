import { Check, Copy } from 'lucide-react'
import { useState } from 'react'
import type { BotCredentials } from '../store/api'

interface CredentialTableProps {
  credentials: BotCredentials
}

const rows = [
  ['Application ID', 'application_id'],
  ['Guild ID', 'guild_id'],
  ['Public Key', 'public_key'],
  ['Bot Token', 'bot_token'],
] as const

export function CredentialTable({ credentials }: CredentialTableProps) {
  const [copiedKey, setCopiedKey] = useState<string | null>(null)

  const copy = async (key: keyof BotCredentials) => {
    await navigator.clipboard.writeText(credentials[key] ?? '')
    setCopiedKey(key)
    window.setTimeout(() => setCopiedKey(null), 1200)
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
      <table className="w-full text-left text-sm">
        <thead className="bg-slate-50 text-xs font-black uppercase tracking-[0.16em] text-slate-500">
          <tr>
            <th className="px-4 py-3">Name</th>
            <th className="px-4 py-3">Environment variable</th>
            <th className="px-4 py-3">Value</th>
            <th className="px-4 py-3 text-right">Copy</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {rows.map(([label, key]) => (
            <tr key={key}>
              <td className="px-4 py-3 font-bold text-slate-800">{label}</td>
              <td className="px-4 py-3 font-mono text-xs text-slate-500">{envNameFor(key)}</td>
              <td className="max-w-[18rem] px-4 py-3">
                <code className="block overflow-hidden text-ellipsis whitespace-nowrap rounded-lg bg-slate-50 px-2 py-1 font-mono text-xs text-slate-700">
                  {credentials[key]}
                </code>
              </td>
              <td className="px-4 py-3 text-right">
                <button
                  type="button"
                  onClick={() => copy(key)}
                  className="inline-flex rounded-lg border border-slate-200 p-2 text-slate-500 transition hover:text-[#5865F2]"
                  aria-label={`Copy ${label}`}
                >
                  {copiedKey === key ? <Check className="h-4 w-4 text-emerald-600" /> : <Copy className="h-4 w-4" />}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function envNameFor(key: (typeof rows)[number][1]) {
  switch (key) {
    case 'application_id':
      return 'DISCORD_APPLICATION_ID'
    case 'guild_id':
      return 'DISCORD_GUILD_ID'
    case 'public_key':
      return 'DISCORD_PUBLIC_KEY'
    case 'bot_token':
      return 'DISCORD_BOT_TOKEN'
  }
}
