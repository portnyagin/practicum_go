package dto

import "errors"

type ShortenResponseDTO struct {
	Result string `json:"result"`
}

type ShortenRequestDTO struct {
	URL string `json:"url"`
}

type UserURLsDTO struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserBatchDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type UserBatchResultDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

var ErrDuplicateKey = errors.New("no rows in result set")
