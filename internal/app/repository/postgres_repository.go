package repository

import (
	"context"
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

func (r *PostgresRepository) Ping(ctx context.Context) (bool, error) {
	row, err := r.handler.QueryRow(ctx, "select 10")
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

func (r *PostgresRepository) FindByUser(ctx context.Context, userID string) ([]model.UserURLs, error) {
	rows, err := r.handler.Query(ctx, database.GetURLsByUserID, userID)
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

func (r *PostgresRepository) Save(ctx context.Context, userID string, originalURL string, shortURL string) error {
	err := r.handler.Execute(ctx, database.InsertURL, userID, nil, originalURL, shortURL)
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

func (r *PostgresRepository) SaveBatch(ctx context.Context, src model.UserBatchURLs) error {
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
	err := r.handler.ExecuteBatch(ctx, database.InsertURL, paramArr)
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

func (r *PostgresRepository) FindByShort(ctx context.Context, shortURL string) (string, error) {
	row, err := r.handler.QueryRow(ctx, database.GetOriginalURLByShort, shortURL)
	if err != nil {
		return "", err
	}
	var res string
	err = row.Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}
