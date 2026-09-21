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

	systemID, _, err := s.dScanner.ResolveSystem(systemName)
	if err != nil {
		return pvp.Result{}, err
	}

	killmails, err := s.zkillClient.RecentPVP(systemID, 3600)
	if err != nil {
		return pvp.Result{}, err
	}

	kills := make([]pvp.Kill, 0, len(killmails))

	for _, killmail := range killmails {
		shipName, _ := s.dScanner.ResolveType(killmail.Victim.ShipTypeID)

		kills = append(kills, pvp.Kill{
			ID:         killmail.KillmailID,
			OccurredAt: killmail.KillmailTime,
			ShipName:   shipName,
			Value:      killmail.ZKB.TotalValue,
			Solo:       killmail.ZKB.Solo,
		})
	}

	return pvp.Result{
		SystemName: systemName,
		Kills:      kills,
	}, nil
}
