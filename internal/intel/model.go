package intel

import "time"

type Threat string

const (
	ThreatFriendly Threat = "friendly"
	ThreatNeutral  Threat = "neutral"
	ThreatUnknown  Threat = "unknown"
	ThreatHostile  Threat = "hostile"
)

type Observation struct {
	ID         string
	PilotName  string
	SystemName string
	ShipName   string
	Threat     Threat
	Source     string
	ObservedAt time.Time
}
