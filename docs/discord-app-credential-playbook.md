# Discord App Credential Playbook

This playbook describes the fastest repeatable way to create a Discord application, generate a bot token, and paste the resulting credentials into the bot-signup admin approval page.

## Prerequisites

- You can log in to the Discord Developer Portal in a browser.
- You have access to MFA for the Discord account.
- The repo dependencies are installed enough to run the helper script.
- On Linux, `xclip` or `wl-copy` is available if you want automatic clipboard copy.

## One-command helper

Run from the repository root:

```bash
APP_NAME=ai-in-action--example node scripts/create_discord_application_credentials.mjs
```

Optional guild override:

```bash
DISCORD_GUILD_ID=822583790773862470 \
APP_NAME=ai-in-action--example \
node scripts/create_discord_application_credentials.mjs
```

The script uses a persistent local Playwright browser profile in `.pw-discord-profile`, so after the first run Discord login may already be remembered.

## Human checkpoints

Discord intentionally blocks full bot-token automation. Expect the script to pause for:

1. **Login** — sign in to the Developer Portal if needed.
2. **hCaptcha** — complete “I am human” after creating the app if prompted.
3. **MFA** — enter your 6-digit authentication code when resetting/generating the bot token.

After each checkpoint, the script continues polling and should resume automatically.

## Output format

At the end, the script prints and copies JSON like this:

```json
{
  "application_id": "1234567890123456789",
  "public_key": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
  "bot_token": "...",
  "guild_id": "822583790773862470"
}
```

Treat this JSON as secret because it contains the bot token.

## Approving a user in bot-signup

1. Open the admin page, for example `http://localhost:5179/admin`.
2. Open the user detail/confirmation page for the waiting user.
3. Paste the JSON bundle into **Paste all credentials**.
4. Confirm that these fields fill in:
   - Application ID
   - Bot Token
   - Guild ID
   - Public Key
5. Click **Approve User**.

The individual credential inputs also accept the full JSON paste; if the pasted text contains recognized keys, the form fills all matching fields.

## Troubleshooting

### Playwright package is missing

If Node cannot import `playwright`, install it for the repo or run through an environment where Playwright is available.

### Clipboard copy fails

Copy the JSON printed by the script manually. The admin form only needs the text.

### Discord does not show the token again

Discord bot tokens are one-time visible. If the script fails after token generation but before you copied it, reset/generate the token again.

### Admin form does not fill a field

Check the JSON key spelling. Supported keys are:

- `application_id`, `applicationId`, `client_id`, `clientId`
- `bot_token`, `botToken`, `token`
- `guild_id`, `guildId`
- `public_key`, `publicKey`

### Why not fully automate with CDP or surf?

CDP and `surf` can drive browser pages, but they do not bypass Discord’s login, hCaptcha, or MFA requirements. Playwright is currently the best fit because it can maintain state, branch on page contents, extract values, and pause cleanly for human-only checkpoints.
