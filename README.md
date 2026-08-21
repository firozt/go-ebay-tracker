# go-ebay-tracker

A Go application and docker container that polls the eBay Browse API for new listings matching a search, and sends alerts to Discord, ntfy, or a local sound webhook when new items appear.

## Features

- Polls eBay's Browse API (OAuth2 client credentials flow) on a configurable interval
- Auto-refreshes access tokens before they expire
- Sends alerts via:
  - **Discord** webhooks
  - **ntfy** push notifications (mobile)
  - **Local sound webhook** — plays an audio alert on a host machine, useful when running the tracker in Docker without audio hardware access
- Tracks seen listing IDs to avoid duplicate alerts (WIP)

## Requirements

- Go 1.22+
- eBay Developer account with **production** API keys (App ID / Cert ID)
  - Sandbox keys will not return real listings
  - New production keysets may be disabled pending eBay's [Marketplace Account Deletion compliance](https://developer.ebay.com/marketplace-account-deletion) — apply for an exemption if you're only reading public listing data
- (Optional) A Discord webhook URL and/or an [ntfy](https://ntfy.sh) topic for alerts

## Setup

1. Clone the repo:
   ```bash
   git clone https://github.com/firozt/go-ebay-tracker.git
   cd go-ebay-tracker
   ```

2. Create a `.env` file in the project root:
   ```env
   APP_ID=your_production_app_id
   CERT_ID=your_production_cert_id
   DEVMODE=0
   TOKEN_URL=https://api.ebay.com/identity/v1/oauth2/token
   DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
   NTFY_TOPIC=your-topic-name
   ```

3. Run it:
   ```bash
   go run ./cmd/go-ebay-tracker
   ```

## eBay API rate limits

The Browse API allows **5,000 calls/day** per application by default. At a 1-minute polling interval, this uses ~1,440 calls/day — well within the limit. Check your remaining quota via eBay's [Developer Analytics API](https://developer.ebay.com/api-docs/developer/analytics/resources/rate_limit/methods/getRateLimits) or the API Reports section of the [developer dashboard](https://developer.ebay.com/my/keys).

## Notes

- Sandbox and production environments use separate credentials and endpoints — mixing them will produce `401 invalid_client` or `Invalid access token` errors.
- If deploying via Docker without host audio passthrough, use the sound-webhook alert method to trigger playback on the host machine instead.

## License

MIT
