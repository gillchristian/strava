package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"cadence-server/strava"
	"cadence-server/store"
)

type AuthHandler struct {
	Store       *store.TokenStore
	Strava      *strava.Client
	ClientID    string
	APIBaseURL  string
	FrontendURL string
}

func (h *AuthHandler) StravaRedirect(w http.ResponseWriter, r *http.Request) {
	redirectURI := h.APIBaseURL + "/auth/callback"
	u := fmt.Sprintf(
		"https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=activity:read_all&approval_prompt=auto",
		h.ClientID,
		url.QueryEscape(redirectURI),
	)
	http.Redirect(w, r, u, http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	if err := h.Strava.ExchangeCodeForTokens(code); err != nil {
		log.Printf("OAuth callback error: %v", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, h.FrontendURL+"/?auth=success", http.StatusFound)
}

func (h *AuthHandler) Status(w http.ResponseWriter, r *http.Request) {
	tokens, err := h.Store.GetTokens()
	if err != nil {
		log.Printf("Status check error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if tokens == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false,
			"athleteId":     nil,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"athleteId":     tokens.AthleteID,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.ClearTokens(); err != nil {
		log.Printf("Logout error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
