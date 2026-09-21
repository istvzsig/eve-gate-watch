package pvp

import "time"

type Kill struct {
	ID         int64
	OccurredAt time.Time
	ShipName   string
	Value      float64
	Solo       bool
}

type Result struct {
	SystemName string
	Kills      []Kill
}
