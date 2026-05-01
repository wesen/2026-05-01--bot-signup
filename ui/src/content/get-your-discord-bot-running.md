---
Title: "Get Your Discord Bot Running"
Slug: "get-your-discord-bot-running"
Short: "End-to-end guide for signing up, receiving credentials, inviting the bot to a server, and following the tutorial."
Topics:
- discord
- bot-signup
- oauth
- credentials
- tutorial
Commands: []
Flags: []
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This guide explains the full path from signup to a running Discord bot. First you clone or install the Discord bot runner, then the signup site handles identity through Discord OAuth, an admin approves your request, and your profile page gives you the exact environment variables and setup links you need.

## 1. Clone or install the Discord bot runner

Start by getting the Discord bot runner from GitHub:

```bash
git clone https://github.com/go-go-golems/discord-bot.git
cd discord-bot
```

If you prefer installing the CLI directly, use the install instructions from the repository:

[https://github.com/go-go-golems/discord-bot](https://github.com/go-go-golems/discord-bot)

You need this repository or installed CLI before the tutorial commands can run locally. The signup platform gives you credentials; the `discord-bot` project is what uses those credentials to run your bot.

## 2. Sign up with Discord

Open the signup site and click **Continue with Discord**. Discord asks you to authorize the signup application with the `identify` and `email` scopes, which lets the site know who you are without creating a password.

After Discord redirects you back, the site creates your account and shows your signup status. New accounts usually start in the waiting-list state.

## 3. Wait for credentials

Your account remains on the waiting list until an admin approves it. Approval means the admin has created or assigned the Discord bot credentials for you.

While you wait, keep the status page open or return to it later from the landing page. If you are signed in, the landing page shows your current status and links back to your waiting-list status page.

## 4. Open your profile and copy the environment variables

After approval, open your profile page:

[Open your bot profile](/profile)

The profile page shows a credentials table and a copyable `.envrc` block. Copy that block into your shell environment or local `.envrc` file:

```bash
export DISCORD_BOT_TOKEN='...'
export DISCORD_APPLICATION_ID='...'
export DISCORD_PUBLIC_KEY='...'
export DISCORD_GUILD_ID='...'
```

These values are secret runtime configuration. Do not paste the bot token into public chat, screenshots, issue trackers, or commits.

## 5. Request server access with the bot invite URL

Use the invite URL from your profile page to add the bot to the Discord server. The URL has this shape:

```text
https://discord.com/oauth2/authorize?client_id=<DISCORD_APPLICATION_ID>&permissions=861140978752&integration_type=0&scope=applications.commands+bot
```

The profile page fills in the `client_id` with your assigned `DISCORD_APPLICATION_ID` so you can copy the correct invite URL directly:

[Open your bot profile to copy the invite URL](/profile)

When Discord asks for a server, choose the server represented by your `DISCORD_GUILD_ID`. If you want to restrict the bot to one channel, configure Discord channel permissions after the bot joins the server: give the bot role access to the intended channel and remove or deny access elsewhere. Avoid granting `Administrator` if you want channel restrictions to work.

## 6. Read the tutorial

After the bot is invited and your environment variables are available, follow the bot tutorial:

[Read the Discord bot tutorial](/tutorial)

The tutorial explains how the JavaScript bot files work, how to run them, and how to test commands in Discord.

## 7. You are done

At this point you have:

- a signed-in signup account,
- approved Discord bot credentials,
- a local `.envrc` or shell environment with the required values,
- a bot invited to the target Discord server,
- and the tutorial for building and running the bot.

You can now start coding and testing your bot.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| The profile page says credentials are not available. | Your account is still waiting for admin approval. | Return to the waiting-list page and wait for approval. |
| Discord says the invite URL is invalid. | The `client_id` is missing or does not match your application ID. | Copy the invite URL from your profile page after approval. |
| The bot joins the server but cannot respond in a channel. | The bot role lacks channel permissions. | Allow View Channel, Send Messages, Read Message History, and Use Application Commands in the intended channel. |
| Slash commands do not appear. | Commands may not be synced yet, or the bot lacks `applications.commands` authorization. | Re-run the bot with command sync enabled and confirm the invite URL includes `scope=applications.commands+bot`. |
| The bot token was committed or shared. | The token is secret and must be rotated after exposure. | Regenerate the bot token in the Discord Developer Portal and update your `.envrc`. |

## See Also

- [Profile page](/profile)
- [Discord bot tutorial](/tutorial)
- [Discord Developer Portal](https://discord.com/developers/applications)
