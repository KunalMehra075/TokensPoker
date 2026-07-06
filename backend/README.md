# FreeTokensPoker API

Go + Gin + MongoDB backend with a WebSocket realtime hub. Clean architecture:
handler -> service -> repository -> Mongo. Business logic never lives in handlers.

## Run

```bash
docker run -d --name ftp-mongo -p 27017:27017 mongo:7   # if Mongo is not already up
cp .env.example .env                                     # optional, defaults work
go run ./cmd/server
```

Or everything in containers:

```bash
docker compose up --build
```

## Layout

```
cmd/server/main.go        Entry point + graceful shutdown
internal/
  api/router.go           Routes + dependency injection
  config/                 Env-driven configuration
  middleware/             CORS, auth (JWT), logging, rate limiting, recovery
  auth/                   JWT manager (identity tokens, no OTP)
  models/                 Domain models + estimation mode catalog
  dto/                    Request/response shapes + error envelope
  repositories/           Mongo data access (indexes on email, roomCode, etc.)
  services/               Business logic, emits realtime events
  handlers/               Thin HTTP handlers
  realtime/               WebSocket hub, client pumps, event protocol
  apperr/                 Typed errors with HTTP mapping
  logger/ utils/          Structured logging + helpers
```

## Auth model

There is no OTP and no password. `POST /api/auth/login` takes `{ name, email }`, upserts a user
keyed by email, and returns a JWT used to attribute actions. Email is an identity label, not a
secret.

## Build and vet

```bash
go build ./...
go vet ./...
```

## Tests

```bash
go test ./...            # unit tests everywhere; DB-backed tests skip if no Mongo
go test ./... -race      # same, with the race detector (as CI runs it)
```

Two layers of coverage:

- **Unit tests** (always run, no dependencies): the JWT manager (`auth`), env config
  (`config`), typed errors (`apperr`), the estimation mode catalog (`models`), the realtime
  event contract (`realtime`), and the room-service helpers (`services`).
- **Integration tests** (`internal/services/flow_integration_test.go`): drive the real service
  stack against MongoDB through the full lifecycle (create room, join, create task, private
  votes, reveal, final decision, plus owner-only and voting-closed guards). They **skip** unless
  a database is reachable, so a machine without Mongo still passes. Point them at one with:

  ```bash
  TEST_MONGODB_URI=mongodb://localhost:27017 go test ./internal/services/ -v
  ```

  `TEST_MONGODB_URI` (falling back to `MONGODB_URI`) is used; each run works in a throwaway
  database that is dropped on completion. CI runs these against a `mongo:7` service container.
