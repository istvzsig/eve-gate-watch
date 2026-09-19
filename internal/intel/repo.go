package intel

import (
	"context"
	"sync"
)

type Repository interface {
	CreateObservation(ctx context.Context, observation Observation) error
	RecentBySystem(ctx context.Context, systemName string, limit int) ([]Observation, error)
}

type MemoryRepository struct {
	mu           sync.RWMutex
	observations []Observation
}

var _ Repository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) CreateObservation(
	ctx context.Context,
	observation Observation,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.observations = append(r.observations, observation)

	return nil
}

func (r *MemoryRepository) RecentBySystem(
	ctx context.Context,
	systemName string,
	limit int,
) ([]Observation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Observation

	for i := len(r.observations) - 1; i >= 0 && len(result) < limit; i-- {
		observation := r.observations[i]

		if observation.SystemName != systemName {
			continue
		}

		result = append(result, observation)
	}

	return result, nil
}
