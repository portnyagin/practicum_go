package model

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
)

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
	FindByUser(ctx context.Context, userID string) ([]UserURLs, error)
	FindByShort(ctx context.Context, userID string, shortURL string) (string, error)
	Save(ctx context.Context, userID string, originalURL string, shortURL string) error
	SaveBatch(ctx context.Context, data UserBatchURLs) error
	Ping(ctx context.Context) (bool, error)
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
	NoRowFound      DatabaseError = DatabaseError{Err: errors.New("no rows in result set")}
)

type DeleteRepository interface {
	BatchDelete(ctx context.Context, userID string, URLList []BatchDeleteURL) error
}

type BatchDeleteURL = string
