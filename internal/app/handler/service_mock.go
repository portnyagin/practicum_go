package handler

import (
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) ZipURL(url string) (string, error) {
	args := s.Called(url)
	return args.String(0), args.Error(1)
}

func (s *ServiceMock) ZipURLv2(url string) (*dto.ShortenResponseDTO, error) {
	args := s.Called(url)
	return &dto.ShortenResponseDTO{Result: args.String(0)}, args.Error(1)
}

func (s *ServiceMock) UnzipURL(key string) (string, error) {
	args := s.Called(key)
	return args.String(0), args.Error(1)
}

//----------------------------------------------------------------
type UserServiceMock struct {
	mock.Mock
}

func (s *UserServiceMock) GetURLsByUser(userID string) ([]dto.UserURLsDTO, error) {
	args := s.Called(userID)
	// TODO:
	return []dto.UserURLsDTO{}, args.Error(1)
}

func (s *UserServiceMock) Ping() bool {
	args := s.Called()
	// TODO:
	return args.Bool(0)
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
