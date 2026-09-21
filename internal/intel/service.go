package intel

import (
	"context"
	"fmt"
	"log"
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
) (IntelReport, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return IntelReport{}, fmt.Errorf("system name is required")
	}

	observations, err := s.repository.RecentBySystem(
		ctx,
		systemName,
		10,
	)
	if err != nil {
		return IntelReport{}, err
	}

	return IntelReport{
		System:       systemName,
		Observations: observations,
	}, nil
}

func (s *Service) RecordDscan(
	ctx context.Context,
	input dscan.DScanInput,
) error {
	systemName := strings.TrimSpace(input.SystemName)

	if systemName == "" {
		return fmt.Errorf("system name is required")
	}

	for _, ship := range input.Ships {
		observation := Observation{
			Pilot:      Pilot{Name: "Unknown Pilot"},
			ShipName:   ship.Name,
			SystemName: systemName,
			Source:     SourceDScan,
			Threat:     ThreatUnknown,
			ObservedAt: time.Now(),
		}

		if err := s.repository.CreateObservation(ctx, observation); err != nil {
			return err
		}

		log.Printf(
			"stored intel: system=%s ship=%s",
			systemName,
			ship.Name,
		)
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
