package app

import (
	"encoding/base64"
	"errors"
)

type ZipService struct {
	store map[string]string
}

func NewZipService() *ZipService {
	var s ZipService
	s.store = make(map[string]string)
	return &s
}

func (s *ZipService) encode(str string) string {
	sha := base64.StdEncoding.EncodeToString([]byte(str))
	return sha
}

func (s *ZipService) ZipURL(url string) (string, error) {
	const baseURL string = "http://localhost:8080/"
	if url == "" {
		return "", errors.New("url is empty")
	}
	key := s.encode(url)
	s.store[key] = url
	return baseURL + key, nil
}

func (s *ZipService) UnzipURL(key string) (string, error) {
	if val, ok := s.store[key]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}
