package dscan

type Scanner interface {
	DScan(systemName string) (DScanResult, error)
	ResolveSystem(systemName string) (int64, string, error)
	ResolveType(typeID int64) (string, error)
}
