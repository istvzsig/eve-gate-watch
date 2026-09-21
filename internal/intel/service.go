package intel

import (
	"fmt"
	"strings"

	"github.com/istvzsig/eve-gate-watch/internal/dscan"
)

type Service struct {
	dScanner dscan.Scanner
}

func NewService(dScanner dscan.Scanner) *Service {
	return &Service{
		dScanner: dScanner,
	}
}

func (s *Service) DScan(
	systemName string,
) (dscan.DScanResult, error) {
	systemName = strings.TrimSpace(systemName)

	if systemName == "" {
		return dscan.DScanResult{}, fmt.Errorf("system name is required")
	}

	return s.dScanner.DScan(systemName)
}
