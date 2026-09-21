package zkill

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Killmail struct {
	KillmailID   int64     `json:"killmail_id"`
	KillmailTime time.Time `json:"killmail_time"`
	Victim       Victim    `json:"victim"`
	ZKB          ZKB       `json:"zkb"`
}

type Victim struct {
	ShipTypeID int64 `json:"ship_type_id"`
}

type ZKB struct {
	TotalValue float64 `json:"totalValue"`
	NPC        bool    `json:"npc"`
	Solo       bool    `json:"solo"`
}

func (c *Client) RecentPVP(
	systemID int64,
	pastSeconds int,
) ([]Killmail, error) {
	url := fmt.Sprintf(
		"%s/kills/solarSystemID/%d/pastSeconds/%d/",
		c.baseURL,
		systemID,
		pastSeconds,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"eve-gate-watch/0.1 github.com/iszsig/eve-gate-watch",
	)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))

		return nil, fmt.Errorf(
			"zKillboard returned %s: %s",
			response.Status,
			string(body),
		)
	}

	var killmails []Killmail

	if err := json.NewDecoder(response.Body).Decode(&killmails); err != nil {
		return nil, fmt.Errorf("decode zKillboard response: %w", err)
	}

	return killmails, nil
}
