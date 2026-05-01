import { ArrowLeft, LogOut, RefreshCw } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../../auth/useAuth'
import { AdminStats } from '../../components/AdminStats'
import { AdminUserTable } from '../../components/AdminUserTable'
import { useGetStatsQuery, useGetWaitlistQuery, useRejectUserMutation } from '../../store/api'

export function AdminDashboard() {
  const { logout } = useAuth()
  const { data: stats, refetch: refetchStats } = useGetStatsQuery()
  const { data: waitlist, isLoading, refetch: refetchWaitlist } = useGetWaitlistQuery()
  const [rejectUser, { isLoading: isRejecting }] = useRejectUserMutation()

  const refresh = () => {
    void refetchStats()
    void refetchWaitlist()
  }

  const reject = async (id: number) => {
    if (!window.confirm('Reject this request?')) return
    await rejectUser(id).unwrap()
  }

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(88,101,242,0.14),_transparent_35%),#fbfbff] px-6 py-8">
      <section className="mx-auto max-w-6xl">
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <Link to="/" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
              <ArrowLeft className="h-4 w-4" /> Home
            </Link>
            <h1 className="mt-4 text-4xl font-black tracking-tight text-slate-950">Admin Dashboard</h1>
            <p className="mt-2 text-slate-600">Review waiting-list requests and assign Discord bot credentials.</p>
          </div>
          <div className="flex gap-3">
            <button onClick={refresh} type="button" className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-bold text-slate-600 hover:text-[#5865F2]">
              <RefreshCw className="h-4 w-4" /> Refresh
            </button>
            <button onClick={() => void logout()} type="button" className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-bold text-slate-600 hover:text-[#5865F2]">
              <LogOut className="h-4 w-4" /> Log out
            </button>
          </div>
        </div>

        <div className="mt-8"><AdminStats stats={stats} /></div>

        <section className="mt-8 rounded-3xl border border-slate-200 bg-white/70 p-6 shadow-xl shadow-indigo-100/60 backdrop-blur">
          <div className="mb-5 flex items-center justify-between gap-4">
            <div>
              <h2 className="text-2xl font-black text-slate-950">Waiting List</h2>
              <p className="mt-1 text-sm text-slate-500">Oldest requests should usually be reviewed first.</p>
            </div>
            <span className="rounded-full bg-indigo-50 px-3 py-1 text-xs font-bold text-[#5865F2]">
              {waitlist?.total ?? 0} waiting
            </span>
          </div>
          {isLoading ? (
            <div className="rounded-2xl bg-slate-50 p-8 text-center text-slate-500">Loading waitlist…</div>
          ) : (
            <AdminUserTable users={waitlist?.users ?? []} onReject={reject} isBusy={isRejecting} />
          )}
        </section>
      </section>
    </main>
  )
}
