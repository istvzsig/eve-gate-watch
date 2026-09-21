package dscan

import "github.com/iszsig/eve-gate-watch/internal/ship"

type DScanInput struct {
	SystemName string
	Ships      []ship.Ship
}

type DScanResult struct {
	SystemName string
	Ships      []ship.Ship
}
