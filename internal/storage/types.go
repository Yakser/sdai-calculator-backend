package storage

type Storage interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
}
