package app

type ShortenResponseDTO struct {
	Result string `json:"result"`
}

type ShortenRequestDTO struct {
	URL string `json:"url"`
}
