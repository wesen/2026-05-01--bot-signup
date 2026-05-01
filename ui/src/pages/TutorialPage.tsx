import { ArrowLeft } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import { Link } from 'react-router-dom'
import remarkGfm from 'remark-gfm'
import tutorialMarkdown from '../content/tutorial.md?raw'

export function TutorialPage() {
  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10">
      <article className="mx-auto max-w-4xl rounded-3xl border border-slate-200 bg-white p-6 shadow-xl sm:p-10">
        <Link to="/" className="mb-8 inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
          <ArrowLeft className="h-4 w-4" /> Home
        </Link>
        <div className="tutorial-markdown">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{tutorialMarkdown}</ReactMarkdown>
        </div>
      </article>
    </main>
  )
}
