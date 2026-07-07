# Deploying the backend to AWS EC2

The Go backend runs as a **single container** on one EC2 instance, behind **Caddy** which
terminates TLS (auto Let's Encrypt) and transparently proxies WebSockets. Config lives in
`backend/deploy/` (`docker-compose.yml`, `Caddyfile`, `.env.example`). The production image is
already verified locally (builds, boots, connects to Mongo, `/health` + `/api/modes` respond, and
the `wss://` upgrade returns `101 Switching Protocols`).

> Single instance is required: the realtime hub is in-memory with no backplane
> (`internal/realtime/hub.go`). Run exactly one `api` container. Scaling out needs a Redis
> backplane first (see the hosting plan's Phase 3).

Split of work: **you** do the EC2 console + DNS steps (web UI); **I** do everything on the box over
SSH (Docker, deploy, verify).

---

## Part A — You: launch the EC2 instance (AWS console)

EC2 → Launch instance:
- **Name:** `freetokenspoker-api`
- **AMI:** **Amazon Linux 2023** (labeled "Free tier eligible" in the wizard), **64-bit (Arm)**
  architecture. AL2023 is the current generation (AL2 hits end-of-life June 30 2026 — do not use it);
  no newer Amazon Linux ships in 2025/2026. SSH user is **`ec2-user`**.
  - To always pull the newest AL2023 image, its public SSM alias is
    `/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-arm64` (`...-x86_64` for x86).
- **Instance type:** **`t4g.small`** (Arm Graviton, 2 GB — free-trial-eligible through Dec 31 2026 and
  enough RAM to build the image). `t3.small` with the x86_64 AMI is the x86 equivalent. Avoid the
  1 GB micros (too tight for the build).
- **Key pair:** create one and **download the `.pem`**.
- **Network / Security group:** create one with these inbound rules:
  - SSH **22** — source *My IP* (preferred) or `0.0.0.0/0`
  - HTTP **80** — source `0.0.0.0/0`  (ACME challenge + redirect)
  - HTTPS **443** — source `0.0.0.0/0`  (HTTPS + WSS)
- **Storage:** 20 GB gp3 (the 8 GB default is tight once Docker images land).
- Launch. (AL2023 has no restrictive OS firewall, so the security group is the only gate.)

Then allocate an **Elastic IP** (EC2 → Elastic IPs → Allocate → Associate with the instance) so the
public IP is stable across stop/start. Note this IP.

## Part B — You: point DNS at the instance

Add a DNS record for the API host (assumed **`api.freetokenspoker.com`** — change in
`backend/deploy/Caddyfile` if you prefer another). On Cloudflare:
- **A record**: `api` → `<Elastic IP>`, **Proxy status = DNS only (grey cloud)** so Caddy can obtain
  a Let's Encrypt certificate. WebSockets still work fine DNS-only.

(A real domain is required for valid TLS/`wss://`; Let's Encrypt won't issue for the
`*.compute.amazonaws.com` hostname.)

## Part C — You: provision + deploy (run these on the VM)

Ensure `dig +short api.freetokenspoker.com` already returns your Elastic IP (grey-cloud/DNS-only)
before you start, so Caddy can obtain the certificate.

**Step 0 — SSH in** (default user for Amazon Linux is `ec2-user`):
```bash
ssh -i /path/to/your-key.pem ec2-user@<ELASTIC_IP>
```

**Step 1 — install Docker + Compose + git, then reconnect** (re-login applies the docker group):
```bash
sudo dnf update -y
sudo dnf install -y docker docker-compose-plugin git
sudo systemctl enable --now docker
sudo usermod -aG docker ec2-user
exit
# reconnect, then sanity-check:
# ssh -i /path/to/your-key.pem ec2-user@<ELASTIC_IP>
docker version && docker compose version
```

**Step 2 — get the code:**
```bash
git clone https://github.com/KunalMehra075/TokensPoker.git ~/TokensPoker
cd ~/TokensPoker/backend/deploy
```
If the repo is private, git prompts for credentials: use a GitHub Personal Access Token as the
password.

**Step 3 — create `.env`** (secrets live only on the VM; this file is gitignored):
```bash
cp .env.example .env
openssl rand -base64 48        # copy this output -> JWT_SECRET
nano .env                     # fill the values below, then Ctrl+O, Enter, Ctrl+X
```
Set in `.env`:
- `MONGODB_URI=` your **rotated** Atlas connection string
- `MONGODB_DB=freetokenspoker`
- `JWT_SECRET=` the generated value
- `CORS_ORIGINS=https://freetokenspoker.com,https://www.freetokenspoker.com,https://<project>.pages.dev`
- keep `APP_ENV=production` and `PORT=8080`

**Step 4 — launch** (builds the image on the VM; starts `api` + Caddy; Caddy fetches TLS automatically):
```bash
docker compose up -d --build
```

## Part E — Wire up the frontend

Cloudflare Pages → environment variables:
```
VITE_API_URL = https://api.freetokenspoker.com
```
Redeploy Pages. `WS_URL` derives automatically (`src/constants/index.ts` turns `https` → `wss` and
appends `/ws`). Ensure the Pages domain(s) are in `CORS_ORIGINS` in the VM `.env`.

## Verify

```bash
docker compose ps                                # api = healthy, caddy = up
docker compose logs caddy | grep -i certificate  # Caddy obtained the Let's Encrypt cert
docker compose logs api | tail                   # Mongo connected on boot, no restarts
curl -sS https://api.freetokenspoker.com/health  # {"status":"ok","service":"freetokenspoker-api"}
```
Then in the browser: open the deployed frontend, create a room, join from a second
browser/incognito with the code, vote, and confirm the live reveal propagates (exercises
`/ws?token=`). Console shows no CORS errors and a persistent `wss://` connection.

### Troubleshooting
- **Caddy can't get a cert** → the `api` A record must be **DNS-only (grey cloud)** on Cloudflare and
  ports **80 + 443** open in the security group. Watch `docker compose logs -f caddy`.
- **`api` keeps restarting / Mongo error** → Atlas → Network Access must allowlist **`0.0.0.0/0`**, and
  re-check the rotated password in `MONGODB_URI`.
- **`docker: permission denied`** → you skipped the reconnect after step 1; re-SSH, or prefix `sudo`.
- **Redeploy after a code change:** `cd ~/TokensPoker/backend/deploy && git pull && docker compose up -d --build`.

---

## Cost note (using your AWS credits)

`t3.small` on-demand is ~$15/mo, `t4g.small` ~$12/mo, plus ~$3.65/mo for the Elastic IP and a little
EBS — all covered by your credits for now. When credits run low, downsize to `t4g.small` and/or
consider a Compute Savings Plan, or revisit Fly (~$2/mo) per the hosting plan.
