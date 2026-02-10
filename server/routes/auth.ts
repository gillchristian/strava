import { Router } from 'express';
import { exchangeCodeForTokens } from '../lib/strava.js';
import { getTokens, clearTokens } from '../lib/tokenStore.js';

const router = Router();

router.get('/auth/strava', (_req, res) => {
  const clientId = process.env.STRAVA_CLIENT_ID;
  const baseUrl = process.env.API_BASE_URL || `http://localhost:${process.env.PORT || 3001}`;
  const redirectUri = `${baseUrl}/auth/callback`;
  const scope = 'activity:read_all';

  const url = `https://www.strava.com/oauth/authorize?client_id=${clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${scope}&approval_prompt=auto`;

  res.redirect(url);
});

router.get('/auth/callback', async (req, res) => {
  const code = req.query.code as string;

  if (!code) {
    res.status(400).send('Missing authorization code');
    return;
  }

  try {
    await exchangeCodeForTokens(code);

    const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:5173';
    res.redirect(`${frontendUrl}/?auth=success`);
  } catch (err) {
    console.error('OAuth callback error:', err);
    res.status(500).send('Authentication failed');
  }
});

router.get('/auth/status', (_req, res) => {
  const tokens = getTokens();
  res.json({
    authenticated: !!tokens,
    athleteId: tokens?.athlete_id ?? null,
  });
});

router.post('/auth/logout', (_req, res) => {
  clearTokens();
  res.json({ ok: true });
});

export default router;
