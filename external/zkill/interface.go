package zkill

type SystemResolver interface {
	ResolveSystem(systemName string) (int64, string, error)
}

type TypeResolver interface {
	ResolveType(typeID int64) (string, error)
}
