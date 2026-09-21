package intel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/iszsig/eve-gate-watch/internal/dscan"
)

type Service struct {
	repository Repository
	dScanner   dscan.Scanner
}

func NewService(repository Repository, dScanner dscan.Scanner) *Service {
	return &Service{
		repository: repository,
		dScanner:   dScanner,
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
			"🧑‍🚀 %s — %s — %s\n",
			observation.Pilot,
			observation.ShipName,
			observation.Threat,
		)
	}

	return builder.String(), nil
}

func (s *Service) RecordDscan(
	ctx context.Context,
	input dscan.DScanInput,
) error {
	// Local
	// Killmails
	// Manual reports
	// ESI/public sources

	systemName := strings.TrimSpace(input.SystemName)

	if systemName == "" {
		return fmt.Errorf("system name is required")
	}

	for _, ship := range input.Ships {
		observation := Observation{
			Pilot:      Pilot{Name: "Unknown Pilot"},
			ShipName:   ship.Name,
			SystemName: systemName,
			Source:     "dscan",
			Threat:     ThreatUnknown,
			ObservedAt: time.Now(),
		}
		if err := s.repository.CreateObservation(ctx, observation); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DScan(
	ctx context.Context,
	systemName string,
) (dscan.DScanResult, error) {
	result, err := s.dScanner.DScan(systemName)
	if err != nil {
		return dscan.DScanResult{}, err
	}

	input := dscan.DScanInput{
		SystemName: result.SystemName,
		Ships:      result.Ships,
	}

	if err := s.RecordDscan(ctx, input); err != nil {
		return dscan.DScanResult{}, err
	}

	return result, nil
}
