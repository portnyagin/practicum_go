package repository

import (
	"github.com/portnyagin/practicum_go/internal/app/model"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Find(key string) (string, error) {
	args := r.Called(key)
	return args.String(0), nil
}

func (r *RepositoryMock) Save(key string, value string) error {
	return nil
}

func (r *RepositoryMock) FindByUser(key string) ([]model.UserURLs, error) {
	return nil, nil
}

func MockEncode(str string) string {
	return str
}

func (r *RepositoryMock) Ping() (bool, error) {
	args := r.Called()
	return args.Bool(0), args.Error(1)
}
