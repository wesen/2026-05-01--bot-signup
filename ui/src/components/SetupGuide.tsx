import { Check, Copy } from 'lucide-react'
import { useState } from 'react'
import ReactMarkdown, { type Components } from 'react-markdown'
import { PrismLight as SyntaxHighlighter } from 'react-syntax-highlighter'
import bash from 'react-syntax-highlighter/dist/esm/languages/prism/bash'
import text from 'react-syntax-highlighter/dist/esm/languages/prism/markup'
import { oneDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import remarkGfm from 'remark-gfm'
import guideMarkdown from '../content/get-your-discord-bot-running.md?raw'
import { discordBotInviteURL, type BotCredentials } from '../store/api'

SyntaxHighlighter.registerLanguage('bash', bash)
SyntaxHighlighter.registerLanguage('text', text)

interface SetupGuideProps {
  credentials?: BotCredentials
}

export function SetupGuide({ credentials }: SetupGuideProps) {
  const genericInviteURL = 'https://discord.com/oauth2/authorize?client_id=<DISCORD_APPLICATION_ID>&permissions=861140978752&integration_type=0&scope=applications.commands+bot'
  const inviteURL = credentials ? discordBotInviteURL(credentials.application_id) : genericInviteURL
  const markdown = stripFrontmatter(guideMarkdown)
    .replaceAll('<DISCORD_APPLICATION_ID>', credentials?.application_id ?? '<DISCORD_APPLICATION_ID>')
    .replaceAll(genericInviteURL, inviteURL)

  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-xl shadow-indigo-100/60">
      <div className="mb-6">
        <p className="text-sm font-black uppercase tracking-[0.2em] text-[#5865F2]">Setup guide</p>
        <h2 className="mt-2 text-2xl font-black text-slate-950">Get your Discord bot running</h2>
      </div>
      <div className="tutorial-markdown profile-guide-markdown">
        <ReactMarkdown remarkPlugins={[remarkGfm]} components={markdownComponents}>
          {markdown}
        </ReactMarkdown>
      </div>
    </section>
  )
}

const markdownComponents: Components = {
  code({ className, children, ...props }) {
    const code = String(children).replace(/\n$/, '')
    const language = /language-(\w+)/.exec(className ?? '')?.[1]

    if (!language) {
      return (
        <code className={className} {...props}>
          {children}
        </code>
      )
    }

    return <GuideCodeBlock code={code} language={language} />
  },
}

function GuideCodeBlock({ code, language }: { code: string; language: string }) {
  const [copied, setCopied] = useState(false)

  async function copyCode() {
    await navigator.clipboard.writeText(code)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1200)
  }

  return (
    <div className="tutorial-code-block">
      <div className="tutorial-code-toolbar">
        <span>{language}</span>
        <button type="button" onClick={copyCode} aria-label="Copy guide code to clipboard">
          {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
          {copied ? 'Copied' : 'Copy'}
        </button>
      </div>
      <SyntaxHighlighter language={language} style={oneDark} PreTag="div" customStyle={{ margin: 0, background: 'transparent' }}>
        {code}
      </SyntaxHighlighter>
    </div>
  )
}

function stripFrontmatter(markdown: string) {
  return markdown.replace(/^---\n[\s\S]*?\n---\n+/, '')
}
