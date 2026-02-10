package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cadence-server/strava"
	"cadence-server/store"
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
	tokens, err := h.Store.GetTokens()
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

	now := time.Now().Unix()
	thirtyDaysAgo := now - 30*24*60*60

	activities, err := h.Strava.FetchActivities(thirtyDaysAgo, now)
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
