package dscan

import (
	"fmt"
	"strings"

	"github.com/iszsig/eve-gate-watch/internal/ship"
)

type FakeDScan struct {
	Ships []ship.Ship
}

func NewFakeDScan(ships []ship.Ship) *FakeDScan {
	s := append([]ship.Ship(nil), ships...)

	return &FakeDScan{
		Ships: s,
	}
}

var _ Scanner = (*FakeDScan)(nil)

func (dscan *FakeDScan) DScan(systemName string) (DScanResult, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return DScanResult{}, fmt.Errorf("system name is required")
	}

	s := append([]ship.Ship(nil), dscan.Ships...)

	return DScanResult{
		SystemName: systemName,
		Ships:      s,
	}, nil
}
