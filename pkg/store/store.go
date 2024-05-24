package store

type Store interface {
	ListObjects(string) ([]string, error)
	DownloadObject(string, string) error
}
