package systems

import (
	"strings"

	"github.com/istvzsig/eve-gate-watch/internal/data"
)

type Service struct {
	repo data.Repository
}

func NewService(repo data.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Search(query string) []string {
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		return []string{}
	}

	systems := s.repo.LoadSystems()

	results := make([]string, 0)

	for _, system := range systems {
		if strings.HasPrefix(strings.ToLower(system), query) {
			results = append(results, system)
		}
	}

	return results
}
