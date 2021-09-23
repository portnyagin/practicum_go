package service

type UserService struct {
	repository Repository
}

func NewUserService(repo Repository) *UserService {
	var s UserService
	s.repository = repo
	return &s
}

func (s *UserService) GetURLsByUser(userID string) ([]UserURLs, error) {
	_, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	return []UserURLs{UserURLs{}}, nil
}

func (s *UserService) Ping() bool {
	// TODO:
	return false
}