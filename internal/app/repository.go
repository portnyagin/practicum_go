package app

import "errors"

type BaseRepository struct {
	store map[string]string
}

func NewBaseRepository() *BaseRepository {
	var r BaseRepository
	r.store = make(map[string]string)
	return &r
}

func (r *BaseRepository) Find(key string) (string, error) {
	if val, ok := r.store[key]; ok {
		return val, nil
	} else {
		return "", errors.New("Can't find value")
	}
}

func (r *BaseRepository) Save(key string, value string) error {
	r.store[key] = value
	return nil
}
