package store

import (
	"database/sql"
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

func NewTokenStore(dbPath string) (*TokenStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER PRIMARY KEY CHECK(id = 1),
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			expires_at INTEGER NOT NULL,
			athlete_id INTEGER NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &TokenStore{db: db}, nil
}

func (s *TokenStore) GetTokens() (*Tokens, error) {
	row := s.db.QueryRow("SELECT access_token, refresh_token, expires_at, athlete_id FROM tokens WHERE id = 1")

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

func (s *TokenStore) SetTokens(t Tokens) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO tokens (id, access_token, refresh_token, expires_at, athlete_id) VALUES (1, ?, ?, ?, ?)",
		t.AccessToken, t.RefreshToken, t.ExpiresAt, t.AthleteID,
	)
	return err
}

func (s *TokenStore) ClearTokens() error {
	_, err := s.db.Exec("DELETE FROM tokens WHERE id = 1")
	return err
}

func (s *TokenStore) IsTokenExpired() bool {
	tokens, err := s.GetTokens()
	if err != nil || tokens == nil {
		return true
	}
	bufferSeconds := int64(300) // 5 minutes
	return time.Now().Unix() >= tokens.ExpiresAt-bufferSeconds
}
