package intel

import (
	"fmt"
	"strings"

	"github.com/istvzsig/eve-gate-watch/external/zkill"
	"github.com/istvzsig/eve-gate-watch/internal/dscan"
	"github.com/istvzsig/eve-gate-watch/internal/pvp"
)

type Service struct {
	dScanner    dscan.Scanner
	zkillClient *zkill.Client
}

func NewService(dScanner dscan.Scanner, zkillClient *zkill.Client) *Service {
	return &Service{
		dScanner:    dScanner,
		zkillClient: zkillClient,
	}
}

func (s *Service) DScan(systemName string) (dscan.DScanResult, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return dscan.DScanResult{}, fmt.Errorf("system name is required")
	}

	return s.dScanner.DScan(systemName)
}

func (s *Service) RecentPVP(systemName string) (pvp.Result, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return pvp.Result{}, fmt.Errorf("system name is required")
	}

	killmails, err := s.zkillClient.RecentPVP(systemName, 3600)
	if err != nil {
		return pvp.Result{}, err
	}

	kills := make([]pvp.Kill, 0, len(killmails))

	for _, killmail := range killmails {
		kills = append(kills, pvp.Kill{
			ID:         killmail.KillmailID,
			OccurredAt: killmail.KillmailTime,
			ShipName:   killmail.ShipName,
			Value:      killmail.ZKB.TotalValue,
			Solo:       killmail.ZKB.Solo,
		})
	}

	return pvp.Result{
		SystemName: strings.ToTitle(systemName),
		Kills:      kills,
	}, nil
}
