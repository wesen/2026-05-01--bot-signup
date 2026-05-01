#!/usr/bin/env node
/**
 * Semi-automates creating Discord Developer Portal app credentials.
 *
 * It intentionally pauses for human-only steps such as login, hCaptcha, and MFA.
 * Usage:
 *   APP_NAME=ai-in-action--example node scripts/create_discord_application_credentials.mjs
 *
 * Optional:
 *   DISCORD_GUILD_ID=822583790773862470 APP_NAME=... node scripts/...
 *   HEADLESS=1 APP_NAME=... node scripts/...   # not recommended because Discord often needs human checks
 */
import { chromium } from 'playwright'
import { spawnSync } from 'node:child_process'

const appName = process.env.APP_NAME || process.argv[2]
const guildId = process.env.DISCORD_GUILD_ID || '822583790773862470'
const headless = process.env.HEADLESS === '1'

if (!appName) {
  console.error('Missing APP_NAME. Example: APP_NAME=ai-in-action--slono node scripts/create_discord_application_credentials.mjs')
  process.exit(2)
}

const userDataDir = process.env.PLAYWRIGHT_USER_DATA_DIR || '.pw-discord-profile'
const browser = await chromium.launchPersistentContext(userDataDir, {
  headless,
  viewport: { width: 1440, height: 1000 },
})
const page = browser.pages()[0] || await browser.newPage()

async function waitForHuman(message, predicate, timeoutMs = 10 * 60 * 1000) {
  console.log(`\nACTION REQUIRED: ${message}`)
  const started = Date.now()
  while (Date.now() - started < timeoutMs) {
    if (await predicate().catch(() => false)) return
    await page.waitForTimeout(1000)
  }
  throw new Error(`Timed out waiting for: ${message}`)
}

function copyToClipboard(text) {
  const copy = spawnSync('xclip', ['-selection', 'clipboard'], { input: text })
  if (copy.status === 0) return true
  const wl = spawnSync('wl-copy', [], { input: text })
  return wl.status === 0
}

await page.goto('https://discord.com/developers/applications', { waitUntil: 'domcontentloaded' })
await waitForHuman('Log in to Discord Developer Portal if needed.', async () => {
  return await page.getByRole('button', { name: /New Application/i }).isVisible({ timeout: 1000 })
})

await page.getByRole('button', { name: /New Application/i }).click()
await page.getByRole('textbox', { name: /Name/i }).fill(appName)
await page.getByRole('checkbox').check({ force: true })
await page.getByRole('button', { name: /^Create$/i }).click()

await waitForHuman('Complete hCaptcha if Discord asks.', async () => {
  return /\/developers\/applications\/\d+\/information/.test(page.url())
})

const applicationId = page.url().match(/applications\/(\d+)\//)?.[1]
if (!applicationId) throw new Error(`Could not determine application ID from URL: ${page.url()}`)

const infoText = await page.locator('main').innerText()
const publicKey = infoText.match(/Public Key\s+([a-f0-9]{64,})/i)?.[1]
if (!publicKey) throw new Error('Could not parse Public Key from General Information page')

await page.goto(`https://discord.com/developers/applications/${applicationId}/bot`, { waitUntil: 'domcontentloaded' })
await page.getByRole('button', { name: /Reset Token/i }).click()
await page.getByRole('button', { name: /Yes, do it!/i }).click()

await waitForHuman('Complete MFA if Discord asks.', async () => {
  const text = await page.locator('main').innerText()
  return /A new token was generated/i.test(text) && /Token\s+For security purposes/i.test(text)
})

const botText = await page.locator('main').innerText()
const token = botText.match(/Token\s+For security purposes[^\n]*\n([^\s]+)\s+Copy/i)?.[1]
  || botText.match(/(M[\w.-]{40,})/)?.[1]
if (!token) throw new Error('Could not parse one-time bot token from Bot page')

const credentials = {
  application_id: applicationId,
  public_key: publicKey,
  bot_token: token,
  guild_id: guildId,
}
const json = JSON.stringify(credentials, null, 2)
console.log('\nDiscord credentials bundle:\n')
console.log(json)
console.log(copyToClipboard(json) ? '\nCopied to clipboard.' : '\nCould not copy to clipboard; copy JSON above manually.')

await browser.close()
