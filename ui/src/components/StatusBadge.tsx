import type { UserStatus } from '../store/api'

const styles: Record<UserStatus, string> = {
  waiting: 'bg-amber-50 text-amber-700 ring-amber-200',
  approved: 'bg-emerald-50 text-emerald-700 ring-emerald-200',
  rejected: 'bg-rose-50 text-rose-700 ring-rose-200',
  suspended: 'bg-slate-100 text-slate-700 ring-slate-200',
}

const labels: Record<UserStatus, string> = {
  waiting: 'Waiting List',
  approved: 'Approved',
  rejected: 'Rejected',
  suspended: 'Suspended',
}

export function StatusBadge({ status }: { status: UserStatus }) {
  return (
    <span className={`inline-flex rounded-full px-3 py-1 text-xs font-bold ring-1 ${styles[status]}`}>
      {labels[status]}
    </span>
  )
}
