package pvp

type SystemResolver interface {
	ResolveSystem(systemName string) (int64, string, error)
}
