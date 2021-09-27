package model

//type UserURLs struct {
//	ShortURL    string
//	OriginalURL string
//}

type UserURLs struct {
	ID          int
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

type Element struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
}

type UserBatchURLs struct {
	UserID string
	List   []Element
}

type RepositoryV2 interface {
	FindByUser(userID string) ([]UserURLs, error)
	Save(userID string, shortURL string, originalURL string) error
	SaveBatch(UserBatchURLs) error
	ReadBatch(userID string) (*UserBatchURLs, error)
	Ping() (bool, error)
}
