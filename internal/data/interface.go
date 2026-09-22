package data

import "context"

type Repository interface {
	SearchSystems(ctx context.Context, query string) ([]System, error)
}

type System struct {
	ID              int64
	Name            string
	ConstellationID int64
	RegionID        int64
	SecurityStatus  float64
	WormholeClassID *int
}
