package store

type Store interface {
	ListObjects() ([]string, error)
}
