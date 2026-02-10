# Trends

Visualize your Strava activity trends over time.

## Setup

```
npm install
npm --prefix client install
npm --prefix server install
cp client/.env.example client/.env
cp server/.env.example server/.env
```

Add your `STRAVA_CLIENT_ID` and `STRAVA_CLIENT_SECRET` to `server/.env`.

## Dev

```
npm run dev
```

Frontend: http://localhost:5173
Backend: http://localhost:3001

## Deploy

**Frontend** → Vercel (deploy from `client/`). Set `VITE_API_URL` to your fly.io backend URL.

**Backend** → fly.io.

```
fly launch
fly volumes create data -r sjc
fly secrets set STRAVA_CLIENT_ID=... STRAVA_CLIENT_SECRET=... FRONTEND_URL=https://your-app.vercel.app API_BASE_URL=https://your-app.fly.dev
fly deploy
```
