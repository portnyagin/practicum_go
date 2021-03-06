package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/portnyagin/practicum_go/internal/app/repository/basedbhandler"
)

type PostgresqlHandler struct {
	pool *pgxpool.Pool
}

type PostgressRow struct {
	Rows *pgx.Row
}

func (handler *PostgresqlHandler) Execute(ctx context.Context, statement string, args ...interface{}) error {
	conn, err := handler.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	if len(args) > 0 {
		_, err = conn.Exec(ctx, statement, args...)
	} else {
		_, err = conn.Exec(ctx, statement)
	}

	return err
}

func (handler *PostgresqlHandler) ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error {
	conn, err := handler.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	batch := &pgx.Batch{}
	if len(args) > 0 {
		for _, argset := range args {
			batch.Queue(statement, argset...)
		}
	} else {
		return nil
	}
	br := conn.SendBatch(context.Background(), batch)
	ct, err := br.Exec()
	if err != nil {
		return err
	}
	fmt.Println(ct.RowsAffected())
	return nil
}

func (handler *PostgresqlHandler) QueryRow(ctx context.Context, statement string, args ...interface{}) (basedbhandler.Row, error) {
	var row pgx.Row
	conn, err := handler.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if len(args) > 0 {
		row = conn.QueryRow(ctx, statement, args...)
	} else {
		row = conn.QueryRow(ctx, statement)
	}
	return row, nil
}

func (handler *PostgresqlHandler) Query(ctx context.Context, statement string, args ...interface{}) (basedbhandler.Rows, error) {
	var rows pgx.Rows

	conn, err := handler.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}

	if len(args) > 0 {
		rows, err = conn.Query(ctx, statement, args...)
	} else {
		rows, err = conn.Query(ctx, statement)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (handler *PostgresqlHandler) Close() {
	if handler != nil {
		handler.pool.Close()
	}
}

func NewPostgresqlHandler(ctx context.Context, dataSource string) (*PostgresqlHandler, error) {
	// Format DSN
	//("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)

	poolConfig, err := pgxpool.ParseConfig(dataSource)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	postgresqlHandler := new(PostgresqlHandler)
	postgresqlHandler.pool = pool
	//baseHandler.ErrNotFound = pgx.ErrNoRows
	return postgresqlHandler, nil
}
