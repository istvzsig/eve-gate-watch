package source

import "github.com/iszsig/eve-gate-watch/internal/dscan"

type ObservationSource interface {
	DScan(systemName string) (dscan.DScanResult, error)
}
