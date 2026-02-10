import Database from 'better-sqlite3';

interface Tokens {
  access_token: string;
  refresh_token: string;
  expires_at: number;
  athlete_id: number;
}

const db = new Database(process.env.DB_PATH || 'tokens.db');

db.exec(`
  CREATE TABLE IF NOT EXISTS tokens (
    id INTEGER PRIMARY KEY CHECK(id = 1),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    athlete_id INTEGER NOT NULL
  )
`);

export function getTokens(): Tokens | null {
  const row = db.prepare('SELECT access_token, refresh_token, expires_at, athlete_id FROM tokens WHERE id = 1').get() as Tokens | undefined;
  return row ?? null;
}

export function setTokens(tokens: Tokens): void {
  db.prepare(
    'INSERT OR REPLACE INTO tokens (id, access_token, refresh_token, expires_at, athlete_id) VALUES (1, ?, ?, ?, ?)'
  ).run(tokens.access_token, tokens.refresh_token, tokens.expires_at, tokens.athlete_id);
}

export function clearTokens(): void {
  db.prepare('DELETE FROM tokens WHERE id = 1').run();
}

export function isTokenExpired(): boolean {
  const tokens = getTokens();
  if (!tokens) return true;
  const bufferSeconds = 300; // 5 minutes
  return Date.now() / 1000 >= tokens.expires_at - bufferSeconds;
}
