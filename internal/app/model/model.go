package model

//type UserURLs struct {
//	shortURL    string
//	originalURL string
//}

type UserURLs struct {
	id          int
	userId      string
	shortURL    string
	originalURL string
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
