import { ArrowLeft } from 'lucide-react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { ApprovalForm, type ApprovalFormValues } from '../../components/ApprovalForm'
import { StatusBadge } from '../../components/StatusBadge'
import { useApproveUserMutation, useGetAdminUsersQuery } from '../../store/api'

export function AdminUserDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const numericID = Number(id)
  const { data, isLoading } = useGetAdminUsersQuery({ page: 1, per_page: 200 })
  const [approveUser, { isLoading: isApproving }] = useApproveUserMutation()
  const user = data?.users.find((candidate) => candidate.id === numericID)

  const approve = async (values: ApprovalFormValues) => {
    await approveUser({ id: numericID, ...values }).unwrap()
    navigate('/admin')
  }

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(88,101,242,0.14),_transparent_35%),#fbfbff] px-6 py-8">
      <section className="mx-auto max-w-5xl">
        <Link to="/admin" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
          <ArrowLeft className="h-4 w-4" /> Back to admin
        </Link>
        <div className="mt-6 grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
          <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-xl shadow-indigo-100/60">
            {isLoading ? (
              <p className="text-slate-500">Loading user…</p>
            ) : user ? (
              <>
                <StatusBadge status={user.status} />
                <h1 className="mt-5 text-3xl font-black tracking-tight text-slate-950">{user.display_name}</h1>
                <dl className="mt-6 space-y-4 text-sm">
                  <div>
                    <dt className="font-bold text-slate-900">Discord ID</dt>
                    <dd className="mt-1 font-mono text-slate-600">{user.discord_id}</dd>
                  </div>
                  <div>
                    <dt className="font-bold text-slate-900">Email</dt>
                    <dd className="mt-1 text-slate-600">{user.email || 'No email from Discord'}</dd>
                  </div>
                  <div>
                    <dt className="font-bold text-slate-900">Joined</dt>
                    <dd className="mt-1 text-slate-600">{new Date(user.created_at).toLocaleString()}</dd>
                  </div>
                </dl>
              </>
            ) : (
              <p className="text-slate-500">User not found.</p>
            )}
          </section>
          {user && <ApprovalForm onSubmit={approve} isSubmitting={isApproving} />}
        </div>
      </section>
    </main>
  )
}
