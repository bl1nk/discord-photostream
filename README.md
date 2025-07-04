# Discord Photostream Bot

This bot deletes any message in a configured channel that does not have an
image attached.

## Configuration

Create a `.env` file in the root of the project to set necessary environment
variables:

```
DISCORD_BOT_TOKEN=YOUR_BOT_TOKEN
DISCORD_CHANNEL_ID=YOUR_CHANNEL_ID
```

If no `DISCORD_CHANNEL_ID` is set, the bot will list all guilds and channels it
has access to, along with their IDs.

