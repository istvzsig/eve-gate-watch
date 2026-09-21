package dscan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// POST /universe/ids/
// GET  /universe/systems/{system_id}/
// GET  /universe/system_jumps/
// GET  /universe/system_kills/

type ESIScanner struct {
	baseURL    string
	httpClient *http.Client
}

func NewESIScanner(url string) *ESIScanner {
	return &ESIScanner{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

var _ Scanner = (*ESIScanner)(nil)

type resolveResponse struct {
	Systems []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"systems"`
}

type systemResponse struct {
	Name            string     `json:"name"`
	ConstellationID int64      `json:"constellation_id"`
	SecurityStatus  float64    `json:"security_status"`
	Planets         []struct{} `json:"planets"`
	Stargates       []int64    `json:"stargates"`
	Stations        []int64    `json:"stations"`
}

type constellationResponse struct {
	Name     string `json:"name"`
	RegionID int64  `json:"region_id"`
}

type regionResponse struct {
	Name string `json:"name"`
}

type systemJumpsResponse struct {
	SystemID  int64 `json:"system_id"`
	ShipJumps int   `json:"ship_jumps"`
}

type systemKillsResponse struct {
	SystemID  int64 `json:"system_id"`
	ShipKills int   `json:"ship_kills"`
	PodKills  int   `json:"pod_kills"`
	NPCKills  int   `json:"npc_kills"`
}

func (s *ESIScanner) DScan(systemName string) (DScanResult, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return DScanResult{}, fmt.Errorf("system name is required")
	}

	systemID, resolvedName, err := s.resolveSystem(systemName)
	if err != nil {
		return DScanResult{}, err
	}

	system, err := s.getSystem(systemID)
	if err != nil {
		return DScanResult{}, err
	}

	constellation, err := s.getConstellation(system.ConstellationID)
	if err != nil {
		return DScanResult{}, err
	}

	region, err := s.getRegion(constellation.RegionID)
	if err != nil {
		return DScanResult{}, err
	}

	jumps, err := s.getSystemJumps(systemID)
	if err != nil {
		return DScanResult{}, err
	}

	kills, err := s.getSystemKills(systemID)
	if err != nil {
		return DScanResult{}, err
	}

	return DScanResult{
		SystemName:        resolvedName,
		RegionName:        region.Name,
		ConstellationName: constellation.Name,
		SecurityStatus:    system.SecurityStatus,
		Planets:           len(system.Planets),
		Stargates:         len(system.Stargates),
		Stations:          len(system.Stations),
		ShipJumps:         jumps.ShipJumps,
		ShipKills:         kills.ShipKills,
		PodKills:          kills.PodKills,
		NPCKills:          kills.NPCKills,
	}, nil
}

func (s *ESIScanner) resolveSystem(
	systemName string,
) (int64, string, error) {
	body, err := json.Marshal([]string{systemName})
	if err != nil {
		return 0, "", err
	}

	var response resolveResponse

	if err := s.postJSON(
		"/universe/ids/",
		body,
		&response,
	); err != nil {
		return 0, "", err
	}

	if len(response.Systems) == 0 {
		return 0, "", fmt.Errorf("solar system %q not found", systemName)
	}

	return response.Systems[0].ID, response.Systems[0].Name, nil
}

func (s *ESIScanner) getSystem(
	systemID int64,
) (systemResponse, error) {
	var response systemResponse

	err := s.getJSON(
		fmt.Sprintf("/universe/systems/%d/", systemID),
		&response,
	)

	return response, err
}

func (s *ESIScanner) getSystemJumps(
	systemID int64,
) (systemJumpsResponse, error) {
	var response []systemJumpsResponse

	if err := s.getJSON(
		"/universe/system_jumps/",
		&response,
	); err != nil {
		return systemJumpsResponse{}, err
	}

	for _, item := range response {
		if item.SystemID == systemID {
			return item, nil
		}
	}

	return systemJumpsResponse{
		SystemID:  systemID,
		ShipJumps: 0,
	}, nil
}

func (s *ESIScanner) getSystemKills(
	systemID int64,
) (systemKillsResponse, error) {
	var response []systemKillsResponse

	if err := s.getJSON(
		"/universe/system_kills/",
		&response,
	); err != nil {
		return systemKillsResponse{}, err
	}

	for _, item := range response {
		if item.SystemID == systemID {
			return item, nil
		}
	}

	return systemKillsResponse{
		SystemID:  systemID,
		ShipKills: 0,
		PodKills:  0,
		NPCKills:  0,
	}, nil
}

func (s *ESIScanner) getConstellation(
	constellationID int64,
) (constellationResponse, error) {
	var response constellationResponse

	err := s.getJSON(
		fmt.Sprintf("/universe/constellations/%d/", constellationID),
		&response,
	)

	return response, err
}

func (s *ESIScanner) getRegion(
	regionID int64,
) (regionResponse, error) {
	var response regionResponse

	err := s.getJSON(
		fmt.Sprintf("/universe/regions/%d/", regionID),
		&response,
	)

	return response, err
}

func (s *ESIScanner) getJSON(
	path string,
	target any,
) error {
	req, err := http.NewRequest(
		http.MethodGet,
		s.baseURL+path,
		nil,
	)
	if err != nil {
		return err
	}

	return s.doJSON(req, target)
}

func (s *ESIScanner) postJSON(
	path string,
	body []byte,
	target any,
) error {
	req, err := http.NewRequest(
		http.MethodPost,
		s.baseURL+path,
		strings.NewReader(string(body)),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	return s.doJSON(req, target)
}

func (s *ESIScanner) doJSON(
	req *http.Request,
	target any,
) error {
	req.Header.Set("User-Agent", "eve-gate-watch/0.1")

	response, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))

		return fmt.Errorf(
			"ESI returned %s: %s",
			response.Status,
			strings.TrimSpace(string(body)),
		)
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decode ESI response: %w", err)
	}

	return nil
}
