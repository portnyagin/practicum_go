package service

import "github.com/stretchr/testify/mock"

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

func (r *RepositoryMock) FindByUser(key string) ([]RepoRecord, error) {
	return nil, nil
}

func mockEncode(str string) string {
	return str
}
