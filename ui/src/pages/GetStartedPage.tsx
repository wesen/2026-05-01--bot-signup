import { ArrowLeft } from 'lucide-react'
import { Link } from 'react-router-dom'
import { SetupGuide } from '../components/SetupGuide'

export function GetStartedPage() {
  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(88,101,242,0.14),_transparent_35%),#fbfbff] px-6 py-10">
      <section className="mx-auto max-w-4xl">
        <Link to="/" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
          <ArrowLeft className="h-4 w-4" /> Home
        </Link>
        <div className="mt-6">
          <SetupGuide />
        </div>
      </section>
    </main>
  )
}
