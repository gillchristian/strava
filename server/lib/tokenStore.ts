import { readFileSync, writeFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

interface Tokens {
  access_token: string;
  refresh_token: string;
  expires_at: number;
  athlete_id: number;
}

const __dirname = dirname(fileURLToPath(import.meta.url));
const TOKEN_FILE = join(__dirname, '..', '..', 'tokens.json');

let tokens: Tokens | null = null;

function loadFromDisk(): void {
  try {
    const data = readFileSync(TOKEN_FILE, 'utf-8');
    tokens = JSON.parse(data);
  } catch {
    tokens = null;
  }
}

loadFromDisk();

export function getTokens(): Tokens | null {
  return tokens;
}

export function setTokens(newTokens: Tokens): void {
  tokens = newTokens;
  writeFileSync(TOKEN_FILE, JSON.stringify(newTokens, null, 2));
}

export function clearTokens(): void {
  tokens = null;
  try {
    writeFileSync(TOKEN_FILE, '');
  } catch {
    // ignore
  }
}

export function isTokenExpired(): boolean {
  if (!tokens) return true;
  const bufferSeconds = 300; // 5 minutes
  return Date.now() / 1000 >= tokens.expires_at - bufferSeconds;
}
