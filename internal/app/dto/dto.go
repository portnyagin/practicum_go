package dto

type ShortenResponseDTO struct {
	Result string `json:"result"`
}

type ShortenRequestDTO struct {
	URL string `json:"url"`
}

type UserURLsDTO struct {
	short_url    string `json:"short___url"`
	original_url string `json:"original___url"`
}
