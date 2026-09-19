package intel

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) RecentIntel(
	ctx context.Context,
	systemName string,
) (string, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return "", fmt.Errorf("system name is required")
	}

	observations, err := s.repository.RecentBySystem(
		ctx,
		systemName,
		10,
	)
	if err != nil {
		return "", err
	}

	if len(observations) == 0 {
		return fmt.Sprintf(
			"🚨 %s INTEL\n\nNo recent observations.",
			strings.ToUpper(systemName),
		), nil
	}

	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"🚨 %s INTEL\n\n",
		strings.ToUpper(systemName),
	)

	for _, observation := range observations {
		fmt.Fprintf(
			&builder,
			"• %s — %s — %s\n",
			observation.PilotName,
			observation.ShipName,
			observation.Threat,
		)
	}

	return builder.String(), nil
}
