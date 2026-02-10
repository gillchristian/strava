import { Router } from 'express';
import { fetchActivities } from '../lib/strava.js';
import { getTokens } from '../lib/tokenStore.js';

const router = Router();

const RUN_TYPES = new Set(['Run', 'TrailRun', 'VirtualRun']);

router.get('/api/activities', async (_req, res) => {
  const tokens = getTokens();
  if (!tokens) {
    res.status(401).json({ error: 'Not authenticated' });
    return;
  }

  try {
    const now = Math.floor(Date.now() / 1000);
    const thirtyDaysAgo = now - 30 * 24 * 60 * 60;

    const activities = await fetchActivities(thirtyDaysAgo, now);
    const runs = activities.filter(
      (a: { type?: string; sport_type?: string }) =>
        RUN_TYPES.has(a.type ?? '') || RUN_TYPES.has(a.sport_type ?? '')
    );

    res.json(runs);
  } catch (err) {
    console.error('Activities fetch error:', err);
    res.status(500).json({ error: 'Failed to fetch activities' });
  }
});

export default router;
