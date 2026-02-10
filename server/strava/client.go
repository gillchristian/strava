package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"cadence-server/store"
)

const (
	stravaAPI   = "https://www.strava.com/api/v3"
	stravaOAuth = "https://www.strava.com/oauth"
)

type Client struct {
	ClientID     string
	ClientSecret string
	Store        *store.TokenStore
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Athlete      struct {
		ID int64 `json:"id"`
	} `json:"athlete"`
}

func (c *Client) ExchangeCodeForTokens(code string) error {
	resp, err := http.PostForm(stravaOAuth+"/token", url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token exchange failed: %d %s", resp.StatusCode, body)
	}

	var data tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("token exchange decode failed: %w", err)
	}

	return c.Store.SetTokens(store.Tokens{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresAt:    data.ExpiresAt,
		AthleteID:    data.Athlete.ID,
	})
}

func (c *Client) RefreshAccessToken() error {
	tokens, err := c.Store.GetTokens()
	if err != nil {
		return fmt.Errorf("failed to get tokens: %w", err)
	}
	if tokens == nil {
		return fmt.Errorf("no tokens to refresh")
	}

	resp, err := http.PostForm(stravaOAuth+"/token", url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {tokens.RefreshToken},
	})
	if err != nil {
		return fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed: %d %s", resp.StatusCode, body)
	}

	var data tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("token refresh decode failed: %w", err)
	}

	return c.Store.SetTokens(store.Tokens{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresAt:    data.ExpiresAt,
		AthleteID:    tokens.AthleteID,
	})
}

func (c *Client) GetValidAccessToken() (string, error) {
	if c.Store.IsTokenExpired() {
		if err := c.RefreshAccessToken(); err != nil {
			return "", err
		}
	}
	tokens, err := c.Store.GetTokens()
	if err != nil {
		return "", err
	}
	if tokens == nil {
		return "", fmt.Errorf("no tokens available")
	}
	return tokens.AccessToken, nil
}

func (c *Client) FetchActivities(after, before int64) ([]json.RawMessage, error) {
	accessToken, err := c.GetValidAccessToken()
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"after":    {strconv.FormatInt(after, 10)},
		"before":   {strconv.FormatInt(before, 10)},
		"per_page": {"100"},
	}

	req, err := http.NewRequest("GET", stravaAPI+"/athlete/activities?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("strava API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("strava API error: %d %s", resp.StatusCode, body)
	}

	var activities []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, fmt.Errorf("activities decode failed: %w", err)
	}

	return activities, nil
}
