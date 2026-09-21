package intel

import (
	"time"
)

type Threat string

const (
	ThreatFriendly Threat = "friendly"
	ThreatNeutral  Threat = "neutral"
	ThreatUnknown  Threat = "unknown"
	ThreatHostile  Threat = "hostile"
)

type Pilot struct {
	Name string
}

type ObservationSource string

const (
	SourceDScan    ObservationSource = "dscan"
	SourceKillmail ObservationSource = "killmail"
	SourceManual   ObservationSource = "manual"
	SourceESI      ObservationSource = "esi"
)

type Observation struct {
	ID         string
	Pilot      Pilot
	SystemName string
	ShipName   string
	Threat     Threat
	Source     ObservationSource
	ObservedAt time.Time
}
