import { Ban, CheckCircle2, ExternalLink, XCircle } from 'lucide-react'
import { Link } from 'react-router-dom'
import type { User } from '../store/api'
import { StatusBadge } from './StatusBadge'

interface AdminUserTableProps {
  users: User[]
  onReject?: (id: number) => void
  onDisable?: (id: number) => void
  isBusy?: boolean
}

export function AdminUserTable({ users, onReject, onDisable, isBusy = false }: AdminUserTableProps) {
  if (users.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-8 text-center text-slate-500">
        No users match this view.
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white">
      <table className="min-w-full divide-y divide-slate-200 text-left text-sm">
        <thead className="bg-slate-50 text-xs font-bold uppercase tracking-wide text-slate-500">
          <tr>
            <th className="px-4 py-3">User</th>
            <th className="px-4 py-3">Discord ID</th>
            <th className="px-4 py-3">Status</th>
            <th className="px-4 py-3">Joined</th>
            <th className="px-4 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {users.map((user) => (
            <tr key={user.id} className="align-middle">
              <td className="px-4 py-4">
                <div className="font-bold text-slate-950">{user.display_name}</div>
                <div className="text-xs text-slate-500">{user.email || 'No email from Discord'}</div>
              </td>
              <td className="px-4 py-4 font-mono text-xs text-slate-600">{user.discord_id}</td>
              <td className="px-4 py-4"><StatusBadge status={user.status} /></td>
              <td className="px-4 py-4 text-slate-600">{new Date(user.created_at).toLocaleDateString()}</td>
              <td className="px-4 py-4">
                <div className="flex justify-end gap-2">
                  <Link
                    to={`/admin/users/${user.id}`}
                    className="inline-flex items-center gap-1 rounded-lg bg-[#5865F2] px-3 py-2 text-xs font-bold text-white hover:bg-[#4752C4]"
                  >
                    <CheckCircle2 className="h-4 w-4" /> {user.status === 'waiting' ? 'Approve' : 'View'}
                  </Link>
                  {onReject && user.status === 'waiting' && (
                    <button
                      type="button"
                      disabled={isBusy}
                      onClick={() => onReject(user.id)}
                      className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-3 py-2 text-xs font-bold text-slate-600 hover:text-rose-600 disabled:opacity-50"
                    >
                      <XCircle className="h-4 w-4" /> Reject
                    </button>
                  )}
                  {onDisable && user.status !== 'suspended' && (
                    <button
                      type="button"
                      disabled={isBusy}
                      onClick={() => onDisable(user.id)}
                      className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-3 py-2 text-xs font-bold text-slate-600 hover:text-rose-600 disabled:opacity-50"
                    >
                      <Ban className="h-4 w-4" /> Disable
                    </button>
                  )}
                  <a
                    href={`https://discord.com/users/${user.discord_id}`}
                    target="_blank"
                    rel="noreferrer"
                    className="rounded-lg border border-slate-200 p-2 text-slate-500 hover:text-[#5865F2]"
                    aria-label={`Open Discord profile for ${user.display_name}`}
                  >
                    <ExternalLink className="h-4 w-4" />
                  </a>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
