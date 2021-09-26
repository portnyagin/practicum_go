package service

import (
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"github.com/portnyagin/practicum_go/internal/app/model"
)

type UserService struct {
	repository model.RepositoryV2
}

func NewUserService(repo model.RepositoryV2) *UserService {
	var s UserService
	s.repository = repo
	return &s
}

func (s *UserService) mapUserURLsDTO(src *model.UserURLs) (*dto.UserURLsDTO, error) {
	return &dto.UserURLsDTO{ShortURL: src.ShortURL, OriginalURL: src.OriginalURL}, nil
}

func (s *UserService) GetURLsByUser(userID string) ([]dto.UserURLsDTO, error) {
	if userID == "" {
		return nil, errors.New("user_id is empty")
	}
	resArr, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	var resDtoList []*dto.UserURLsDTO
	for _, rec := range resArr {
		dto, err := s.mapUserURLsDTO(&rec)
		if err != nil {
			return nil, errors.New("can't map result to UserURLsDTO")
			break
		}
		resDtoList = append(resDtoList, dto)
	}
	return []dto.UserURLsDTO{}, nil
}

func (s *UserService) Save(userID string, originalURL string, shortURL string) error {
	err := s.repository.Save(userID, shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Ping() bool {
	// TODO:
	res, _ := s.repository.Ping()
	return res
}
