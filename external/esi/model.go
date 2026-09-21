package esi

type constellationResponse struct {
	Name     string `json:"name"`
	RegionID int64  `json:"region_id"`
}

type resolveResponse struct {
	Systems []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"systems"`
}

type regionResponse struct {
	Name string `json:"name"`
}

type systemResponse struct {
	Name            string     `json:"name"`
	ConstellationID int64      `json:"constellation_id"`
	SecurityStatus  float64    `json:"security_status"`
	Planets         []struct{} `json:"planets"`
	Stargates       []int64    `json:"stargates"`
	Stations        []int64    `json:"stations"`
}

type systemJumpsResponse struct {
	SystemID  int64 `json:"system_id"`
	ShipJumps int   `json:"ship_jumps"`
}

type systemKillsResponse struct {
	SystemID  int64 `json:"system_id"`
	ShipKills int   `json:"ship_kills"`
	PodKills  int   `json:"pod_kills"`
	NPCKills  int   `json:"npc_kills"`
}

type typeResponse struct {
	Name string `json:"name"`
}
