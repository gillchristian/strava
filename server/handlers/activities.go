package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cadence-server/store"
	"cadence-server/strava"
)

type ActivitiesHandler struct {
	Store  *store.TokenStore
	Strava *strava.Client
}

var runTypes = map[string]bool{
	"Run":        true,
	"TrailRun":   true,
	"VirtualRun": true,
}

func (h *ActivitiesHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	sessionToken := getSessionToken(r)
	if sessionToken == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not authenticated"})
		return
	}

	tokens, err := h.Store.GetTokensBySession(sessionToken)
	if err != nil {
		log.Printf("Activities token check error: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	if tokens == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not authenticated"})
		return
	}

	accessToken, refreshed, err := h.Strava.GetValidAccessToken(tokens)
	if err != nil {
		log.Printf("Activities access token error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get access token"})
		return
	}

	if refreshed != nil {
		if err := h.Store.UpdateTokens(*refreshed); err != nil {
			log.Printf("Activities token update error: %v", err)
		}
	}

	now := time.Now().Unix()
	thirtyDaysAgo := now - 30*24*60*60

	activities, err := h.Strava.FetchActivities(accessToken, thirtyDaysAgo, now)
	if err != nil {
		log.Printf("Activities fetch error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch activities"})
		return
	}

	// Filter to run types
	var runs []json.RawMessage
	for _, raw := range activities {
		var activity struct {
			Type      string `json:"type"`
			SportType string `json:"sport_type"`
		}
		if err := json.Unmarshal(raw, &activity); err != nil {
			continue
		}
		if runTypes[activity.Type] || runTypes[activity.SportType] {
			runs = append(runs, raw)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if runs == nil {
		runs = []json.RawMessage{}
	}
	json.NewEncoder(w).Encode(runs)
}
