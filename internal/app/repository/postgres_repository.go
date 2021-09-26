package repository

import (
	"fmt"
	"github.com/portnyagin/practicum_go/internal/app/database"
	"github.com/portnyagin/practicum_go/internal/app/model"
	"github.com/portnyagin/practicum_go/internal/app/repository/basedbhandler"
)

type PostgresRepository struct {
	handler basedbhandler.DBHandler
}

func NewPostgresRepository(handler basedbhandler.DBHandler) (*PostgresRepository, error) {
	var repo PostgresRepository
	repo.handler = handler
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

func InitDatabase(h basedbhandler.DBHandler) error {
	err := h.Execute(database.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("Database structure created successfully")
	return nil
}

func ClearDatabase(h basedbhandler.DBHandler) error {
	err := h.Execute(database.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
