package model

//type UserURLs struct {
//	ShortURL    string
//	OriginalURL string
//}

type UserURLs struct {
	Id          int
	UserID      string
	ShortURL    string
	OriginalURL string
}

type RepoRecord = string

type Repository interface {
	Find(key string) (RepoRecord, error)
	Save(key string, value string) error
	FindByUser(key string) ([]UserURLs, error)
	Ping() (bool, error)
}

type RepositoryV2 interface {
	FindByUser(key string) ([]UserURLs, error)
	Ping() (bool, error)
}
