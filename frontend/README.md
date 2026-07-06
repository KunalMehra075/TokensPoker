# FreeTokensPoker Web

React 19 + Vite + TypeScript SPA. Tailwind for styling, TanStack Query for server state,
Zustand for auth/client state, and a small reconnecting WebSocket client for realtime.

## Run

```bash
npm install
npm run dev        # http://localhost:5173
```

Point it at the backend with `VITE_API_URL` in `.env` (defaults to `http://localhost:8080`).

## Scripts

```bash
npm run dev        # dev server
npm run build      # type-check + production build to dist/
npm run preview    # serve the production build
npm run test       # run the Vitest suite once
npm run test:watch # watch mode
```

## Tests

Vitest + Testing Library in a jsdom environment (`vitest.config.ts`, setup in `src/test/`).
Covered: the voting math (`lib/voting`), class + avatar helpers (`lib/utils`), the persisted
auth store (`store/auth`), brand/socket constants and invite links (`constants`), and a `Button`
component smoke test. Test files (`*.test.ts[x]`) live next to the code they cover and are
excluded from the production `tsc` build.

## Layout

```
src/
  pages/            Landing, Login, Dashboard, Room, History, NotFound, legal/
  components/
    ui/             Hand-built shadcn-style primitives (button, input, card, dialog, ...)
    layout/         Headers + footer
    room/ task/ voting/   Feature components
    icons/          Brand SVG (no emoji anywhere, per house rules)
  services/         api.ts (axios), socket.ts (WebSocket client)
  store/            Zustand auth store (persisted)
  hooks/            useRoomSocket, useModes, useDocumentTitle
  lib/              utils (cn), voting math, queryClient
  types/ constants/ Shared types + brand/env config
```

## Design system

Manrope typeface, white background, gray-200 borders, blue-600 accent. No gradients, no
glassmorphism, minimal motion. Tokens live as CSS variables in `src/index.css` and map to
Tailwind theme colors in `tailwind.config.js`. The aesthetic target is GitHub / Linear / Stripe.

## Realtime

`src/services/socket.ts` is a reconnecting WebSocket wrapper that speaks the backend's typed
envelope protocol. `useRoomSocket` joins a room and routes events into the TanStack Query cache.
