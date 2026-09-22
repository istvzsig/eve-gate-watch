package systems

import (
	"context"
	"strings"

	"github.com/istvzsig/eve-gate-watch/internal/data"
)

type Service struct {
	repo data.Repository
}

func NewService(repo data.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(ctx context.Context, query string) ([]data.System, error) {
	query = strings.TrimSpace(query)

	if query == "" {
		return []data.System{}, nil
	}

	return s.repo.SearchSystems(ctx, query)
}
