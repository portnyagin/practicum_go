package service

import (
	"encoding/base64"
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"github.com/portnyagin/practicum_go/internal/app/model"
)

type UserService struct {
	repository model.RepositoryV2
	encode     EncodeFunc
	baseURL    string
}

func NewUserService(repo model.RepositoryV2, baseURL string) *UserService {
	var s UserService
	s.repository = repo
	s.encode = func(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	s.baseURL = baseURL
	return &s
}

func (s *UserService) zipURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("URL is empty")
	}
	key := s.encode(url)
	return s.baseURL + key, nil
}

//********** Mappers *****************************************************************/
func (s *UserService) mapUserURLsDTO(src *model.UserURLs) (*dto.UserURLsDTO, error) {
	return &dto.UserURLsDTO{ShortURL: src.ShortURL, OriginalURL: src.OriginalURL}, nil
}

//********** Mappers *****************************************************************/

func (s *UserService) GetURLsByUser(userID string) ([]dto.UserURLsDTO, error) {
	if userID == "" {
		return nil, errors.New("user_id is empty")
	}
	resArr, err := s.repository.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	var resDtoList []dto.UserURLsDTO
	for _, rec := range resArr {
		d, err := s.mapUserURLsDTO(&rec)
		if err != nil {
			return nil, errors.New("can't map result to UserURLsDTO")
		}
		resDtoList = append(resDtoList, *d)
	}
	return resDtoList, nil
}

func (s *UserService) Save(userID string, originalURL string, shortURL string) error {
	err := s.repository.Save(userID, shortURL, originalURL)
	if errors.Is(err, &model.UniqueViolation) {
		return dto.ErrDuplicateKey
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) SaveBatch(userID string, srcDTO []dto.UserBatchDTO) ([]dto.UserBatchResultDTO, error) {
	var (
		res    model.UserBatchURLs
		err    error
		resDTO []dto.UserBatchResultDTO
	)
	res.UserID = userID
	for _, obj := range srcDTO {
		var e model.Element
		e.CorrelationID = obj.CorrelationID
		e.OriginalURL = obj.OriginalURL
		e.ShortURL, err = s.zipURL(obj.OriginalURL)
		if err != nil {
			return nil, err
		}
		res.List = append(res.List, e)
		resDTO = append(resDTO, dto.UserBatchResultDTO{CorrelationID: obj.CorrelationID, ShortURL: e.ShortURL})
	}
	err = s.repository.SaveBatch(res)
	if errors.Is(err, &model.UniqueViolation) {
		return nil, dto.ErrDuplicateKey
	}
	if err != nil {
		return nil, err
	}
	return resDTO, nil
}

func (s *UserService) Ping() bool {
	res, _ := s.repository.Ping()
	return res
}
