package postgres

import (
	"avito-internship/pkg/e"
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgDatabase struct {
	Pool *pgxpool.Pool
}

func NewPgDatabase(pool *pgxpool.Pool) *PgDatabase {
	return &PgDatabase{Pool: pool}
}

func Connect() (*PgDatabase, error) {
	const op = "PgDatabase.Connect"

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, e.Wrap(op, err)
	}

	return NewPgDatabase(pool), nil
}

func (db *PgDatabase) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
