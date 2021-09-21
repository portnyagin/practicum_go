package service

type UserService struct {
	repository Repository
}

func NewUserService(repo Repository) *UserService {
	var s UserService
	s.repository = repo
	return &s
}

func (s *UserService) GetURLsByUser(userID string) ([]string, error) {
	// TODO: Проверить куку
	res, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
