package sde

type Region struct {
	ID   int64
	Name string
}

type Constellation struct {
	ID       int64
	Name     string
	RegionID int64
}

type SolarSystem struct {
	ID              int64
	Name            string
	ConstellationID int64
	RegionID        int64
	SecurityStatus  float64
	WormholeClassID *int
}
