package zkill

import "time"

type Killmail struct {
	KillmailID   int64     `json:"killmail_id"`
	KillmailTime time.Time `json:"killmail_time"`
	Victim       Victim    `json:"victim"`
	ZKB          ZKB       `json:"zkb"`
}

type Victim struct {
	ShipTypeID int64 `json:"ship_type_id"`
}

type ZKB struct {
	TotalValue float64 `json:"totalValue"`
	NPC        bool    `json:"npc"`
	Solo       bool    `json:"solo"`
}
