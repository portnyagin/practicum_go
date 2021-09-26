package dto

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
