import { useState, useEffect } from 'react';
import { apiFetch, API_URL } from '../lib/api';

interface AuthStatus {
  authenticated: boolean;
  athleteId: number | null;
}

export function useAuth() {
  const [authenticated, setAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiFetch<AuthStatus>('/auth/status')
      .then((data) => setAuthenticated(data.authenticated))
      .catch(() => setAuthenticated(false))
      .finally(() => setLoading(false));
  }, []);

  function login() {
    window.location.href = API_URL + '/auth/strava';
  }

  async function logout() {
    await fetch(API_URL + '/auth/logout', { method: 'POST' });
    setAuthenticated(false);
  }

  return { authenticated, loading, login, logout };
}
