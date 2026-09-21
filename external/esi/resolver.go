package esi

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *ESIClient) ResolveSystem(
	systemName string,
) (int64, string, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return 0, "", fmt.Errorf("system name is required")
	}

	body, err := json.Marshal([]string{systemName})
	if err != nil {
		return 0, "", err
	}

	var response resolveResponse

	if err := c.postJSON(
		"/universe/ids/",
		body,
		&response,
	); err != nil {
		return 0, "", err
	}

	if len(response.Systems) == 0 {
		return 0, "", fmt.Errorf(
			"solar system %q not found",
			systemName,
		)
	}

	system := response.Systems[0]

	return system.ID, system.Name, nil
}

func (c *ESIClient) ResolveType(typeID int64) (string, error) {
	if typeID <= 0 {
		return "", fmt.Errorf("type id is required")
	}

	var response typeResponse

	if err := c.getJSON(
		fmt.Sprintf("/universe/types/%d/", typeID),
		&response,
	); err != nil {
		return "", err
	}

	if response.Name == "" {
		return "", fmt.Errorf(
			"type %d has no name",
			typeID,
		)
	}

	return response.Name, nil
}
