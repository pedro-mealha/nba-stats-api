package nba

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	API interface {
		GetScoreboard(context.Context, GetScoreboardCommand) (ScoreboardData, error)
		GetBoxscore(context.Context, GetBoxscoreCommand) (BoxscoreData, error)
	}

	// Client is the NBA API client
	Client struct {
		baseURL     string
		cdnURL      string
		wnbaCdnURL  string
		client      httpClient
	}
)

// New creates a new instance of Client
func New(client httpClient, baseURL, cdnURL, wnbaCdnURL string) *Client {
	return &Client{baseURL: baseURL, cdnURL: cdnURL, client: client, wnbaCdnURL: wnbaCdnURL}
}

// GetScoreboard get scoreboard for a specific day
func (c *Client) GetScoreboard(ctx context.Context, cmd GetScoreboardCommand) (ScoreboardData, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/stats/scoreboardv3", c.baseURL),
		nil,
	)
	if err != nil {
		return ScoreboardData{}, err
	}

	q := req.URL.Query()
	q.Set("LeagueID", string(cmd.LeagueID))
	q.Set("GameDate", cmd.Date)
	req.URL.RawQuery = q.Encode()

	resp, err := c.doRequest(req)
	if err != nil {
		return ScoreboardData{}, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	var s ScoreboardData
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return ScoreboardData{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return s, nil
}

// GetScoreboard get scoreboard for a specific day
func (c *Client) GetBoxscore(ctx context.Context, cmd GetBoxscoreCommand) (BoxscoreData, error) {
	cdnURL := c.cdnURL
	if cmd.LeagueID == WNBA {
		cdnURL = c.wnbaCdnURL
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/static/json/liveData/boxscore/boxscore_%s.json", cdnURL, cmd.GameID),
		nil,
	)
	if err != nil {
		return BoxscoreData{}, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return BoxscoreData{}, fmt.Errorf("failed to request nba api: %w", err)
	}

	defer resp.Body.Close()

	var b BoxscoreData
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return BoxscoreData{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return b, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	// Ignore it. Skip it.
	req.Header.Set("Referer", c.baseURL)
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.2")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("failed to get scoreboard with status code %d", resp.StatusCode)
	}

	return resp, nil
}
