# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
npm run dev          # Run frontend (localhost:5173) + backend (localhost:3001) concurrently
npm run build        # TypeScript check + Vite build (client only)
npm run start        # Production server (NODE_ENV=production tsx index.ts in server/)
npm run lint         # ESLint across the client
```

## Architecture

Cadence — a monthly snapshot of your running metrics. Full-stack TypeScript app that visualizes Strava running activity data. Split deployment: React frontend on Vercel, Express backend on Fly.io. Monorepo with separate `client/` and `server/` packages.

**Frontend (`client/`)** — React 19 + Vite + Tailwind CSS v4 + Recharts. Custom hooks (`useAuth`, `useActivities`, `useChartData`) encapsulate auth flow, data fetching with localStorage caching, and chart data normalization (0-1 range). API client in `client/src/lib/api.ts`.

**Backend (`server/`)** — Express 5 with two route groups:
- `routes/auth.ts` — Strava OAuth 2.0 flow (redirect, callback, token exchange)
- `routes/activities.ts` — Fetches running activities from Strava API (filters to Run/TrailRun/VirtualRun, 30-day window)
- `lib/tokenStore.ts` — SQLite (better-sqlite3) for persistent token storage with auto-refresh (5-min buffer)
- `lib/strava.ts` — Strava API client

**Data flow:** User authenticates via Strava OAuth → tokens stored in SQLite → backend fetches activities from Strava API → frontend caches and visualizes with Recharts line charts.

## Environment

Each package has its own `.env.example`. Copy and fill in:
- `server/.env` — `STRAVA_CLIENT_ID`, `STRAVA_CLIENT_SECRET`, and optionally `FRONTEND_URL`, `API_BASE_URL`, `DB_PATH`, `PORT`
- `client/.env` — `VITE_API_URL` (backend URL)

## Deployment

Frontend deploys to Vercel from `client/` (set `VITE_API_URL` to backend URL). Backend deploys to Fly.io via Docker (GitHub Actions on push to main). Fly.io uses a persistent volume at `/data` for the SQLite database.
