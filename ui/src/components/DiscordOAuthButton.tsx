import { MessageCircle } from 'lucide-react'

interface DiscordOAuthButtonProps {
  onClick: () => void
  children?: string
  fullWidth?: boolean
}

export function DiscordOAuthButton({
  onClick,
  children = 'Continue with Discord',
  fullWidth = false,
}: DiscordOAuthButtonProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`inline-flex items-center justify-center gap-2 rounded-xl bg-[#5865F2] px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-indigo-200 transition hover:-translate-y-0.5 hover:bg-[#4752C4] focus:outline-none focus:ring-4 focus:ring-indigo-200 ${
        fullWidth ? 'w-full' : ''
      }`}
    >
      <MessageCircle className="h-5 w-5" aria-hidden="true" />
      {children}
    </button>
  )
}
