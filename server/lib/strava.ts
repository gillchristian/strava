import { getTokens, setTokens, isTokenExpired } from './tokenStore.js';

const STRAVA_API = 'https://www.strava.com/api/v3';
const STRAVA_OAUTH = 'https://www.strava.com/oauth';

export async function exchangeCodeForTokens(code: string) {
  const body = new URLSearchParams({
    client_id: process.env.STRAVA_CLIENT_ID!,
    client_secret: process.env.STRAVA_CLIENT_SECRET!,
    code,
    grant_type: 'authorization_code',
  });

  const res = await fetch(`${STRAVA_OAUTH}/token`, {
    method: 'POST',
    body,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Token exchange failed: ${res.status} ${text}`);
  }

  const data = await res.json();
  setTokens({
    access_token: data.access_token,
    refresh_token: data.refresh_token,
    expires_at: data.expires_at,
    athlete_id: data.athlete.id,
  });

  return data;
}

export async function refreshAccessToken() {
  const tokens = getTokens();
  if (!tokens) throw new Error('No tokens to refresh');

  const body = new URLSearchParams({
    client_id: process.env.STRAVA_CLIENT_ID!,
    client_secret: process.env.STRAVA_CLIENT_SECRET!,
    grant_type: 'refresh_token',
    refresh_token: tokens.refresh_token,
  });

  const res = await fetch(`${STRAVA_OAUTH}/token`, {
    method: 'POST',
    body,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Token refresh failed: ${res.status} ${text}`);
  }

  const data = await res.json();
  setTokens({
    access_token: data.access_token,
    refresh_token: data.refresh_token,
    expires_at: data.expires_at,
    athlete_id: tokens.athlete_id,
  });

  return data;
}

export async function getValidAccessToken(): Promise<string> {
  if (isTokenExpired()) {
    await refreshAccessToken();
  }
  const tokens = getTokens();
  if (!tokens) throw new Error('No tokens available');
  return tokens.access_token;
}

export async function fetchActivities(after: number, before: number) {
  const accessToken = await getValidAccessToken();

  const params = new URLSearchParams({
    after: String(after),
    before: String(before),
    per_page: '100',
  });

  const res = await fetch(`${STRAVA_API}/athlete/activities?${params}`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Strava API error: ${res.status} ${text}`);
  }

  return res.json();
}
