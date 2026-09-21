package data

import "sync"

type MemoryRepository struct {
	mu      sync.RWMutex
	systems []string
}

func NewMemoryRepository(systems []string) *MemoryRepository {
	return &MemoryRepository{
		systems: systems,
	}
}

func (r *MemoryRepository) LoadSystems() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	systems := make([]string, len(r.systems))
	copy(systems, r.systems)

	return systems
}
