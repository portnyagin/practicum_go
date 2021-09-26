package repository

import (
	"github.com/portnyagin/practicum_go/internal/app/model"
	"github.com/portnyagin/practicum_go/internal/app/repository/base_db_handler"
)

type PostgresRepository struct {
	handler base_db_handler.DbHandler
}

func NewPostgresRepository(handler base_db_handler.DbHandler) (*PostgresRepository, error) {
	var repo PostgresRepository
	var err error
	// TODO:
	repo.handler = handler
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *PostgresRepository) Ping() (bool, error) {
	row, err := r.handler.QueryRow("select 10")
	if err != nil {
		return false, err
	}
	var res int
	err = row.Scan(&res)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PostgresRepository) FindByUser(key string) ([]model.UserURLs, error) {
	return nil, nil
}
