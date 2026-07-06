# FreeTokensPoker

Planning Poker for the AI era. Engineering teams already debate which model to use and how
many tokens a feature will burn. FreeTokensPoker turns that into a structured estimation ritual:
everyone votes in private, the room reveals together, and the team discusses the differences.

This repo contains the full V1 MVP: a Go backend and a React frontend, both implemented and
verified end to end.

## What it does

- Frictionless identity: enter a name and email, no password and no OTP. Email is an identity
  label so your rooms, votes, and history persist.
- Rooms with a short shareable code and a shareable invite link (`/join/:code`). An invitee
  sees the room name, signs in (the code is held through login), then confirms the join.
- The owner runs the session, members vote.
- Four estimation modes, served from the backend so cards are never hardcoded:
  - AI Tokens (500K to 50M)
  - AI Cost ($1 to $250)
  - Engineering Days (1 to 21)
  - Best AI Model (GPT, Claude, Gemini, ...)
- Private vote, then a simultaneous reveal. Cards stay hidden until the owner reveals.
- Final decision per task, then the round is archived to history.
- Real-time everything: presence, votes coming in, reveal, and final decision are pushed live
  over WebSockets.
- History across all your rooms.

## Tech stack

| Layer     | Choice |
|-----------|--------|
| Frontend  | React 19, Vite, TypeScript, Tailwind CSS, TanStack Query, Zustand, React Router, React Hook Form, Zod, Axios |
| Backend   | Go 1.24+, Gin, MongoDB, JWT, gorilla/websocket |
| Realtime  | WebSocket hub with a typed JSON event protocol (see note below) |
| Database  | MongoDB |

### Note on realtime transport

The architecture doc named Socket.IO. The Go Socket.IO server ports lag the Socket.IO v4 wire
protocol and are not well maintained, which makes them a reliability risk. This build instead
uses a small, fully controlled `gorilla/websocket` hub with a typed envelope protocol
(`{event, roomId, payload}`) behind a `Broadcaster` interface. It delivers the same event-driven
collaboration the doc asked for, and it is verified working across two live browser sessions.
Swapping the transport later only touches `internal/realtime` on the backend and
`src/services/socket.ts` on the frontend.

## Repository layout

```
backend/    Go API + realtime hub (clean architecture: handler -> service -> repository -> Mongo)
frontend/   React SPA (pages, components/ui, services, store, hooks)
docs/        PRD, architecture, strategy
```

## Running locally

Prerequisites: Go 1.24+, Node 20+, Docker (for MongoDB).

### 1. Start MongoDB

```bash
docker run -d --name ftp-mongo -p 27017:27017 -v ftp_mongo_data:/data/db mongo:7
```

### 2. Start the backend

```bash
cd backend
cp .env.example .env        # optional; defaults work out of the box
go run ./cmd/server
# API on http://localhost:8080
```

### 3. Start the frontend

```bash
cd frontend
npm install
npm run dev
# App on http://localhost:5173
```

Open `http://localhost:5173`, create a room, and share the code. To see the realtime reveal,
open a second browser (or an incognito window), join with the code, and vote.

### Run the whole backend with Docker Compose

```bash
cd backend
docker compose up --build   # starts Mongo + API together
```

## Configuration

### Backend (`backend/.env`)

| Variable           | Default                          | Purpose |
|--------------------|----------------------------------|---------|
| `APP_ENV`          | `development`                    | `production` switches Gin to release mode |
| `PORT`             | `8080`                           | HTTP port |
| `MONGODB_URI`      | `mongodb://localhost:27017`      | Mongo connection string |
| `MONGODB_DB`       | `freetokenspoker`                | Database name |
| `JWT_SECRET`       | `dev-insecure-secret-change-me`  | Signs the identity JWT (change in prod) |
| `JWT_EXPIRY_HOURS` | `168`                            | Token lifetime |
| `CORS_ORIGINS`     | `http://localhost:5173,http://localhost:4173` | Browser origin allowlist |
| `RATE_LIMIT_RPS`   | `20`                             | Per-IP request rate |

### Frontend (`frontend/.env`)

| Variable        | Default                  | Purpose |
|-----------------|--------------------------|---------|
| `VITE_API_URL`  | `http://localhost:8080`  | Backend base URL (WS URL is derived from it) |
| `VITE_WS_URL`   | derived                  | Override the WebSocket URL if needed |

## API surface

```
GET    /health
GET    /api/modes                 Estimation mode catalog
GET    /api/invite/:code          Public room preview for an invite link
POST   /api/auth/login            { name, email } -> { token, user }
GET    /api/me

POST   /api/rooms                 Create a room (owner)
GET    /api/rooms/:id
POST   /api/rooms/join            { roomCode }
DELETE /api/rooms/:id             Owner only

POST   /api/tasks                 Create estimation (owner)
GET    /api/tasks/:id
GET    /api/rooms/:id/tasks
PATCH  /api/tasks/:id/reveal      Owner only
PATCH  /api/tasks/:id/final       { finalValue } owner only

POST   /api/votes                 { taskId, selectedCard }
PATCH  /api/votes                 Change a vote

GET    /api/history

GET    /ws?token=<jwt>            WebSocket upgrade
```

Realtime events (server to client): `presence`, `member_joined`, `task_created`,
`vote_received`, `votes_revealed`, `final_decision`, `task_closed`.

## Verification

- Backend: full REST flow plus realtime broadcast verified end to end (auth, room, join, task,
  hidden votes pre-reveal, reveal, final decision, history).
- Frontend: a two-browser end-to-end run (owner + member) covering login, room creation, join,
  voting, live reveal, final decision, archival, and history. Production build and strict
  TypeScript both pass.

## Deployment notes

- Frontend is a static SPA. Build with `npm run build`; deploy `dist/` to Cloudflare Pages or
  Vercel. `public/_redirects` and `public/_headers` are included for Cloudflare Pages
  (SPA fallback, preview noindex, security headers).
- Backend is a single Go binary (see `backend/Dockerfile`). Deploy to Cloud Run, Fly.io,
  Railway, or any container host. Provision MongoDB Atlas and set `MONGODB_URI`, `JWT_SECRET`,
  and `CORS_ORIGINS` to the live frontend origin.
- SEO: `index.html` ships real meta tags, Open Graph, and JSON-LD (WebApplication + FAQPage).
  `robots.txt` and `sitemap.xml` are in `frontend/public/`. Public routes (`/`, `/about`,
  `/contact`, `/privacy`, `/terms`) are real, footer-linked pages.

## Out of scope for V1

Organizations, analytics, billing, AI API integrations, Slack/Jira/GitHub, SSO, notifications.
These are deliberately excluded to keep the MVP small and the workflow sharp.
