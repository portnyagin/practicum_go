package app

import (
	"encoding/base64"
	"errors"
)

type EncodeFunc func(str string) string

type ZipService struct {
	repository Repository
	encode     EncodeFunc
}

type Repository interface {
	Find(key string) (string, error)
	Save(key string, value string) error
}

func NewZipService(repo Repository) *ZipService {
	var s ZipService
	s.repository = repo
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return &s
}

func (s *ZipService) ZipURL(url string) (string, error) {
	const baseURL string = "http://localhost:8080/"
	if url == "" {
		return "", errors.New("URL is empty")
	}
	key := s.encode(url)
	s.repository.Save(key, url)
	return baseURL + key, nil
}

func (s *ZipService) UnzipURL(key string) (string, error) {
	res, err := s.repository.Find(key)
	return res, err
}
