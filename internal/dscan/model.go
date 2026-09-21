package dscan

type DScanResult struct {
	SystemName        string  `json:"system"`
	RegionName        string  `json:"region"`
	ConstellationName string  `json:"constellation"`
	SecurityStatus    float64 `json:"security_status"`
	Planets           int     `json:"planets"`
	Stargates         int     `json:"stargates"`
	Stations          int     `json:"stations"`
	ShipJumps         int     `json:"ship_jumps"`
	ShipKills         int     `json:"ship_kills"`
	PodKills          int     `json:"pod_kills"`
	NPCKills          int     `json:"npc_kills"`
}
