package handler

import (
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) ZipURL(url string) (string, error) {
	args := s.Called(url)
	return args.String(0), args.Error(1)
}

func (s *ServiceMock) ZipURLv2(url string) (*ShortenResponseDTO, error) {
	args := s.Called(url)
	return &ShortenResponseDTO{Result: args.String(0)}, args.Error(1)
}

func (s *ServiceMock) UnzipURL(key string) (string, error) {
	args := s.Called(key)
	return args.String(0), args.Error(1)
}

//----------------------------------------------------------------
type UserServiceMock struct {
	mock.Mock
}

func (s *UserServiceMock) GetURLsByUser(userID string) ([]string, error) {
	args := s.Called(userID)
	return []string{args.String(0)}, args.Error(1)
}

//----------------------------------------------------------------

type CryptoServiceMock struct {
	mock.Mock
}

func (s *CryptoServiceMock) Validate(token string) (bool, string) {
	args := s.Called(token)
	return args.Bool(0), args.String(1)
}

func (s *CryptoServiceMock) GetNewUserToken() (string, string, error) {
	args := s.Called()
	return args.String(0), args.String(1), args.Error(2)
}