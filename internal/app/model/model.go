package model

import "github.com/jackc/pgerrcode"

type UserURLs struct {
	ID          int
	UserID      string
	ShortURL    string
	OriginalURL string
}

type RepoRecord = string

type FileRepository interface {
	Find(key string) (RepoRecord, error)
	Save(key string, value string) error
	FindByUser(key string) ([]UserURLs, error)
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

type DBRepository interface {
	FindByUser(userID string) ([]UserURLs, error)
	FindByShort(shortURL string) (string, error)
	Save(userID string, originalURL string, shortURL string) error
	SaveBatch(data UserBatchURLs) error
	Ping() (bool, error)
}

type DatabaseError struct {
	Err  error
	Code string
}

func (t *DatabaseError) Error() string {
	return t.Err.Error()
}

var (
	UniqueViolation DatabaseError = DatabaseError{Code: pgerrcode.UniqueViolation}
)
