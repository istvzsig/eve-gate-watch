package dscan

type Scanner interface {
	DScan(systemName string) (DScanResult, error)
}
