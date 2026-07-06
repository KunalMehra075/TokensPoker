# Deploying the backend to Fly.io

The Go backend runs as a **single always-on Fly machine** in `iad` (US-East). It serves the REST
API and the realtime WebSocket hub. Config lives in `fly.toml`. The production image has been
verified locally (builds, boots, connects to Mongo, `/health` + `/api/modes` respond, and the
`wss://` upgrade returns `101 Switching Protocols`).

> Single instance is required: the realtime hub is in-memory with no backplane
> (`internal/realtime/hub.go`). `fly.toml` pins `min = max = 1`. Do not scale out until a Redis
> pub/sub backplane is added (see the plan's Phase 2/3).

## One-time prerequisites

1. **Install flyctl** and sign in:
   ```bash
   brew install flyctl        # or: curl -L https://fly.io/install.sh | sh
   fly auth login             # interactive; in Claude Code run: ! fly auth login
   ```
2. **MongoDB Atlas (free M0):** create a cluster, a DB user, and database `freetokenspoker`.
   Network Access → allowlist `0.0.0.0/0` (Fly egress IPs are dynamic; the DB user/password is the
   security boundary). Copy the `mongodb+srv://...` connection string.

## Deploy

Run everything from this `backend/` directory.

1. **Create the app** (uses the existing `fly.toml`, no code deploy yet):
   ```bash
   fly launch --no-deploy --copy-config --name freetokenspoker-api --region iad
   ```
   If the name is taken, pick another and update `app` in `fly.toml`.

2. **Set secrets / env** (replace the Mongo URI; keep `JWT_SECRET` strong and private):
   ```bash
   fly secrets set \
     APP_ENV=production \
     MONGODB_URI="mongodb+srv://kunalmehra:thegameison@cluster0.odybzpp.mongodb.net/jobhuntapp?retryWrites=true&w=majority" \
     MONGODB_DB=freetokenspoker \
     JWT_SECRET="TokensPokerAlphaSecret679" \
     JWT_EXPIRY_HOURS=168 \
     RATE_LIMIT_RPS=20 \
     RATE_LIMIT_BURST=40 \
     CORS_ORIGINS="https://freetokenspoker.com,https://www.freetokenspoker.com,https://freetokenspoker.pages.dev"
   ```
   `CORS_ORIGINS` gates both REST CORS and the WebSocket `CheckOrigin` (`internal/realtime/handler.go`),
   so it must list every frontend origin: the custom domain, `www`, and the `*.pages.dev` preview.

3. **Deploy:**
   ```bash
   fly deploy
   fly scale count 1   # ensure exactly one machine
   ```
   Note the URL, e.g. `https://freetokenspoker-api.fly.dev`.

## Wire up the frontend

In Cloudflare Pages → the project's environment variables:
```
VITE_API_URL = https://freetokenspoker-api.fly.dev
```
Redeploy Pages. The WebSocket URL derives automatically (`src/constants/index.ts` turns `https` →
`wss` and appends `/ws`), so no separate WS variable is needed. When the custom domain is live on
Pages, make sure it is in `CORS_ORIGINS` and re-run the `fly secrets set` above (it redeploys).

## Verify

```bash
curl https://freetokenspoker-api.fly.dev/health     # {"status":"ok","service":"freetokenspoker-api"}
fly status                                          # exactly one machine, running, in iad
fly logs                                             # Mongo connected on boot, no restarts
```
Then, in the browser: open the deployed frontend, create a room, join from a second
browser/incognito with the code, vote, and confirm the live reveal propagates (exercises
`/ws?token=`). The console should show no CORS errors and the network tab a persistent `wss://`
connection.
