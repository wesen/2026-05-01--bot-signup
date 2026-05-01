import { ArrowLeft } from 'lucide-react'
import { Link } from 'react-router-dom'

export function TutorialPage() {
  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10">
      <article className="mx-auto max-w-3xl rounded-3xl border border-slate-200 bg-white p-8 leading-7 text-slate-700 shadow-xl">
        <Link to="/" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-[#5865F2]">
          <ArrowLeft className="h-4 w-4" /> Home
        </Link>
        <h1 className="mt-6 text-4xl font-black tracking-tight text-slate-950">Build and Run Your Discord Bot</h1>
        <p className="mt-4">
          This placeholder tutorial page will be replaced in Phase 9 with the full markdown tutorial from the go-go-golems/discord-bot project.
        </p>
        <pre className="mt-6 overflow-x-auto rounded-2xl bg-slate-950 p-5 text-sm text-slate-100"><code>{`const { defineBot } = require("discord")

module.exports = defineBot(({ command }) => {
  command("hello", { description: "Say hello" }, async () => {
    return { content: "Hello from your bot!" }
  })
})`}</code></pre>
      </article>
    </main>
  )
}
