package service

import (
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

func (s *UserService) GetURLsByUser(userID string) ([]dto.UserURLsDTO, error) {
	_, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	return []dto.UserURLsDTO{}, nil
}

func (s *UserService) Ping() bool {
	// TODO:
	res, _ := s.repository.Ping()
	return res
}
