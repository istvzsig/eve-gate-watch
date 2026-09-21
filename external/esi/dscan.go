package esi

import (
	"fmt"

	"github.com/istvzsig/eve-gate-watch/internal/dscan"
)

func (c *ESIClient) DScan(
	systemName string,
) (dscan.DScanResult, error) {
	systemID, resolvedName, err := c.ResolveSystem(systemName)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	system, err := c.getSystem(systemID)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	constellation, err := c.getConstellation(
		system.ConstellationID,
	)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	region, err := c.getRegion(
		constellation.RegionID,
	)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	jumps, err := c.getSystemJumps(systemID)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	kills, err := c.getSystemKills(systemID)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	return dscan.DScanResult{
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

func (c *ESIClient) getSystem(
	systemID int64,
) (systemResponse, error) {
	var response systemResponse

	err := c.getJSON(
		fmt.Sprintf("/universe/systems/%d/", systemID),
		&response,
	)

	return response, err
}

func (c *ESIClient) getConstellation(
	constellationID int64,
) (constellationResponse, error) {
	var response constellationResponse

	err := c.getJSON(
		fmt.Sprintf(
			"/universe/constellations/%d/",
			constellationID,
		),
		&response,
	)

	return response, err
}

func (c *ESIClient) getRegion(
	regionID int64,
) (regionResponse, error) {
	var response regionResponse

	err := c.getJSON(
		fmt.Sprintf("/universe/regions/%d/", regionID),
		&response,
	)

	return response, err
}

func (c *ESIClient) getSystemJumps(
	systemID int64,
) (systemJumpsResponse, error) {
	var response []systemJumpsResponse

	if err := c.getJSON(
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

func (c *ESIClient) getSystemKills(
	systemID int64,
) (systemKillsResponse, error) {
	var response []systemKillsResponse

	if err := c.getJSON(
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
