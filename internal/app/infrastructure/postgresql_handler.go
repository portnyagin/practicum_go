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
	ctx  context.Context
}

type PostgressRow struct {
	Rows *pgx.Row
}

func (handler *PostgresqlHandler) Execute(statement string, args ...interface{}) error {
	conn, err := handler.pool.Acquire(handler.ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	if len(args) > 0 {
		_, err = conn.Exec(handler.ctx, statement, args...)
	} else {
		_, err = conn.Exec(handler.ctx, statement)
	}

	return err
}

func (handler *PostgresqlHandler) ExecuteBatch(statement string, args [][]interface{}) error {
	conn, err := handler.pool.Acquire(handler.ctx)
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

func (handler *PostgresqlHandler) QueryRow(statement string, args ...interface{}) (basedbhandler.Row, error) {
	var row pgx.Row
	conn, err := handler.pool.Acquire(handler.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if len(args) > 0 {
		row = conn.QueryRow(handler.ctx, statement, args...)
	} else {
		row = conn.QueryRow(handler.ctx, statement)
	}

	return row, nil
}

func (handler *PostgresqlHandler) Query(statement string, args ...interface{}) (basedbhandler.Rows, error) {
	var rows pgx.Rows

	conn, err := handler.pool.Acquire(handler.ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}

	if len(args) > 0 {
		rows, err = conn.Query(handler.ctx, statement, args...)
	} else {
		rows, err = conn.Query(handler.ctx, statement)
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
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	postgresqlHandler := new(PostgresqlHandler)
	postgresqlHandler.ctx = ctx
	postgresqlHandler.pool = pool
	//baseHandler.ErrNotFound = pgx.ErrNoRows
	return postgresqlHandler, nil
}