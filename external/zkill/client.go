package zkill

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL        string
	httpClient     *http.Client
	systemResolver SystemResolver
	typeResolver   TypeResolver
}

func NewClient(
	baseURL string,
	systemResolver SystemResolver,
	typeResolver TypeResolver,
) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		systemResolver: systemResolver,
		typeResolver:   typeResolver,
	}
}

func (c *Client) RecentPVP(
	systemName string,
	pastSeconds int,
) ([]Killmail, error) {
	systemID, _, err := c.systemResolver.ResolveSystem(systemName)
	if err != nil {
		return nil, err
	}

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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read zKillboard response: %w", err)
	}

	var killmails []Killmail

	if err := json.Unmarshal(body, &killmails); err == nil {
		return killmails, nil
	}

	var apiError struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error != "" {
		return nil, fmt.Errorf("zKillboard API error: %s", apiError.Error)
	}

	return nil, fmt.Errorf(
		"decode zKillboard response: %s",
		strings.TrimSpace(string(body)),
	)
}
