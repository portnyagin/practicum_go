package service

import (
	"context"
	"github.com/portnyagin/practicum_go/internal/app/model"
	"github.com/stretchr/testify/mock"
)

type FileRepositoryMock struct {
	mock.Mock
}

func (r *FileRepositoryMock) Find(key string) (string, error) {
	args := r.Called(key)
	return args.String(0), nil
}

func (r *FileRepositoryMock) Save(key string, value string) error {
	return nil
}

func (r *FileRepositoryMock) FindByUser(key string) ([]model.UserURLs, error) {
	return nil, nil
}

func MockEncode(str string) string {
	return str
}

/*
//////////////////////////////////////////////////////////////////////////////////////////////////
*/
type DBRepositoryMock struct {
	mock.Mock
}

func (r *DBRepositoryMock) FindByUser(ctx context.Context, userID string) ([]model.UserURLs, error) {
	args := r.Called(userID)
	res := model.UserURLs{ID: 1, UserID: userID, OriginalURL: args.String(0), ShortURL: args.String(1)}
	return []model.UserURLs{res}, args.Error(2)
}

func (r *DBRepositoryMock) FindByShort(ctx context.Context, userID string, shortURL string) (string, error) {
	args := r.Called(userID, shortURL)
	return args.String(0), args.Error(1)
}

func (r *DBRepositoryMock) Save(ctx context.Context, userID string, originalURL string, shortURL string) error {
	args := r.Called(userID, originalURL, shortURL)
	return args.Error(0)
}

func (r *DBRepositoryMock) SaveBatch(ctx context.Context, data model.UserBatchURLs) error {
	args := r.Called(data)
	return args.Error(0)
}

func (r *DBRepositoryMock) Ping(ctx context.Context) (bool, error) {
	args := r.Called()
	return args.Bool(0), args.Error(1)
}
