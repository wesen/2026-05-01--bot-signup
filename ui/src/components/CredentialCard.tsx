import { Check, Copy, Eye, EyeOff } from 'lucide-react'
import { useState } from 'react'

interface CredentialCardProps {
  label: string
  value: string
  secret?: boolean
}

export function CredentialCard({ label, value, secret = false }: CredentialCardProps) {
  const [revealed, setRevealed] = useState(!secret)
  const [copied, setCopied] = useState(false)
  const displayValue = revealed ? value : '••••••••••••••••••••••••'

  const copy = async () => {
    await navigator.clipboard.writeText(value)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1200)
  }

  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
      <div className="mb-2 flex items-center justify-between gap-3">
        <h3 className="text-sm font-bold text-slate-700">{label}</h3>
        <div className="flex gap-2">
          {secret && (
            <button
              type="button"
              onClick={() => setRevealed((v) => !v)}
              className="rounded-lg border border-slate-200 p-2 text-slate-500 transition hover:text-[#5865F2]"
              aria-label={revealed ? `Hide ${label}` : `Show ${label}`}
            >
              {revealed ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          )}
          <button
            type="button"
            onClick={copy}
            className="rounded-lg border border-slate-200 p-2 text-slate-500 transition hover:text-[#5865F2]"
            aria-label={`Copy ${label}`}
          >
            {copied ? <Check className="h-4 w-4 text-emerald-600" /> : <Copy className="h-4 w-4" />}
          </button>
        </div>
      </div>
      <code className="block overflow-hidden text-ellipsis rounded-xl bg-slate-50 px-3 py-2 font-mono text-xs text-slate-700">
        {displayValue}
      </code>
    </div>
  )
}
