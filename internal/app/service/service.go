package service

import (
	"encoding/base64"
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/model"
	//"github.com/portnyagin/practicum_go/internal/app/handler"
)

type EncodeFunc func(str string) string

type ZipService struct {
	repository model.Repository
	encode     EncodeFunc
	baseURL    string
}

func NewZipService(repo model.Repository, baseURL string) *ZipService {
	var s ZipService
	s.repository = repo
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	s.baseURL = baseURL
	return &s
}

func (s *ZipService) ZipURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("URL is empty")
	}
	key := s.encode(url)
	err := s.repository.Save(key, url)
	if err != nil {
		return "", err
	}
	return s.baseURL + key, nil
}

func (s *ZipService) UnzipURL(key string) (string, error) {
	res, err := s.repository.Find(key)
	return res, err
}
