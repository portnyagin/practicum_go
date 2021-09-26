package model

//type UserURLs struct {
//	short_url    string
//	original_url string
//}

type UserURLs struct {
	id           int
	user_id      string
	short_url    string
	original_url string
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
