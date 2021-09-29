package repository

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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

func (r *PostgresRepository) FindByUser(userID string) ([]model.UserURLs, error) {
	rows, err := r.handler.Query(database.GetURLsByUserID2, userID)
	if err != nil {
		return nil, err
	}
	var resArr []model.UserURLs
	for rows.Next() {
		var rec model.UserURLs
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.ShortURL, &rec.OriginalURL)
		resArr = append(resArr, rec)
		if err != nil {
			return nil, err
		}
	}
	return resArr, nil
}

func (r *PostgresRepository) Save(userID string, shortURL string, originalURL string) error {
	err := r.handler.Execute(database.InsertUserURLWithoutCorrelationID, userID, shortURL, originalURL)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &model.UniqueViolation
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) SaveBatch(src model.UserBatchURLs) error {
	// TODO
	var paramArr [][]interface{}
	for _, obj := range src.List {
		var paramLine []interface{}
		paramLine = append(paramLine, src.UserID)
		paramLine = append(paramLine, obj.CorrelationID)
		paramLine = append(paramLine, obj.OriginalURL)
		paramLine = append(paramLine, obj.ShortURL)
		paramArr = append(paramArr, paramLine)
	}
	err := r.handler.ExecuteBatch(database.InsertUserURL2, paramArr)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &model.UniqueViolation
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) ReadBatch(userID string) (*model.UserBatchURLs, error) {
	if userID == "" {
		return nil, errors.New("userID is empty")
	}
	rows, err := r.handler.Query(database.AllUserURLsWithCorrelationIDByUserID, userID)
	if err != nil {
		return nil, err
	}
	var res model.UserBatchURLs
	res.UserID = userID

	for rows.Next() {
		var e model.Element
		err := rows.Scan(&e.CorrelationID, &e.OriginalURL, &e.ShortURL)
		if err != nil {
			return nil, err
		}
		res.List = append(res.List, e)
	}
	return &res, nil
}