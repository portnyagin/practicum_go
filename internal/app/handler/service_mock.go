package handler

import (
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (s *UserServiceMock) GetURLsByUser(userID string) ([]dto.UserURLsDTO, error) {
	args := s.Called(userID)
	var res []dto.UserURLsDTO
	res = append(res, dto.UserURLsDTO{OriginalURL: args.String(0), ShortURL: args.String(0)})
	return res, args.Error(1)
}

func (s *UserServiceMock) Ping() bool {
	args := s.Called()
	// TODO:
	return args.Bool(0)
}

func (s *UserServiceMock) Save(userID string, originalURL string, shortURL string) error {
	args := s.Called(userID, originalURL, shortURL)
	if originalURL == "bad_URL" {
		return args.Error(0)
	} else {
		return nil
	}
}

func (s *UserServiceMock) SaveBatch(userID string, srcDTO []dto.UserBatchDTO) ([]dto.UserBatchResultDTO, error) {
	args := s.Called(userID, srcDTO)
	var res []dto.UserBatchResultDTO
	res = append(res, dto.UserBatchResultDTO{CorrelationID: args.String(0), ShortURL: args.String(1)})
	return res, args.Error(2)
}

func (s *UserServiceMock) GetURLByShort(shortURL string) (string, error) {
	args := s.Called(shortURL)
	return args.String(0), args.Error(1)
}

func (s *UserServiceMock) ZipURL(url string) (string, error) {
	args := s.Called(url)
	return args.String(0), args.Error(1)
}

//----------------------------------------------------------------

type CryptoServiceMock struct {
	mock.Mock
}

func (s *CryptoServiceMock) Validate(token string) (bool, string) {
	args := s.Called(token)
	return args.Bool(0), token
}

func (s *CryptoServiceMock) GetNewUserToken() (string, string, error) {
	args := s.Called()
	return args.String(0), args.String(1), args.Error(2)
}
