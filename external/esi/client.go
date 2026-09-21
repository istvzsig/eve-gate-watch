package esi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ESIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewESIClient(baseURL string) *ESIClient {
	return &ESIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ESIClient) getJSON(path string, target any) error {
	req, err := http.NewRequest(
		http.MethodGet,
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return err
	}

	return c.doJSON(req, target)
}

func (c *ESIClient) postJSON(
	path string,
	body []byte,
	target any,
) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.doJSON(req, target)
}

func (c *ESIClient) doJSON(
	req *http.Request,
	target any,
) error {
	req.Header.Set(
		"User-Agent",
		"eve-gate-watch/0.1",
	)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(
			io.LimitReader(response.Body, 4096),
		)

		return fmt.Errorf(
			"ESI returned %s: %s",
			response.Status,
			strings.TrimSpace(string(body)),
		)
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf(
			"decode ESI response: %w",
			err,
		)
	}

	return nil
}
