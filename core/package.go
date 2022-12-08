package core

type Package interface {
	CreateWrapper() (Wrapper, error)
	Manifest(validation bool) (any, error)
}
