package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	AthleteID    int64  `json:"athlete_id"`
}

type TokenStore struct {
	db *sql.DB
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewTokenStore(dbPath string) (*TokenStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Check if old single-user schema exists and migrate
	var hasIDCheck bool
	row := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='tokens'")
	var tableSql sql.NullString
	if err := row.Scan(&tableSql); err == nil && tableSql.Valid {
		// Old schema has CHECK(id = 1); drop and recreate
		if strings.Contains(tableSql.String, "CHECK") {
			hasIDCheck = true
		}
	}

	if hasIDCheck {
		if _, err := db.Exec("DROP TABLE tokens"); err != nil {
			return nil, err
		}
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			athlete_id INTEGER PRIMARY KEY,
			session_token TEXT NOT NULL UNIQUE,
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			expires_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &TokenStore{db: db}, nil
}

func (s *TokenStore) GetTokensBySession(sessionToken string) (*Tokens, error) {
	row := s.db.QueryRow(
		"SELECT access_token, refresh_token, expires_at, athlete_id FROM tokens WHERE session_token = ?",
		sessionToken,
	)

	var t Tokens
	err := row.Scan(&t.AccessToken, &t.RefreshToken, &t.ExpiresAt, &t.AthleteID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TokenStore) SetTokens(t Tokens, sessionToken string) error {
	_, err := s.db.Exec(
		`INSERT INTO tokens (athlete_id, session_token, access_token, refresh_token, expires_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(athlete_id) DO UPDATE SET
			session_token = excluded.session_token,
			access_token = excluded.access_token,
			refresh_token = excluded.refresh_token,
			expires_at = excluded.expires_at`,
		t.AthleteID, sessionToken, t.AccessToken, t.RefreshToken, t.ExpiresAt,
	)
	return err
}

func (s *TokenStore) UpdateTokens(t Tokens) error {
	_, err := s.db.Exec(
		"UPDATE tokens SET access_token = ?, refresh_token = ?, expires_at = ? WHERE athlete_id = ?",
		t.AccessToken, t.RefreshToken, t.ExpiresAt, t.AthleteID,
	)
	return err
}

func (s *TokenStore) ClearTokensBySession(sessionToken string) error {
	_, err := s.db.Exec("DELETE FROM tokens WHERE session_token = ?", sessionToken)
	return err
}

func IsTokenExpired(tokens *Tokens) bool {
	bufferSeconds := int64(300) // 5 minutes
	return time.Now().Unix() >= tokens.ExpiresAt-bufferSeconds
}
