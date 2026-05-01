import type { Stats } from '../store/api'

export function AdminStats({ stats }: { stats?: Stats }) {
  const items = [
    ['Total Users', stats?.total_users ?? '—'],
    ['Waiting', stats?.waiting_users ?? '—'],
    ['Approved', stats?.approved_users ?? '—'],
    ['Bots Running', stats?.bots_running ?? '—'],
  ] as const

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {items.map(([label, value]) => (
        <div key={label} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="text-sm font-bold text-slate-500">{label}</div>
          <div className="mt-2 text-3xl font-black text-slate-950">{value}</div>
        </div>
      ))}
    </div>
  )
}
